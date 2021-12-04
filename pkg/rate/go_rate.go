package rate

import (
	"fmt"
	"github.com/go-kratos/aegis/ratelimit"
	"golang.org/x/time/rate"
	"time"
)

type GoRate struct {
	Limiter *rate.Limiter
}

func NewGoRate() *GoRate {
	r := rate.Every(1 * time.Second)
	limiter := rate.NewLimiter(r, 1)
	return &GoRate{Limiter: limiter}
}

func (gr *GoRate) Allow() (df ratelimit.DoneFunc, err error) {
	if !gr.Limiter.Allow() {
		fmt.Println("被限流了.....")
		return nil, ratelimit.ErrLimitExceed
	}

	return func(ratelimit.DoneInfo) {
	}, nil
}
