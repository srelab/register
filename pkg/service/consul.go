package service

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/srelab/register/pkg/g"
	"os"
)

type Consul struct {
	Name    string
	Address string
	Port    int
}

func (consul *Consul) config() *api.Config {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", g.Config().Consul.Host, g.Config().Consul.Port)

	return config
}

func (consul *Consul) id() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s:%s:%s:%d", hostname, consul.Name, consul.Address, consul.Port)
}

func (consul *Consul) Register(eid string) error {
	client, err := api.NewClient(consul.config())
	if err != nil {
		return err
	}

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:                consul.id(),
		Name:              consul.Name,
		Address:           consul.Address,
		EnableTagOverride: false,
		Port:              consul.Port,

		Checks: api.AgentServiceChecks{
			&api.AgentServiceCheck{
				Name:     "tcp@" + eid[0:12],
				TCP:      fmt.Sprintf("%s:%d", consul.Address, consul.Port),
				Interval: "10s",
				Timeout:  "3s",
				DeregisterCriticalServiceAfter: "3m",
			},
			&api.AgentServiceCheck{
				Name:     "ping@" + eid[0:12],
				Args:     []string{"/welab.co/bin/consul-health", "ping"},
				Interval: "10s",
				Timeout:  "3s",
				DeregisterCriticalServiceAfter: "3m",
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (consul *Consul) UnRegister() error {
	client, err := api.NewClient(consul.config())
	if err != nil {
		return err
	}

	err = client.Agent().ServiceDeregister(consul.id())

	if err != nil {
		return err
	}

	return nil
}
