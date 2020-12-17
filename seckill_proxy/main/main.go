package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "seckill_proxy/router"
)

func main() {
	err := InitConfig()
	if err != nil {
		panic(err)
		return
	}

	logs.Debug(" init sec")
	err = InitSec()
	if err != nil {
		panic(err)
		return
	}
	beego.Run()
}
