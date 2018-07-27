package store

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/srelab/register/pkg/util"
	"sync"
)

type container struct {
	mutex *sync.RWMutex
	m     map[string]map[string]interface{}
}

func (c container) Get(id string, client *docker.Client) (map[string]interface{}, error) {
	if info, ok := c.m[id]; ok {
		return info, nil
	}

	if err := c.Set(client, id); err != nil {
		return nil, fmt.Errorf("set container(%v) info failed:%v", id, err)
	}

	if info, ok := c.m[id]; ok {
		return info, nil
	}

	return nil, errors.New("container not found in store")
}

func (c container) Add(id string, info map[string]interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.m[id] = info
}

func (c container) Set(client *docker.Client, id string) error {
	inspect, err := client.InspectContainer(id)
	if err != nil {
		return fmt.Errorf("get container(%v) info failed", id)
	}

	info, err := util.GetContainerInfo(inspect)
	if err != nil {
		return err
	}

	c.Add(id, info)
	return nil
}

func (c container) Remove(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.m, id)
}

var Container = container{mutex: new(sync.RWMutex), m: make(map[string]map[string]interface{})}
