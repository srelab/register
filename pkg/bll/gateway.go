package bll

import (
	"fmt"
	"github.com/srelab/register/pkg/g"
	"github.com/go-resty/resty"
	"net/http"
	"time"
)

type GatewayEntry struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
}

func (entry *GatewayEntry) URL(route string) string {
	u := fmt.Sprintf("http://%s:%s", g.Config().Gateway.Host, g.Config().Gateway.Port)
	return u + route
}

func (entry *GatewayEntry) request() *resty.Request {
	resty.SetRetryCount(3).SetRetryWaitTime(5 * time.Second).SetRetryMaxWaitTime(20 * time.Second)
	return resty.R()
}

func (entry *GatewayEntry) Register() error {
	resp, err := entry.request().
		SetHeader("Content-Type", "application/json").
		SetBody(entry).
		Post(entry.URL(fmt.Sprintf("/upstreams/%s/register", entry.Name)))

	if err != nil || resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server request data error: %v", err)
	}

	return nil
}


func (entry *GatewayEntry) UnRegister() error {
	resp, err := entry.request().
		SetHeader("Content-Type", "application/json").
		SetBody(entry).
		Post(entry.URL(fmt.Sprintf("/upstreams/%s/unregister", entry.Name)))

	if err != nil || resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("server request data error: %v", err)
	}

	return nil
}
