package config

var ENV string = "DEVELOPMENT"

var CONFIG = map[string]map[string]string{
	"DEVELOPMENT": {
		"AMQP_ADDRESS":    "amqp://admin:qpwoeiruty@ip:5672/message",
		"AMQP_QUEUE":      "smsexchange-dev",
		"MONGO_ADDRESS":   "ip:27017",
		"REDIS_ADDRESS":   "ip:6379",
		"MAX_POOL_SIZE":   "5",
		"REDIS_PASSWORD":  "111111",
	},
	"TEST": {

	},
	"PRODUCTION": {

	},
}
