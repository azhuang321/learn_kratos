package data

import (
	v1 "chat/api/group/service/v1"
	"chat/app/user/service/internal/biz"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"go.uber.org/zap"
	"sync"
	"time"
)

//todo  是否应该使用 zap log  到日志收集环节再看
type userRepo struct {
	data *Data
	log  *zap.Logger
}

func NewUserRepo(data *Data, logger *zap.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  logger,
	}
}

// Create 测试mysql有链接池 最低保持两个空闲
func (u userRepo) Create(ctx context.Context, user *biz.User) error {
	_, err := u.data.db.User.Create().SetUsername(user.Username).SetPassword(user.Password).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

var localExecDict = make(map[string]map[string]interface{})

type DemoListener struct {
	localTrans       *sync.Map
	transactionIndex int32
}

// ExecuteLocalTransaction 本地事务 基于回调,发送消息成功后 立即调用此回调方法,依据返回的状态 来确认 发送的消息是否可进行消费
//本地事务执行成功，返回COMMIT，MQ服务器会把HALF MESSAGE塞到真实的队列中，并同时向RMQ_SYS_TRANS_OP_HALF_TOPIC塞一条消息。此时消费端可以消费消息。
//本地事务执行失败，返回ROLLBACK，表示消息需要回滚，会直接删掉HALF MESSAGE，并并同时向RMQ_SYS_TRANS_OP_HALF_TOPIC塞一条消息。
//本地事务返回未知，返回UN_KNOW，不做任何处理。MQ服务器会通过TransactionalMessageCheckService定时调用checkLocalTransaction检查本地事务。（最多重试15次，超过了默认丢弃此消息）
func (dl *DemoListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("ExecuteLocalTransaction.............................................")
	msgBody := msg.Body
	msgs := make(map[string]interface{})
	_ = json.Unmarshal(msgBody, &msgs)
	fmt.Println(msgs)
	return primitive.UnknowState

	//1.创建用户
	if false {
	}

	//创建用户
	if true {
		//return primitive.CommitMessageState
	}
	time.Sleep(time.Second * 60)
	//return primitive.RollbackMessageState
	return 5
}

// CheckLocalTransaction 当长时间没有得到此事务消息的本地事务状态,就会进行消息回查,当出现宕机的异常情况 消息没有反馈明确状态时,通过回查来确认用户是否创建成功
func (dl *DemoListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("CheckLocalTransaction.............................................")
	msgBody := msg.Body
	msgs := make(map[string]interface{})
	_ = json.Unmarshal(msgBody, &msgs)
	fmt.Println(msgs, "..........................................")

	return primitive.CommitMessageState

	//查询本地数据库,看一下用户创建是否成功
	if true {
		return primitive.RollbackMessageState //
	} else {
		return primitive.CommitMessageState //
	}
}

func NewDemoListener() *DemoListener {
	return &DemoListener{
		localTrans: new(sync.Map),
	}
}

func (u userRepo) GroupInfo(ctx context.Context, user *biz.User) error {
	userId := "1"
	localExecDict[userId] = map[string]interface{}{}

	e, b := sentinel.Entry("breaker_err_count")
	if b != nil {
		u.log.Error("breaker")
		return ratelimit.ErrLimitExceed
	} else {
		rlog.SetLogLevel("error")
		p, err := rocketmq.NewTransactionProducer(
			NewDemoListener(),
			producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
			producer.WithRetry(1),
		)
		if err != nil {
			u.log.Error(err.Error())
			return err
		}

		err = p.Start()
		if err != nil {
			u.log.Error(err.Error())
			return err
		}
		msg := map[string]string{
			"uid":      "1",
			"group_id": "1",
		}
		msgBody, _ := json.Marshal(msg)
		fmt.Println("before sendMessage.............................................")

		res, err := p.SendMessageInTransaction(context.TODO(), primitive.NewMessage("create_group", msgBody))
		if err != nil {
			u.log.Error(err.Error())
			return err
		}
		fmt.Println("after sendMessage.............................................")

		if res.Status != primitive.SendOK {
			return errors.New("发送消息失败")
		}
		fmt.Printf("send status:%+v,message_id:%v", res.Status, res.MsgID)

		time.Sleep(time.Second * 600)

		for {
			if _, ok := localExecDict[userId]; ok {
				u.log.Warn("回调结束...")
				break
			}
			u.log.Warn("等待rocket mq 回调中")
			time.Sleep(time.Second)
		}

		err = p.Shutdown()
		if err != nil {
			fmt.Printf("shutdown producer error: %s", err.Error())
		}
		fmt.Println(localExecDict[userId])

		_, err = u.data.gc.GetGroupInfo(ctx, &v1.GetGroupInfoRequest{Id: 1})
		if err != nil {
			u.log.Error(err.Error()) //todo  uberrate 测试出现错误
			sentinel.TraceError(e, err)
		}
		e.Exit()
	}

	return nil
}
