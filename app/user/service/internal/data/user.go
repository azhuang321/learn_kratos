package data

import (
	v1 "chat/api/group/service/v1"
	"chat/app/user/service/internal/biz"
	"context"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"go.uber.org/zap"
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

func (u userRepo) GroupInfo(ctx context.Context, user *biz.User) error {
	e, b := sentinel.Entry("breaker_err_count")
	if b != nil {
		u.log.Error("breaker")
		return ratelimit.ErrLimitExceed
	} else {
		_, err := u.data.gc.GetGroupInfo(ctx, &v1.GetGroupInfoRequest{Id: 1})
		if err != nil {
			u.log.Error(err.Error()) //todo  uberrate 测试出现错误
			sentinel.TraceError(e, err)
		}
		e.Exit()
	}

	return nil
}
