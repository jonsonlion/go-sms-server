package api

import (
	"net/http"
	"server/app/domain"
	"encoding/json"
	"server/utils/logger"
	"server/verify/captcha"
	"server/verify/store"
	"time"
	"github.com/gorilla/mux"
	"server/verify"
	"strconv"
	"fmt"
)

var (
	errorStr = "{\"code\":-1,\"msg\":\"处理失败%s\"}"
	success = "{\"code\":1,\"msg\":\"处理成功\"}"
)

func init(){
	captcha.SetCustomStore( store.NewRedisStore(5 * time.Minute) ) //缓存5分钟，图形验证码
}

/**
 * 设置默认的header
 */
func setCommonHeaders(w *http.ResponseWriter)(http.ResponseWriter){
	(*w).Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	(*w).Header().Set("Pragma", "no-cache")
	(*w).Header().Set("Expires", "0")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	return *w
}

//获取pathparam
func getPathParam(name string,req *http.Request)interface{}{
	vars := mux.Vars(req)
	return vars[name]
}

//默认首页
func Index(w http.ResponseWriter, req *http.Request) {
	setCommonHeaders(&w)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(success))
}

//get
//获取图形验证码token
func GetImageCaptchaToken(w http.ResponseWriter, req *http.Request){
	channel := getPathParam("channel",req)
	result := domain.CaptchaToken{}
	result.Code = 1
	result.Msg = "success"
	result.Token = verify.GenerateFlowToken()
	logger.Info(nil, "GetImageCaptchaToken channel: %s, token:%s",channel, result.Token)
	setCommonHeaders(&w)
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(result)
	if nil == err{
		verify.StoreImageCaptchaToken(result.Token,"TODO")//存储什么内容？
		w.Write(b)
	}else{
		w.Write([]byte(errorStr))
	}
}

//get or post
//根据图形验证码token，获取图形验证码
func GetImageCaptcha(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "image/png")
	token := req.FormValue("token")
	logger.Info(nil, "image token %s", token)
	width := 200
	hight := 80
	if "" != req.FormValue("width"){
		width,_ = strconv.Atoi(req.FormValue("width"))
	}
	if "" != req.FormValue("hight"){
		hight,_ = strconv.Atoi(req.FormValue("hight"))
	}
	b,err := verify.ValidateImageCodeToken(token)
	if b{
		id := captcha.NewLenCustom(4, token)
		err = captcha.WriteImage(w,id,width,hight)
	}
	if nil != err{
		logger.Error(nil,"GetImageCaptcha captcha.NewImage error %s",err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(errorStr,err)))
	}
}

//发送短信验证码
func SendSmsCaptcha(w http.ResponseWriter, req *http.Request){

}

//发送短信验证码
func ValidateSmsCaptcha(w http.ResponseWriter, req *http.Request){

}


























