package rate

import (
	"github.com/go-kratos/aegis/ratelimit"
	uRate "go.uber.org/ratelimit"
)

type UberRate struct {
	Limiter uRate.Limiter
}

func NewUberRate() *UberRate {
	limiter := uRate.New(1) // per second
	return &UberRate{Limiter: limiter}
}

func (ur *UberRate) Allow() (df ratelimit.DoneFunc, err error) {
	ur.Limiter.Take()
	//if !ur.Limiter.Allow() {
	//	fmt.Println("被限流了.....")
	//	return nil, ratelimit.ErrLimitExceed
	//}

	return func(ratelimit.DoneInfo) {
	}, nil
}
