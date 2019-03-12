package main

import (
	_ "redis_queue/routers"
	"github.com/astaxie/beego"
	tasks2 "redis_queue/tasks"
)

func main() {
	tasks := new(tasks2.TaskService)
	tasks.StartTasks()
	beego.Run()
}

