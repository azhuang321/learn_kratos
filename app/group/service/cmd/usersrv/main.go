package main

import (
	v1 "chat/api/group/service/v1"
	myConf "chat/pkg/conf"
	logConf "chat/pkg/log/conf"
	"flag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/go-uuid"

	"chat/app/group/service/internal/conf"
	zapLog "chat/pkg/log"
)

// go build -ldflags "-X main.version=x.y.z"
var (
	// name is the name of the compiled software.
	name string
	// version is the version of the compiled software.
	version string
	// flagConf is the config flag.
	flagConf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, r registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(name),
		kratos.Version(version),
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
	var rc conf.Registry
	var lc logConf.Log
	var tc conf.Trace

	myConf.NewConf(flagConf, &bc, &rc, &lc, &tc)

	name = v1.Group_ServiceDesc.ServiceName
	version = bc.Project.Version
	if generateUUID, err := uuid.GenerateUUID(); err == nil {
		id += "-" + generateUUID
	}

	logger := zapLog.Logger(bc.Project.Mode, lc.ZapLog)
	zapLogger := logger.GetZapLogger()
	newLogger := logger.GetLogger(id, name, version)

	url := "http://127.0.0.1:14268/api/traces"
	if os.Getenv("jaeger_url") != "" {
		url = os.Getenv("jaeger_url")
	}
	err := setTracerProvider(url)
	if err != nil {
		logger.Log(log.LevelError, err)
		return
	}

	app, cleanup, err := initApp(bc.Server, bc.Data, newLogger, zapLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

// set trace provider
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
			semconv.ServiceNameKey.String(name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}
