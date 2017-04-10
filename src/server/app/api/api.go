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

const (
	errorStr = "{\"code\":-1,\"msg\":\"处理失败\",\"reason\":\"%s\"}"
	success = "{\"code\":1,\"msg\":\"处理成功\"}"
)
const (
	imgCaptchaLen = 4
	smsCaptchaLen = 6
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
	result.Code = domain.SUCCESS
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
	channel := getPathParam("channel",req)
	token := req.FormValue("token")
	logger.Info(nil, "GetImageCaptcha channel:%s, token:%s", channel, token)
	width := 200
	height := 80
	if "" != req.FormValue("width"){
		width,_ = strconv.Atoi(req.FormValue("width"))
	}
	if "" != req.FormValue("height"){
		height,_ = strconv.Atoi(req.FormValue("height"))
	}
	b,err := verify.ValidateImageCodeToken(token)
	if b{
		id := captcha.NewLenCustom(imgCaptchaLen, token)
		err = captcha.WriteImage(w,id,width,height)
	}
	if nil != err{
		logger.Error(nil,"GetImageCaptcha captcha.NewImage error %s",err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(errorStr,err)))
	}
}

//发送短信验证码
func SendSmsCaptcha(w http.ResponseWriter, req *http.Request){
	channel := getPathParam("channel",req)
	token := req.FormValue("token")
	phone := req.FormValue("phone")
	imgcaptcha := req.FormValue("imgcaptcha")
	logger.Info(nil, "SendSmsCaptcha channel: %s, phone:%s, token:%s, imgcaptcha:%s",channel,phone, token, imgcaptcha)
	setCommonHeaders(&w)
	w.Header().Set("Content-Type", "application/json")

	result := domain.SmsCaptcha{}
	result.Code = domain.SUCCESS
	result.Msg = "success"
	smsToken , err, code := verify.SendSms(token, imgcaptcha, phone, smsCaptchaLen)
	if nil != err{
		result.Code = code
		result.Msg = "发送验证码失败"
		result.Reason = fmt.Sprintf("%s", err)
		logger.Error(nil, "SendSmsCaptcha error channel: %s, phone:%s, token:%s, imgcaptcha:%s, error:%s",channel,phone, token, imgcaptcha,err)
	}else{
		result.Token = smsToken
		verify.StoreImageCaptchaToken(result.Token,"TODO")//存储什么内容？
	}
	b, err := json.Marshal(result)
	if nil == err{
		w.Write(b)
	}else{
		w.Write([]byte(errorStr))
	}
}

//发送短信验证码
func ValidateSmsCaptcha(w http.ResponseWriter, req *http.Request){
	channel := getPathParam("channel",req)
	token := req.FormValue("token")
	phone := req.FormValue("phone")
	captcha := req.FormValue("captcha")
	logger.Info(nil, "ValidateSmsCaptcha channel: %s, phone:%s, token:%s, captcha:%s",channel,phone, token, captcha)
	setCommonHeaders(&w)
	w.Header().Set("Content-Type", "application/json")

	result := domain.SmsCaptcha{}
	result.Code = domain.SUCCESS
	result.Msg = "success"
	r := verify.ValidateSmsCaptcha(phone, token, captcha)
	if r{
		result.Code = domain.SUCCESS
		result.Msg = "验证成功"
	}else{
		result.Code = domain.INVALID_SMS_CAPTCHA
		result.Msg = "短信验证码错误"
		logger.Info(nil, "SendSmsCaptcha failed channel: %s, phone:%s, token:%s, captcha:%s",channel,phone, token, captcha)
	}
	b, err := json.Marshal(result)
	if nil == err{
		w.Write(b)
	}else{
		w.Write([]byte(errorStr))
	}
}


























