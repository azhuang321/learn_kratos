package data

import (
	v1 "chat/api/group/service/v1"
	"chat/app/user/service/internal/conf"
	"chat/app/user/service/internal/data/ent"
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector/wrr"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	grpc2 "google.golang.org/grpc"

	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewEntClient, NewUserRepo, NewGroupClient, NewDiscovery)

// Data .
type Data struct {
	db *ent.Client
	gc v1.GroupClient
}

// NewData .
func NewData(db *ent.Client, gc *grpc2.ClientConn, c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		db.Close()
		gc.Close()
		fmt.Println("closing the data resources")
	}
	return &Data{db: db, gc: v1.NewGroupClient(gc)}, cleanup, nil
}

func NewEntClient(c *conf.Data, logger log.Logger) *ent.Client {
	log := log.NewHelper(logger)
	drv, err := sql.Open(
		c.Database.Driver,
		c.Database.Source,
	)
	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		log.WithContext(ctx).Info(i...)
		tracer := otel.Tracer("ent.")
		kind := trace.SpanKindServer
		_, span := tracer.Start(ctx,
			"Query",
			trace.WithAttributes(
				attribute.String("sql", fmt.Sprint(i...)),
			),
			trace.WithSpanKind(kind),
		)
		span.End()
	})
	client := ent.NewClient(ent.Driver(sqlDrv))
	if err != nil {
		fmt.Println(err)
		log.Errorf("failed opening connection to sqlite: " + err.Error())
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Errorf("failed creating schema resources:" + err.Error())
	}
	return client
}

func NewGroupClient(r registry.Discovery) *grpc2.ClientConn {
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///"+v1.Group_ServiceDesc.ServiceName),
		grpc.WithDiscovery(r),

		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(), //todo 替换
		),
		grpc.WithTimeout(2*time.Second),
		grpc.WithBalancerName(wrr.Name),
		grpc.WithOptions(grpc2.WithStatsHandler(&tracing.ClientHandler{})), // for tracing remote ip recording
	)
	if err != nil {
		panic(err)
	}

	return conn
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
