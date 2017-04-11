package amqp

import (
	"github.com/streadway/amqp"
	"server/config"
	"server/utils/logger"
)

func StartAmqpConsumer(message chan []byte) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(nil, "StartAmqpConsumer error", err)
		}
	}()
	conn, err := amqp.Dial(config.CONFIG[config.ENV]["AMQP_ADDRESS"])
	defer conn.Close()
	if err != nil {
		logger.Error(nil, "amqp start error %s", err)
		panic("amqp start error")
	}
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(nil, "amqp fail to open a channel %s", err)
		panic("amqp fail to open a channel")
	}
	defer ch.Close()
	msgs, err := ch.Consume(
		config.CONFIG[config.ENV]["AMQP_QUEUE"], // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logger.Error(nil, "consumer listen error", err)
	}
	for d := range msgs {
		message <- d.Body
	}
}
