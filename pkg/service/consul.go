package service

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
)

type Consul struct {
	Name    string
	Address string
	Port    string
}

func (consul *Consul) config() *api.Config {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", consul.Address, consul.Port)

	return config
}

func (consul *Consul) id() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s:%s:%s:%s", hostname, consul.Name, consul.Address, consul.Port)
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
