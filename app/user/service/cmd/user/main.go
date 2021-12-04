package main

import (
	"chat/pkg/trace"
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/go-uuid"

	v1 "chat/api/user/service/v1"
	"chat/app/user/service/internal/conf"
	myConf "chat/pkg/conf"
	zapLog "chat/pkg/log"
	logConf "chat/pkg/log/conf"
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

	name = v1.User_ServiceDesc.ServiceName
	version = bc.Project.Version
	if generateUUID, err := uuid.GenerateUUID(); err == nil {
		id += "-" + generateUUID
	}

	logger := zapLog.Logger(bc.Project.Mode, lc.ZapLog)
	zapLogger := logger.GetZapLogger()
	newLogger := logger.GetLogger(id, name, version)

	err := trace.SetTracerProvider(name, tc.Jaeger.Url)
	if err != nil {
		panic(err)
	}

	app, cleanup, err := initApp(bc.Server, &rc, bc.Data, newLogger, zapLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
