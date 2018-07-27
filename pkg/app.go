package pkg

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/srelab/gpool"
	"github.com/srelab/register/pkg/g"
	"github.com/srelab/register/pkg/logger"
	"github.com/srelab/register/pkg/service"
	"github.com/srelab/register/pkg/store"
)

func Start() error {
	client, err := docker.NewClient(g.Config().Docker.Endpoint)
	if err != nil {
		logger.Fatal("failed to get docker client:", err)
	}

	listener := make(chan *docker.APIEvents)
	if err := client.AddEventListener(listener); err != nil {
		logger.Fatal("failed to listener docker event:", err)
	}

	defer func() {
		err = client.RemoveEventListener(listener)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	p := gpool.NewPool(g.Config().Concurrency*2, g.Config().Concurrency)
	defer p.Release()

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		logger.Fatal("failed to get containers:", err)
	}

	for _, container := range containers {
		if err := store.Container.Set(client, container.ID); err != nil {
			logger.Errorf("failed to container [id:%v] information:%v", container.ID, err)
			continue
		}
	}

	for {
		select {
		case e, ok := <-listener:
			if !ok {
				logger.Fatal("'%s' not found or permission denied...", g.Config().Docker.Endpoint)
			}

			p.JobQueue <- func() {
				if err := handle(e, client); err != nil {
					logger.Error(err)
					(&service.Privilege{
						Host: g.Config().Privilege.Host,
						Port: g.Config().Privilege.Port,
					}).WechatMsgSend(err)
				}
			}
		}
	}

	return nil
}

func handle(event *docker.APIEvents, client *docker.Client) error {
	switch event.Action {
	case "start", "unpause":
		info, err := store.Container.Get(event.ID, client)
		if err != nil {
			return err
		}

		gateway := &service.Gateway{
			Name: info["SERVICE_NAME"].(string),
			Host: info["DOCKER_ADDRESS"].(string),
			Port: info["SERVICE_PORT"].(int),
		}

		logger.Infof("[Event][%s][%s][%s] - try to register for a service",
			event.Action, gateway.Name, gateway.Port)

		if err := gateway.Register(); err != nil {
			return fmt.Errorf("failed to register service[%s] to API gateway, because %s", gateway.Name, err)
		}

		consul := &service.Consul{
			Name:    gateway.Name,
			Address: g.Config().Consul.Host,
			Port:    g.Config().Consul.Port,
		}

		if err := consul.Register(); err != nil {
			return fmt.Errorf("failed to register service[%s] to consul, because %s", gateway.Name, err)
		}
	case "pause", "die":
		info, err := store.Container.Get(event.ID, client)
		if err != nil {
			return err
		}

		gateway := &service.Gateway{
			Name: info["SERVICE_NAME"].(string),
			Host: info["DOCKER_ADDRESS"].(string),
			Port: info["SERVICE_PORT"].(int),
		}

		logger.Infof("[Event][%s][%s][%s] - try to unregister for a service",
			event.Action, gateway.Name, gateway.Port)

		if err := gateway.UnRegister(); err != nil {
			return fmt.Errorf("failed to unregister service[%s] from API gateway, because %s", gateway.Name, err)
		}

		consul := &service.Consul{
			Name:    gateway.Name,
			Address: g.Config().Consul.Host,
			Port:    g.Config().Consul.Port,
		}

		if err := consul.UnRegister(); err != nil {
			return fmt.Errorf("failed to unregister service[%s] from consul, because %s", gateway.Name, err)
		}
	}

	return nil
}
