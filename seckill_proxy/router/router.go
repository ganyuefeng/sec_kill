package router

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"seckill_proxy/controller"
)

func init() {
	logs.Debug("enter router init")
	beego.Router("/seckill", &controller.SkillController{}, "*:SecKill")
	beego.Router("/secinfo", &controller.SkillController{}, "*:SecInfo")
}
