### 基于redis的消息队列

#### 运行环境

    golang
    redis
    beego

#### 配置

##### app.conf
    
    appname = redis_queue           #项目名称
    httpport = 1200                 #监听端口
    autorender = false              #不渲染模板
    
    include runmode.conf            #导入环境变量
    include config.conf             #数据库等配置
    include queue_list.conf         #队列及消费者配置
    
##### config.conf

    [dev]                           #dev环境下redis配置
    redis_host = 127.0.0.1:6379     #redis host
    redis_password =                #redis password
    
    [prod]                          #prod环境下redis配置
    redis_host = 127.0.0.1:6379     #redis host
    redis_password =                #redis password
    
##### queue_list.conf

    ######### 队列配置
    TEST_QUEUE	= "TEST"                            #topic名称
    TEST_CALLBACK	= "http://xxxx.com/callback"    #消息转发地址
    TEST_REDIS	= "test_queue_redis"                #存在redis中key名
    
    [consumers]                                     #消费者配置
    ######### 消费者配置                             
    TEST = on                                       #topic为TEST的开启消费者

##### runmode.conf

    runmode = dev                   #运行环境
    
##### 消费频率配置（tasks/queue_consumer/consumer.go）

    //注册
    func (t *ConsumerTask) Register () {
    	TaskList := make(map[string]string)
    	TaskList["consumer"]  = "0/2 * * * * *"         //每2秒消费一次
    	for key,value := range TaskList{
    		toolbox.AddTask(key,t.dispatch(key, value))
    	}
    }

#### 运行

    bee run
    
#### 测试

    curl http://127.0.0.1:1200/message/new -X -POST -d "topic=TEST&data1=1&data2=2&data3=3"

    输出：{"code": 200, "data": null, "message": "success"}
    
    然后消费者会向 TEST_CALLBACK（http://xxxx.com/callback）发起一条为POST的请求参数为 data1=1&data2=2&data3=3
    
    curl http://127.0.0.1:1200/message/new?topic=TEST&data1=1&data2=2&data3=3
    
    输出：{"code": 200, "data": null, "message": "success"}
    
    然后消费者会向 TEST_CALLBACK（http://xxxx.com/callback）发起一条为GET的请求参数为 data1=1&data2=2&data3=3
    
#### 生产者接口代码
    
```golang
func (c *IndexCtl) PushQueue() {
    params := c.Input()                                         //获取输入参数
    service.Log("push_queue", params)                           //记录日志
    var topic string                                            //声明topic
    if val, ok := params["topic"]; ok {                         //判断入参是否存在topic
        topic = strings.ToLower(val[0])                         //存在赋值给topic
    } else {
        c.ReturnData(499, "请指定topic", nil)                    //否则抛出错误
    }
    run_topics, _ := beego.AppConfig.GetSection("consumers")    //获取所有topic
    if val, ok := run_topics[topic]; !ok || val != "on" {       //判断当前topic是否存在并且为on状态
        c.ReturnData(499, "topic不存在", nil)                    //不存在或者不为on 抛出错误
    }
    delete(params, "topic")                                     //删除topic参数
    method := c.Ctx.Request.Method                              //获取当前请求类型
    redis_queue := &structs.RedisQueue{                         //结构化
        Method: method,
        Params: params,
    }
    is := new(service.IndexService)                             
    is.LPushRedis(topic, redis_queue)                           //lpush到redis的topic列表中
    c.ReturnData(200, "success", nil)
}
```

#### 消费者代码

```golang
func (t *ConsumerTask) consumer() {
    topics, _ := beego.AppConfig.GetSection("consumers")                                //拿到所有消费者
    for k := range topics {
        queue_name := beego.AppConfig.String(strings.ToUpper(k) + "_QUEUE")             //获取消费者配置信息
        queue_redis := beego.AppConfig.String(strings.ToUpper(k) + "_REDIS")            //获取消费者配置信息
        queue_callback := beego.AppConfig.String(strings.ToUpper(k) + "_CALLBACK")      //获取消费者配置信息
        t.doConsumer(queue_name, queue_redis, queue_callback)                           //开始消费
    }
}

func (t *ConsumerTask) doConsumer(queue_name, queue_redis, queue_callback string)  {
    redis_tool := &service.RedisTool{}                                                  //拿到redis链接
    res, err := redis_tool.LPop(queue_redis)                                            //取出列表中数据
    if res == nil || err != nil {
        return                                                                          //没有数据或者出错中断当前消费
    }
    fmt.Println(queue_name + ": 消费中")
    fmt.Println("消息转发至:" + queue_callback)
    service.Log("do_consumer_log", "queue_name:" + queue_name + " queue_redis:" + queue_redis + " queue_callback:" + queue_callback)
    var r structs.RedisQueue
    err = json.Unmarshal(res.([]byte), &r)
    if err != nil || r.Method == "" {
        return
    }
    t.request(queue_callback, r)                                                        //转发请求
}

func (t *ConsumerTask) request(uri string, param structs.RedisQueue) {
    var req *httplib.BeegoHTTPRequest
    switch strings.ToUpper(param.Method) {                                              //请求方式
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
    for k, v := range param.Params {                                                    //数据
        if len(v) < 1 {
            continue
        }
        req.Param(k, v[0])
    }
    _, err := req.Bytes()                                                               //发送请求
    if err != nil {
        service.Log("queue_request_err", "url:" + uri + " error:" + err.Error() + " params:" + fmt.Sprint(param.Params))
        fmt.Println(err.Error())
    }
    fmt.Println("消费完成")
}
```