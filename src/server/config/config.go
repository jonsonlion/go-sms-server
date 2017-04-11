package config

var ENV string = "DEVELOPMENT"

var CONFIG = map[string]map[string]string{
	"DEVELOPMENT": {
		"AMQP_ADDRESS":    "amqp://admin:qpwoeiruty@192.168.19.141:5672/message",
		"AMQP_QUEUE":      "smsexchange-dev",
		"MONGO_ADDRESS":   "192.168.18.88:27017",
		"REDIS_ADDRESS":   "192.168.18.102:6379",
		"MAX_POOL_SIZE":   "5",
		"REDIS_PASSWORD":  "111111",
	},
	"TEST": {

	},
	"PRODUCTION": {

	},
}
