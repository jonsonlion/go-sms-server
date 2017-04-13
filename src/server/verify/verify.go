package verify

import (
	"server/utils/uuid"
	"server/service/redis"
	"fmt"
	"server/utils/logger"
	"server/utils/random"
	"encoding/json"
	"errors"
	"server/verify/captcha"
	"server/utils/convert"
	"server/verify/store"
	"server/app/domain"
	"server/service/sms"
	"server/service/mongo/mongointercept"
	"time"
)

var (
	REDIS_KEY_IMAGE_CAPTCHA_TOKEN_EXPIRE = 15  //15s
	REDIS_KEY_SMS_CAPTCHA_EXPIRE = 1800        //30 min
	REDIS_KEY_IMAGE_CAPTCHA_TOKEN = "sms:key:image:captcha:token:%s"  //图形验证码缓存
	REDIS_KEY_SMS_CAPTCHA = "sms:key:sms:captcha:%s"   //短信验证码缓存
	REDIS_KEY_SMS_INTERCEPT_IP = "sms:key:sms:intercept:ip:%s"   //ip拦截
	REDIS_KEY_SMS_INTERCEPT_IP_EXPIRE = 3600   //ip拦截1小时
	REDIS_KEY_SMS_INTERCEPT_IP_LIMIT = 10   //ip拦截 次数
	REDIS_KEY_SMS_INTERCEPT_1MINUTE_PHONE = "sms:key:sms:intercept:1minute:phone:%s"   //手机号拦截
	REDIS_KEY_SMS_INTERCEPT_PHONE = "sms:key:sms:intercept:phone:%s"   //手机号拦截
	REDIS_KEY_SMS_INTERCEPT_PHONE_EXPIRE = 1800  //手机号拦截
	REDIS_KEY_SMS_INTERCEPT_PHONE_LIMIT= 5  //手机号拦截 5个

	REDIS_KEY_SMS_INTERCEPT_MAN_MADE= "sms:key:sms:intercept:manmade:phone:%s"  //人工认证的验证码
)

//生成流水token
func GenerateFlowToken()string{
	return uuid.GetUUIDWithoutLine()
}

//存储token，记录获取图片token的有效性
func StoreImageCaptchaToken(token string, value string){
	key := fmt.Sprintf(REDIS_KEY_IMAGE_CAPTCHA_TOKEN, token)
	redis.SetAndExpire(key , value, REDIS_KEY_IMAGE_CAPTCHA_TOKEN_EXPIRE)
}

//存储图片验证码token验证
func ValidateImageCodeToken(token string)(bool,error){
	key := fmt.Sprintf(REDIS_KEY_IMAGE_CAPTCHA_TOKEN, token)
	value := redis.GetString(key)
	if "" == value{
		return false, errors.New(fmt.Sprintf("img token invalidate, token:%s",token))
	}
	redis.DEL(key)//验证完后清理掉
	return true,nil
}

//发送短信验证码
// signature string 消息签名
// token string 图片验证码token
// imgcode string 图片验证码值
// return []byte 返回 短信验证码
func SendSms(signature string,channel string, token string, imgcode string, phone string, len int, ip string)(string, error, int){
	//拦截验证
	c, e := interceptFilter(signature, channel, token, imgcode, phone, ip)
	if nil != e {
		return "", e, c
	}
	code := random.RandomNum(len)
	_, err := send(phone,code,signature,channel)
	if nil != err{
		return "", errors.New(fmt.Sprintf("send sms error result: %s" , err)), domain.ERROR
	}
	//添加到缓存
	smstoken := GenerateFlowToken() + "@" + convert.BytesToString(random.RandomNum(5))
	smskey := fmt.Sprintf(REDIS_KEY_SMS_CAPTCHA, smstoken)
	smsvalue := map[string]string{
		"phone" : phone,
		"captcha" : convert.BytesToString(code),
	}
	bv, _ := json.Marshal(smsvalue)
	redis.SetAndExpire(smskey, string(bv), REDIS_KEY_SMS_CAPTCHA_EXPIRE)
	//清除图形验证码缓存
	store.Clear(token)
	return smstoken,nil,domain.SUCCESS
}

//验证短信验证码是否正确
//phone string 手机号
//smstoken发送短信后会生成一个token,验证的时候使用
//captcha string 短信验证码
func ValidateSmsCaptcha(phone string, smstoken string, captcha string)bool{
	//号码和验证码是否经过人工认证
	manmadekey := fmt.Sprintf(REDIS_KEY_SMS_INTERCEPT_MAN_MADE, phone)
	mdv := redis.GetString(manmadekey)
	if "" != mdv {
		if mdv == captcha{
			return true
		}
	}
	smskey := fmt.Sprintf(REDIS_KEY_SMS_CAPTCHA, smstoken)
	smsvalue := map[string]string{}
	value := redis.GetString(smskey)
	json.Unmarshal([]byte(value), &smsvalue)
	logger.Info(nil, "validate sms captcha cachevalue:%s, phone:%s, smstoken:%s, captcha:%s", smsvalue, phone, smstoken, captcha)
	return phone == smsvalue["phone"] && captcha == smsvalue["captcha"]
}

//设置人工设置的验证码
func ManMadeSmsCaptcha(phone string, captcha string)error{
	//添加人工认证的验证码
	manmadekey := fmt.Sprintf(REDIS_KEY_SMS_INTERCEPT_MAN_MADE, phone)
	_, err := redis.SetAndExpire(manmadekey, captcha,REDIS_KEY_SMS_CAPTCHA_EXPIRE)
	return err
}

//拦截验证
func interceptFilter(signature string,channel string, token string, imgcode string, phone string, ip string) (int, error){
	item := mongointercept.Item{}
	item.No = uuid.GetUUIDWithoutLine()
	item.Phone = phone
	item.Channel = channel
	item.Signature = signature
	item.Reason = ""
	item.CreateDate = time.Now()
	// 验证图形验证码是否正确
	if !captcha.Verify(token,convert.StringToBytes(imgcode)){
		return domain.INVALID_IMG_CAPTCHA, errors.New("图形验证码错误")
	}
	//ip次数验证
	ipkey := fmt.Sprintf(REDIS_KEY_SMS_INTERCEPT_IP, ip)
	ripcount,err := redis.IncrAndExpire(ipkey, REDIS_KEY_SMS_INTERCEPT_IP_EXPIRE)
	if nil != err {
		logger.Error(nil, "interceptFilter error ipKey:%s ", ipkey)
	}
	if ripcount > REDIS_KEY_SMS_INTERCEPT_IP_LIMIT{
		//记录限制
		item.Reason = fmt.Sprintf("区域发送超过限制(IP发送限制limit:%s)", REDIS_KEY_SMS_INTERCEPT_IP_LIMIT)
		mongointercept.SaveMessage(&item)
		return domain.SMS_INTERCEPT_IP_LIMIT , errors.New("区域发送超过限制")
	}
	//手机1分钟内发送次数限制
	phone1MinuteKey := fmt.Sprintf(REDIS_KEY_SMS_INTERCEPT_1MINUTE_PHONE, phone)
	phone1MinuteKeycount,err := redis.IncrAndExpire(phone1MinuteKey, 60)  //一分钟只能发送一个
	if nil != err {
		logger.Error(nil, "interceptFilter error phone1MinuteKey:%s ", phone1MinuteKey)
	}
	if phone1MinuteKeycount > 1{
		//记录限制
		item.Reason = fmt.Sprintf("验证码发送过于频繁(手机1分钟内发送次数限制limit:%s)", 1)
		mongointercept.SaveMessage(&item)
		return domain.SMS_INTERCEPT_1MINUTE_PHONE_LIMIT , errors.New("验证码发送过于频繁")
	}
	//手机发送次数限制
	phone1HourKey := fmt.Sprintf(REDIS_KEY_SMS_INTERCEPT_PHONE, phone)
	phone1HourKeycount, err := redis.IncrAndExpire(phone1HourKey, REDIS_KEY_SMS_INTERCEPT_PHONE_EXPIRE)  //一小时
	if nil != err {
		logger.Error(nil, "interceptFilter error phone1HourKey:%s ", phone1HourKey)
	}
	if phone1HourKeycount > REDIS_KEY_SMS_INTERCEPT_PHONE_LIMIT{
		//记录限制
		item.Reason = fmt.Sprintf("验证码发送超过限制(手机发送次数限制limit:%s)", REDIS_KEY_SMS_INTERCEPT_PHONE_LIMIT)
		mongointercept.SaveMessage(&item)
		return domain.SMS_INTERCEPT_PHONE_LIMIT , errors.New("验证码发送超过限制")
	}
	return 0, nil
}

func send(phone string, verify []byte, signature string, channel string)(string,error){
	detail := fmt.Sprintf("感谢您的使用，您的验证码是 %s ,请在%d分钟内使用，谢谢！", convert.BytesToString(verify), REDIS_KEY_SMS_CAPTCHA_EXPIRE/60)
	return sms.SendSms(phone,channel,signature,detail)
}