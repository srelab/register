package pkg

import (
	"fmt"
	"github.com/srelab/register/pkg/bll"
	"github.com/srelab/gpool"
	"github.com/srelab/register/pkg/g"
	"github.com/srelab/register/pkg/logger"
	. "github.com/srelab/register/pkg/store"
	"github.com/fsouza/go-dockerclient"

)

func Start() error {
	client, err := docker.NewClient(g.Config().Docker.Endpoint)
	if err != nil {
		logger.Error("failure to get docker client instance:", err)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		logger.Error("failure to get containers:", err)
	}

	for _, container := range containers {
		if err := ContainerStore.Set(client, container.ID); err != nil {
			logger.Errorf("Get container(%v) information failed:%v", container.ID, err)
			continue
		}
	}

	listener := make(chan *docker.APIEvents)
	if err := client.AddEventListener(listener); err != nil {
		logger.Fatal("add event listener failure:", err)
	}

	defer func() {
		err = client.RemoveEventListener(listener)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	p := gpool.NewPool(g.Config().Concurrency*2, g.Config().Concurrency)
	defer p.Release()

	for {
		select {
		case e, ok := <-listener:
			if !ok {
				logger.Errorf("'%s' not found or permission denied...", g.Config().Docker.Endpoint)
			}
			p.JobQueue <- func() {
				if err := handle(e, client); err != nil {
					logger.Error(err)
				}
			}
		}
	}

	return nil
}

func handle(event *docker.APIEvents, client *docker.Client) error {
	switch event.Action {
	case "start", "unpause":
		logger.Infof("Event:%s --> Try to register service", event.Action)
		info, err := ContainerStore.Get(event.ID, client)
		if err != nil {
			return err
		}

		gw := &bll.GatewayEntry{
			Name:       info["SERVICE_NAME"].(string),
			Host:    	info["DOCKER_ADDRESS"].(string),
			Port:       info["SERVICE_PORT"].(int),
		}

		if err := gw.Register(); err != nil {
			return fmt.Errorf("failed register service to gateway: %s", err)
		}
	case "pause", "die":
		logger.Infof("Event:%s --> Try to deregister service", event.Action)
		info, err := ContainerStore.Get(event.ID, client)
		if err != nil {
			return err
		}

		gw := &bll.GatewayEntry{
			Name:       info["SERVICE_NAME"].(string),
			Host:    	info["DOCKER_ADDRESS"].(string),
			Port:       info["SERVICE_PORT"].(int),
		}

		if err := gw.UnRegister(); err != nil {
			return fmt.Errorf("failed deregister service to gateway: %s", err)
		}
	}

	return nil
}
