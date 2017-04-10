package domain

const(
	SUCCESS = 1  //操作成功
	ERROR = -1  //系统异常
	INVALID_IMG_CAPTCHA = -2 // 图形验证码错误
	INVALID_SMS_CAPTCHA = -3 // 短信验证码错误
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
