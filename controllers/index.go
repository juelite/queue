package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"redis_queue/service"
	"redis_queue/structs"
)

type IndexCtl struct {
	BaseCtx
}

func (c *IndexCtl) PushQueue() {
	params := c.Input()
	service.Log("push_queue", params)
	var topic string
	if val, ok := params["topic"]; ok {
		topic = strings.ToLower(val[0])
	} else {
		c.ReturnData(499, "请指定topic", nil)
	}
	run_topics, _ := beego.AppConfig.GetSection("consumers")
	if val, ok := run_topics[topic]; !ok || val != "on" {
		c.ReturnData(499, "topic不存在", nil)
	}
	delete(params, "topic")
	method := c.Ctx.Request.Method
	redis_queue := &structs.RedisQueue{
		Method: method,
		Params: params,
	}
	is := new(service.IndexService)
	is.LPushRedis(topic, redis_queue)
	c.ReturnData(200, "success", nil)
}
