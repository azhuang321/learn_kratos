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
var ProviderSet = wire.NewSet(NewData, NewEntClient, NewUserRepo, NewGroupService)

// Data .
type Data struct {
	db        *ent.Client
	groupConn *grpc2.ClientConn
}

// NewData .
func NewData(db *ent.Client, groupConn *grpc2.ClientConn, c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		db.Close()
		groupConn.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, groupConn: groupConn}, cleanup, nil
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

func NewGroupService() *grpc2.ClientConn {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	dis := consul.New(client)

	// new grpc client
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///"+v1.Group_ServiceDesc.ServiceName),
		grpc.WithDiscovery(dis),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
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
