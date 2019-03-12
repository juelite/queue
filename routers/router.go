package routers

import (
	"redis_queue/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/message/new", &controllers.IndexCtl{}, "*:PushQueue")
}
