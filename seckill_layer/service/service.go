package service

import "github.com/astaxie/beego/logs"

func Run() (err error) {
	logs.Debug("start run")
	err = RunProcess()
	if err != nil {
		logs.Error("RunProcess failed!!")
	}
	return
}
