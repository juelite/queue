package controllers

import "github.com/astaxie/beego"

type BaseCtx struct {
	beego.Controller
}

func (this *BaseCtx) ReturnData(code int,  message string, data map[string]interface{}) {

	this.Data["json"] = map[string]interface{}{
		"code"    : code ,
		"message" : message ,
		"data"    : data ,
	}
	this.ServeJSON()
	this.StopRun()
}