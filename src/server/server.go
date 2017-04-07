package main

import (
	"os"
	"runtime"
	"server/config"
	"server/utils/logger"
	"strings"
	"net/http"
	"server/app/api"
	"github.com/gorilla/mux"
	"log"
	"time"
	"server/service/redis"
	"strconv"
)

//server 启动 ./server 环境
//server 启动 ./server PRODUCTION httpaddress
func main() {

	defer func() {
		if err := recover(); err != nil {
			logger.Error(nil, "sms server start error", err)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())
	arg_num := len(os.Args)
	logger.Info(logger.ExtraFileds{"args": os.Args}, "runtime config is %s", os.Args)
	httpaddress := "127.0.0.1:9090"
	if arg_num >= 2 {
		config.ENV = strings.ToUpper(os.Args[1])
		if arg_num > 2 {
			if "" != os.Args[2] {
				httpaddress = os.Args[2]
			}
		}
	} else {
		config.ENV = "DEVELOPMENT"
	}
	logger.Info(nil, "server start at %s", config.ENV)
	redis.REDIS_ADDRESS = config.CONFIG[config.ENV]["REDIS_ADDRESS"]
	redis.REDIS_PASSWORD = config.CONFIG[config.ENV]["REDIS_PASSWORD"]
	redis.REDIS_MAX_POOL_SIZE, _ = strconv.Atoi(config.CONFIG[config.ENV]["MAX_POOL_SIZE"])
	httpServer(httpaddress) //阻塞监听
	logger.Info(nil,"api http server shutdown","")
}

func httpServer(httpaddress string){
	defer func() {
		if err := recover(); err != nil {
			logger.Error(nil, "api http server error %s ", err)
		}
	}()
	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/", api.Index)
	s.HandleFunc("/1/{channel}/captcha/token", api.GetImageCaptchaToken) //图形验证码token
	s.HandleFunc("/1/{channel}/captcha", api.GetImageCaptcha) //图形验证码
	s.HandleFunc("/1/{channel}/captcha/sms/send", nil) //TODO 发送短信验证码
	s.HandleFunc("/1/{channel}/captcha/sms/verify", nil) //TODO 验证码验证
	logger.Info(nil,"Server is at %s", httpaddress)
	srv := &http.Server{
		Handler:      r,
		Addr:         httpaddress,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 45 * time.Second,
		ReadTimeout:  45 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
