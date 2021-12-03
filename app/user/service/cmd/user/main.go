package main

import (
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

	app, cleanup, err := initApp(bc.Server, bc.Data, logger, logger.(*zapLog.ZapLogger).Logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
