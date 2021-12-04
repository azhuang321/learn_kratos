package rate

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/go-kratos/aegis/ratelimit"
)

type SentinelRate struct{}

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
			Threshold:              3,
			StatIntervalInMs:       6000,
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
