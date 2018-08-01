package service

import (
	"fmt"
	"github.com/go-resty/resty"
	"time"
)

type Privilege struct {
	Host string
	Port string
}

func (privilege *Privilege) url(route string) string {
	u := fmt.Sprintf("http://%s:%s/api", privilege.Host, privilege.Port)
	return u + route
}

func (privilege *Privilege) request() *resty.Request {
	resty.SetRetryCount(3).SetRetryWaitTime(5 * time.Second).SetRetryMaxWaitTime(20 * time.Second)
	return resty.R()
}

func (privilege *Privilege) WechatMsgSend(errmsg error) {
	privilege.request().
		SetHeader("Host", "sa.wolaidai.com").
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"config": map[string]interface{}{
				"user_list": []string{"freedie.liu", "daniel.dong", "huaijian.wang", "faye.zhou", "andy.zhong"},
				"content":   fmt.Sprintf("[register client] - %s", errmsg),
			},
		}).
		Post(privilege.url("/tasks/wechat/send"))
}
