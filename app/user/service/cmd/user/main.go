package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/go-uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	v1 "chat/api/user/service/v1"
	"chat/app/user/service/internal/conf"
	myConf "chat/pkg/conf"
	zapLog "chat/pkg/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, r registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	flag.Parse()

	var bc conf.Bootstrap
	myConf.NewConf(flagconf, &bc)
	//myConf.NewConsulConf(flagconf, &bc)

	Name = v1.User_ServiceDesc.ServiceName
	Version = bc.Project.Version
	if generateUUID, err := uuid.GenerateUUID(); err == nil {
		id += "-" + generateUUID
	}

	logger := zapLog.Logger(bc.Project.Mode, bc.Log.LogPath, int(bc.Log.MaxSize), int(bc.Log.MaxBackups), int(bc.Log.MaxSize), bc.Log.Compress)
	zapLogger := logger.(*zapLog.ZapLogger).Logger
	logger = log.With(logger, "trace_id", tracing.TraceID())
	logger = log.With(logger, "span_id", tracing.SpanID())
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "service.id", id)
	logger = log.With(logger, "service.name", Name)
	logger = log.With(logger, "service.version", Version)

	url := "http://127.0.0.1:14268/api/traces"
	err := setTracerProvider(url)
	if err != nil {
		logger.Log(log.LevelError, err)
		return
	}

	app, cleanup, err := initApp(bc.Server, bc.Data, logger, zapLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

// Set global trace provider
func setTracerProvider(url string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}
