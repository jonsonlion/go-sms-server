package verify

import (
	"server/utils/uuid"
	"server/service/redis"
	"fmt"
	"server/utils/logger"
	"server/utils/random"
	"encoding/json"
	"errors"
	"github.com/dchest/captcha"
	"server/utils/convert"
)

var (
	REDIS_KEY_IMAGE_CAPTCHA_TOKEN_EXPIRE = 15  //15s
	REDIS_KEY_SMS_CAPTCHA_EXPIRE = 1800        //30 min
	REDIS_KEY_IMAGE_CAPTCHA_TOKEN = "sms:key:image:captcha:token:%s"
	REDIS_KEY_SMS_CAPTCHA = "sms:key:sms:captcha:%s"
)

//生成流水token
func GenerateFlowToken()string{
	return uuid.GetUUIDWithoutLine()
}

//存储token，记录获取图片token的有效性
func StoreImageCaptchaToken(token string, value string){
	key := fmt.Sprintf(REDIS_KEY_IMAGE_CAPTCHA_TOKEN, token)
	redis.Set(key , value)
	redis.EXPIRE(key, REDIS_KEY_IMAGE_CAPTCHA_TOKEN_EXPIRE)
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
// token string 图片验证码token
// imgcode string 图片验证码值
// return []byte 返回 短信验证码
func SendSms(token string, imgcode string, phone string, len int)(string, error){
	if !captcha.Verify(token,convert.StringToBytes(imgcode)){
		return "", errors.New( fmt.Sprintf("img captcha invalidate, token:%s",token))
	}
	code := random.RandomNum(len)
	send(phone,code)
	//添加到缓存
	smstoken := GenerateFlowToken()
	smskey := fmt.Sprintf(REDIS_KEY_SMS_CAPTCHA, smstoken)
	smsvalue := map[string]string{
		"phone" : phone,
		"captcha" : string(code),
	}
	bv, _ := json.Marshal(smsvalue)
	redis.Set(smskey , string(bv))
	redis.EXPIRE(smskey, REDIS_KEY_SMS_CAPTCHA_EXPIRE)
	return smstoken,nil
}

//验证短信验证码是否正确
//phone string 手机号
//smstoken发送短信后会生成一个token,验证的时候使用
//captcha string 短信验证码
func ValidateSmsCaptcha(phone string, smstoken string, captcha string)bool{
	smskey := fmt.Sprintf(REDIS_KEY_SMS_CAPTCHA, smstoken)
	smsvalue := map[string]string{}
	value := redis.GetString(smskey)
	json.Unmarshal([]byte(value), &smsvalue)
	logger.Info(nil, "validate sms captcha cachevalue:%s, phone:%s, smstoken:%s, captcha:%s", smsvalue, phone, smstoken, captcha)
	return phone != smsvalue["phone"] && captcha == smsvalue["captcha"]
}

func send(phone string, verify []byte)string{
	detail := fmt.Sprintf("感谢您的使用，您的验证码是 %s ,请在%d分钟内使用，谢谢！", string(verify), REDIS_KEY_SMS_CAPTCHA_EXPIRE/60)
	logger.Info(nil,"send sms detail: %s", detail)
	return "1"
}