package registry

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	etcdClient "go.etcd.io/etcd/client/v3"

	zaplog "chat/pkg/log"
)

func NewEtcdRegistry(logger log.Logger) registry.Registrar {
	myLogger := logger.(*zaplog.ZapLogger).Logger
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		myLogger.Error(err.Error())
		return nil
	}
	return etcd.New(client)
}
