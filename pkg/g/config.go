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

type GlobalConfig struct {
	Name        string
	Concurrency int

	Log     *LogConfig
	Consul  *ConsulConfig
	Docker  *DockerConfig
	Gateway *GatewayConfig
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
	docker := &DockerConfig{
		Endpoint: ctx.String("docker.endpoint"),
	}

	log := &LogConfig{
		Dir:   ctx.String("log.dir"),
		Level: ctx.String("log.level"),
	}

	gateway := &GatewayConfig{
		Host: ctx.String("gateway.host"),
		Port: ctx.String("gateway.port"),
	}

	config = &GlobalConfig{
		Name:        NAME,
		Concurrency: ctx.Int("concurrency"),
		Docker:      docker,
		Log:         log,
		Gateway:     gateway,
	}
}
