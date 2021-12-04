package rate

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/go-kratos/aegis/ratelimit"
	"time"
)

type DiDiRate struct {
	Limiter *limiter.Limiter
}

// 只支持http
func NewDiDiRate() *DiDiRate {
	lmt := tollbooth.NewLimiter(1, nil)
	// 或创建一个带有可过期令牌桶的限制器 此设置意味着：创建一个 1 requestsecond 限制器，其中的每个令牌桶将在最初设置后 1 小时过期。
	lmt = tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	// 配置查找 IP 地址的位置列表。默认情况下，它是：“RemoteAddr”、“X-Forwarded-For”、“X-Real-IP” 如果您的应用程序在代理之后，请先设置“X-Forwarded-For”。
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})

	// 仅限制 GET 和 POST 请求。
	lmt.SetMethods([]string{"GET", "POST"})

	// 基于基本身份验证用户名的限制。您可以在加载时添加它们，或者稍后在处理请求时添加它们。
	lmt.SetBasicAuthUsers([]string{"bob", "jane", "didip", "vip"})
	// 您也可以稍后删除它们。
	lmt.RemoveBasicAuthUsers([]string{"vip"})

	// 限制包含某些值的请求标头。您可以在加载时添加它们，或者稍后在处理请求时添加它们。
	lmt.SetHeader("X-Access-Token", []string{"abc123", "xyz098"})
	// 您可以一次删除所有条目。
	lmt.RemoveHeader("X-Access-Token")
	// 或者删除特定的。
	lmt.RemoveHeaderEntries("X-Access-Token", []string{"limitless-token"})

	return &DiDiRate{Limiter: lmt}
}

func (gr *DiDiRate) Allow() (df ratelimit.DoneFunc, err error) {
	//if !gr.Limiter.Allow() {
	//	fmt.Println("被限流了.....")
	//	return nil, ratelimit.ErrLimitExceed
	//}

	return func(ratelimit.DoneInfo) {
	}, nil
}
