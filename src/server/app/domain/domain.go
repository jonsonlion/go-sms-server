package domain

const(
	SUCCESS = 1  //操作成功
	ERROR = -1  //系统异常
	INVALID_IMG_CAPTCHA = -2 // 图形验证码错误
	INVALID_SMS_CAPTCHA = -3 // 短信验证码错误
	SMS_INTERCEPT_IP_LIMIT = -5 // 区域IP发送超过限制
	SMS_INTERCEPT_1MINUTE_PHONE_LIMIT = -6 // 验证码发送过于频繁
	SMS_INTERCEPT_PHONE_LIMIT = -7 // 验证码获取超过限制
)

type BaseResult struct {
	Code int 	`json:"code"`
	Msg string	`json:"msg"`
	Reason string	`json:"reason"`
}

type CaptchaToken struct {
	BaseResult
	Token string	`json:"token"`
}

type SmsCaptcha struct {
	BaseResult
	Token string `json:"token"`
}

type ManMadeSmsCaptcha struct {
	BaseResult
	Captcha string `json:"captcha"`
}
