package main

import (
	"os"
	"runtime"
	"server/config"
	"server/utils/logger"
	"strings"
	_"server/service"
)

//server 启动 ./server 环境
//server 启动 ./server PRODUCTION
func main() {

	defer func() {
		if err := recover(); err != nil {
			logger.Error(nil, "sms server start error, restart", err)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())
	arg_num := len(os.Args)
	logger.Info(logger.ExtraFileds{"args": os.Args}, "runtime config is %s", os.Args)
	if arg_num >= 2 {
		config.ENV = strings.ToUpper(os.Args[1])
	} else {
		config.ENV = "DEVELOPMENT"
	}
	logger.Info(nil, "server start at %s", config.ENV)


}
