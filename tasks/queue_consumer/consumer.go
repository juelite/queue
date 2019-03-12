package queue_consumer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/toolbox"
	"redis_queue/service"
	"redis_queue/structs"
)

type ConsumerTask struct {}

//注册
func (t *ConsumerTask) Register () {
	TaskList := make(map[string]string)
	TaskList["consumer"]  = "0/2 * * * * *"
	for key,value := range TaskList{
		toolbox.AddTask(key,t.dispatch(key, value))
	}
}

//模板
func (t *ConsumerTask) dispatch(tname string,spec string) *toolbox.Task {
	tasks := toolbox.NewTask(
		tname,spec, func() error {
			err := t.doDispatch(tname)
			return err
		},
	)
	return tasks
}

/**
 * 根据任务类型通知
 * @param tp string 任务类型
 */
func (t *ConsumerTask) doDispatch(tp string) error {
	switch tp {
	case "consumer":
		t.consumer()
	}
	return nil
}

func (t *ConsumerTask) consumer() {
	topics, _ := beego.AppConfig.GetSection("consumers")
	for k := range topics {
		queue_name := beego.AppConfig.String(strings.ToUpper(k) + "_QUEUE")
		queue_redis := beego.AppConfig.String(strings.ToUpper(k) + "_REDIS")
		queue_callback := beego.AppConfig.String(strings.ToUpper(k) + "_CALLBACK")
		t.doConsumer(queue_name, queue_redis, queue_callback)
	}
}

func (t *ConsumerTask) doConsumer(queue_name, queue_redis, queue_callback string)  {
	redis_tool := &service.RedisTool{}
	res, err := redis_tool.LPop(queue_redis)
	if res == nil || err != nil {
		return
	}
	fmt.Println(queue_name + ": 消费中")
	fmt.Println("消息转发至:" + queue_callback)
	service.Log("do_consumer_log", "queue_name:" + queue_name + " queue_redis:" + queue_redis + " queue_callback:" + queue_callback)
	var r structs.RedisQueue
	err = json.Unmarshal(res.([]byte), &r)
	if err != nil || r.Method == "" {
		return
	}
	t.request(queue_callback, r)
}

func (t *ConsumerTask) request(uri string, param structs.RedisQueue) {
	var req *httplib.BeegoHTTPRequest
	switch strings.ToUpper(param.Method) {
	case "POST":
		req = httplib.Post(uri)
	case "GET":
		req = httplib.Get(uri)
	case "PUT":
		req = httplib.Put(uri)
	case "DELETE":
		req = httplib.Delete(uri)
	case "HEAD":
		req = httplib.Head(uri)
	default:
		return
	}
	for k, v := range param.Params {
		if len(v) < 1 {
			continue
		}
		req.Param(k, v[0])
	}
	_, err := req.Bytes()
	if err != nil {
		service.Log("queue_request_err", "url:" + uri + " error:" + err.Error() + " params:" + fmt.Sprint(param.Params))
		fmt.Println(err.Error())
	}
	fmt.Println("消费完成")
}