package util

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strconv"
	"strings"
)

func GetContainerInfo(container *docker.Container) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	for _, env := range container.Config.Env {
		key, val := env2kv(env, "=")
		info[key] = val
	}

	if !check(info) {
		return nil, errors.New("keys and container info not match")
	}

	info["DOCKER_ADDRESS"] = eth0Address(container.ID)
	return info, nil
}

func env2kv(env, sep string) (string, string) {
	s := strings.Split(env, sep)
	return s[0], s[1]
}

func check(info map[string]interface{}) bool {
	keys := []string{"SERVICE_NAME", "SERVICE_PORT", "CONTEXT_PATH", "COMPATIBLE"}

	for _, key := range keys {
		if _, ok := info[key]; !ok {
			return false
		}

		if key == "SERVICE_PORT" {
			val, err := strconv.Atoi(info[key].(string))

			if err != nil {
				return false
			}

			info[key] = val
		}
	}

	return true
}

func eth0Address(id string) string {
	//for network := range networks {
	//	if obj, exists := networks[network]; exists {
	//		if obj.IPAddress != "" {
	//			return obj.IPAddress
	//		}
	//
	//		continue
	//	}
	//}

	cmd := fmt.Sprintf(`docker exec %s ifconfig eth0 | grep -oP '\d.+(?=  (Bcast:|netmask))'`, id)
	address, err := CmdOutBytes("/bin/sh", "-c", cmd)
	if err != nil {
		return ""
	}

	return strings.TrimSuffix(string(address), "\n")
}
