package sms

import (
	"server/utils/logger"
	"encoding/json"
	"server/service/mongo"
	"time"
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
	item := mongo.Item{}
	item.Phone = phone
	item.Channel = channel
	item.Signature = signature
	item.Message = message
	item.Arrived = 0
	item.CreateDate = time.Now()
	flowId , err := mongo.SaveMessage( &item )
	//调用 短信接口
	r := send(phone, signature, message)
	if err == nil{
		m := mongo.ItemResponse{}
		m.FlowId = flowId
		m.Arrived = 1
		m.Reason = r
		err := mongo.UpdateMessage( &m )
		if err != nil{
			logger.Info(nil,"req-- send sms phone:%s, channel:%s, signature:%s, message: %s, error:%s",phone, channel, signature, message, err)
		}
		logger.Info(nil,"res-- send sms phone:%s, channel:%s, signature:%s, message: %s, result:%s",phone, channel, signature, message, r)
	}
	return "1", nil
}

func send( phone string , signature string,  message string)string{
	return "1"
}
