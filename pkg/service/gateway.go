package service

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/srelab/register/pkg/g"
	"net/http"
	"time"
)

type Gateway struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (gateway *Gateway) url(route string) string {
	u := fmt.Sprintf("http://%s:%s", g.Config().Gateway.Host, g.Config().Gateway.Port)
	return u + route
}

func (gateway *Gateway) request() *resty.Request {
	resty.SetRetryCount(3).SetRetryWaitTime(5 * time.Second).SetRetryMaxWaitTime(20 * time.Second)
	return resty.R()
}

func (gateway *Gateway) Register() error {
	resp, err := gateway.request().
		SetHeader("Content-Type", "application/json").
		SetBody(gateway).
		Post(gateway.url(fmt.Sprintf("/upstreams/%s/register", gateway.Name)))

	if err != nil || resp.StatusCode() != http.StatusOK {
		if err == nil {
			err = errors.New(string(resp.Body()))
		}

		return err
	}

	return nil
}

func (gateway *Gateway) UnRegister() error {
	resp, err := gateway.request().
		SetHeader("Content-Type", "application/json").
		SetBody(gateway).
		Post(gateway.url(fmt.Sprintf("/upstreams/%s/unregister", gateway.Name)))

	if err != nil || resp.StatusCode() != http.StatusOK {
		if err == nil {
			err = errors.New(string(resp.Body()))
		}

		return err
	}

	return nil
}
