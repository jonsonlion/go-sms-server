package domain

type BaseResult struct {
	Code int 	`json:"code"`
	Msg string	`json:"msg"`
}

type CaptchaToken struct {
	BaseResult
	Token string	`json:"token"`
}
