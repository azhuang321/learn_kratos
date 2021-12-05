package rate

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
	"github.com/go-kratos/aegis/ratelimit"
)

type SentinelRate struct{}

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Println()
	fmt.Println()
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
	fmt.Println()
	fmt.Println()
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Println()
	fmt.Println()
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
	fmt.Println()
	fmt.Println()
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Println()
	fmt.Println()
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
	fmt.Println()
	fmt.Println()
}

func NewSentinelRate() *SentinelRate {
	// We should initialize Sentinel first.
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		panic(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "test",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              30,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		panic(err)
	}

	// https://sentinelguard.io/zh-cn/docs/golang/circuit-breaking.html
	// Register a state change listener so that we could observer the state change of the internal circuit breaker.
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorCount=50
		{
			Resource:         "breaker_err_count",
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   5000,  //熔断触发后持续的时间
			MinRequestAmount: 1,     //触发熔断的最小请求数目
			StatIntervalMs:   10000, //统计的时间窗口长度
			Threshold:        2,     //是错误计数的阈值
		},
		{
			Resource:                     "breaker_err_ratio",
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               3000, //熔断触发后持续的时间
			MinRequestAmount:             10,   //触发熔断的最小请求数目
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
			Threshold:                    0.4, // 表示的是错误比例的阈值(小数表示，比如0.1表示10%)
		},
		{
			Resource:                     "slow_request_ratio",
			Strategy:                     circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             10,
			StatIntervalMs:               5000,
			StatSlidingWindowBucketCount: 10,
			MaxAllowedRtMs:               50,  //慢调用熔断策略  如果请求的response time小于或等于MaxAllowedRtMs，那么就不是慢调用
			Threshold:                    0.5, // 如果当前资源的慢调用比例如果高于Threshold，那么熔断器就会断开；否则保持闭合状态
		},
	})
	if err != nil {
		panic(err)
	}

	return &SentinelRate{}
}

func (ur *SentinelRate) Allow() (df ratelimit.DoneFunc, err error) {
	e, b := sentinel.Entry("test", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		fmt.Println("被限流了.....", b.Error())
		return nil, ratelimit.ErrLimitExceed
	}
	return func(ratelimit.DoneInfo) { e.Exit() }, nil
}
