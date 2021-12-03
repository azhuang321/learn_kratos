package conf

import (
	ncf "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v3"

	conf "chat/pkg/conf/proto"
)

func NewNacosConf(configPath string, projectConf interface{}) {
	c := config.New(config.WithSource(file.NewSource(configPath + "/nacos.yaml")))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var nacos conf.Nacos
	if err := c.Scan(&nacos); err != nil {
		panic(err)
	}

	var sc []constant.ServerConfig
	for _, val := range nacos.ServerConfigs {
		sc = append(sc, *constant.NewServerConfig(val.IpAddr, val.Port))
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacos.ClientConfig.NamespaceId,
		TimeoutMs:           nacos.ClientConfig.TimeoutMs,
		NotLoadCacheAtStart: nacos.ClientConfig.NotLoadCacheAtStart,
		LogDir:              nacos.ClientConfig.LogDir,
		CacheDir:            nacos.ClientConfig.CacheDir,
		RotateTime:          nacos.ClientConfig.RotateTime,
		MaxAge:              nacos.ClientConfig.MaxAge,
		LogLevel:            nacos.ClientConfig.LogLevel,
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	dataID := nacos.NodeConfig.DataId
	group := nacos.NodeConfig.Group

	nc := config.New(
		config.WithSource(
			ncf.NewConfigSource(client, ncf.WithGroup(group), ncf.WithDataID(dataID)),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	defer nc.Close()

	if err = nc.Load(); err != nil {
		panic(err)
	}

	if err := nc.Scan(projectConf); err != nil {
		panic(err)
	}
}
