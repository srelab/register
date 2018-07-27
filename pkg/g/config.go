package g

import (
	"github.com/urfave/cli"
	"sync"
)

type DockerConfig struct {
	Endpoint string
}

type LogConfig struct {
	Dir   string
	Level string
}

type ConsulConfig struct {
	Host string
	Port string
}

type GatewayConfig struct {
	Host string
	Port string
}

type PrivilegeConfig struct {
	Host string
	Port string
}

type GlobalConfig struct {
	Name        string
	Concurrency int

	Log       *LogConfig
	Consul    *ConsulConfig
	Docker    *DockerConfig
	Gateway   *GatewayConfig
	Privilege *PrivilegeConfig
}

var (
	config *GlobalConfig
	lock   = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(ctx *cli.Context) {
	config = &GlobalConfig{
		Name:        NAME,
		Concurrency: ctx.Int("concurrency"),
		Consul: &ConsulConfig{
			Host: ctx.String("consul.host"),
			Port: ctx.String("consul.port"),
		},
		Docker: &DockerConfig{
			Endpoint: ctx.String("docker.endpoint"),
		},
		Log: &LogConfig{
			Dir:   ctx.String("log.dir"),
			Level: ctx.String("log.level"),
		},
		Gateway: &GatewayConfig{
			Host: ctx.String("gateway.host"),
			Port: ctx.String("gateway.port"),
		},
		Privilege: &PrivilegeConfig{
			Host: ctx.String("privilege.host"),
			Port: ctx.String("privilege.port"),
		},
	}
}
