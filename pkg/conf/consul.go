package conf

import (
	"gopkg.in/yaml.v3"
	"strconv"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/hashicorp/consul/api"

	conf "chat/pkg/conf/proto"
	ccf "github.com/go-kratos/kratos/contrib/config/consul/v2"
)

func NewConsulConf(configPath string, projectConf interface{}) {
	c := config.New(config.WithSource(file.NewSource(configPath + "/consul.yaml")))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var consul conf.Consul
	if err := c.Scan(&consul); err != nil {
		panic(err)
	}

	client, err := api.NewClient(&api.Config{
		Address: consul.Ip + ":" + strconv.Itoa(int(consul.Port)),
	})

	source, err := ccf.New(client, ccf.WithPath(consul.ConfigPath))
	if err != nil {
		panic(err)
	}

	kvs, err := source.Load()
	if err != nil {
		panic(err)
	}
	if len(kvs) < 1 {
		panic("config err")
	}

	cc := config.New(
		config.WithSource(
			source,
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	defer cc.Close()

	if err = cc.Load(); err != nil {
		panic(err)
	}

	if err := cc.Scan(projectConf); err != nil {
		panic(err)
	}
}
