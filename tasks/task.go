package tasks

import (
	"fmt"

	"github.com/astaxie/beego/toolbox"
	"redis_queue/tasks/queue_consumer"
)

type TaskService struct {}


//初始化
func init()  {
	//注册消费者
	fmt.Println("register queue consumers tasks")
	consumer := new(queue_consumer.ConsumerTask)
	consumer.Register()
}

//开始
func (t *TaskService) StartTasks()  {
	toolbox.StartTask()
}