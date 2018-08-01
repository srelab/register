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

func (consul *Consul) Register() error {
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
