package registry

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	etcdClient "go.etcd.io/etcd/client/v3"
)

func NewEtcdRegistry(logger log.Logger) registry.Registrar {
	logH := log.NewHelper(logger)
	client, err := etcdClient.New(etcdClient.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		logH.Error(err.Error())
		return nil
	}
	return etcd.New(client)
}
