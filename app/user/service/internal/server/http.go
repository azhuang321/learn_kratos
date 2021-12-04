package server

import (
	v1 "chat/api/user/service/v1"
	"chat/app/user/service/internal/conf"
	"chat/app/user/service/internal/service"
	"chat/pkg/rate"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			ratelimit.Server(ratelimit.WithLimiter(rate.NewSentinelRate())),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	openAPIHandler := openapiv2.NewHandler()

	srv := http.NewServer(opts...)
	// /q/swagger-ui/
	srv.HandlePrefix("/q/", openAPIHandler)
	v1.RegisterUserHTTPServer(srv, greeter)
	return srv
}
