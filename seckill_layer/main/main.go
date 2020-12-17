package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"seckill_layer/service"
)

func main() {
	// load config
	err := initConfig("ini", "./conf/seclayer.conf")
	if err != nil {
		logs.Error("init config failed", err)
		panic(fmt.Sprintf("init config failed, err:%v", err))
	}
	logs.Debug("init config success")
	//init logger
	/*
		err = initLogger()
		if err != nil {
			logs.Error("initLogger failed", err)
			panic(fmt.Sprintf("initLogger failed, err:%v", err))
		}
		logs.Debug("initLogger success")

	*/
	//init seckill_layer
	err = service.InitSecLayer(appConfig)
	if err != nil {
		logs.Error("initSecKill failed", err)
		panic(fmt.Sprintf("initSecKill failed, err:%v", err))
	}
	logs.Debug("InitSecLayer success")
	//service start
	err = service.Run()
	if err != nil {
		logs.Error("serviceRun failed", err)
		panic(fmt.Sprintf("serviceRun failed, err:%v", err))
	}

	logs.Info("service start success!")
}
