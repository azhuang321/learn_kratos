package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

func NewConf(configPath string, projectConf ...interface{}) {
	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	if len(projectConf) < 1 {
		return
	}
	for _, v := range projectConf {
		if err := c.Scan(v); err != nil {
			panic(err)
		}
	}
}
