package sms

import (
	"server/utils/logger"
	"encoding/json"
)

type SmsMq struct {
	Phone string `json:"phone"`
	Channel string `json:"channel"`
	Signature string `json:"signature"`
	Message string `json:"message"`
}

//短信mq消息消费者
//{"phone":"11111111", "channel":"", "signature":"yry", "message":"欢迎欢迎"}
func SmsMqConsumer(message chan []byte){
	for msg := range message{
		sm := SmsMq{}
		err := json.Unmarshal(msg, &sm)
		if nil == err{
			SendSms(sm.Phone, sm.Channel, sm.Signature, sm.Message)
		}
	}
}

//发送消息
func SendSms(phone string, channel string, signature string, message string) (string ,error){
	logger.Info(nil,"send sms phone:%s, channel:%s, signature:%s, message: %s",phone, channel, signature, message)
	//TODO send
	logger.Info(nil,"send sms phone:%s, channel:%s, signature:%s, message: %s, result:%s",phone, channel, signature, message, "1")
	return "1", nil
}
