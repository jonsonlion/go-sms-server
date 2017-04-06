package config

var ENV string = "DEVELOPMENT"

var CONFIG = map[string]map[string]string{
	"DEVELOPMENT": {
		"APNS_CERT":       "jxyr-push-prod-Certificates.p12",
		"APNS_ENV":        "prod",
		"APNS_PASSWORD":   "1",
		"ANNS_TOPIC":      "com.jxyr.ygmobile",
		"AMQP_ADDRESS":    "amqp://admin:qpwoeiruty@192.168.19.141:5672/message",
		"AMQP_QUEUE":      "pushqueue-dev",
		"MONGO_ADDRESS":   "192.168.18.88:27017",
		"REDIS_ADDRESS":   "192.168.18.102:6379",
		"MAX_POOL_SIZE":   "5",
		"REDIS_PASSWORD":  "111111",
		"BAD_TOKEN_QUEUE": "push:ios:baddevicetoken:queue:com.jxyr.ygmobile",
	},
	"TEST": {
		"APNS_CERT":       "jxyr-push-dev-Certificates.p12",
		"APNS_ENV":        "dev",
		"APNS_PASSWORD":   "1",
		"ANNS_TOPIC":      "com.jxyr.ygmobile",
		"AMQP_ADDRESS":    "amqp://admin:qpwoeiruty@192.168.19.141:5672/message",
		"AMQP_QUEUE":      "pushqueue-test",
		"MONGO_ADDRESS":   "192.168.18.88:27017",
		"REDIS_ADDRESS":   "192.168.18.102:6379",
		"MAX_POOL_SIZE":   "5",
		"REDIS_PASSWORD":  "111111",
		"BAD_TOKEN_QUEUE": "push:ios:baddevicetoken:queue:com.jxyr.ygmobile",
	},
	"PRODUCTION": {
		"APNS_CERT":       "jxyr-push-prod-Certificates.p12",
		"APNS_ENV":        "prod",
		"APNS_PASSWORD":   "1",
		"ANNS_TOPIC":      "com.jxyr.ygmobile",
		"AMQP_ADDRESS":    "amqp://admin:qpwoeiruty@192.168.19.141:5672/message",
		"AMQP_QUEUE":      "pushqueue",
		"MONGO_ADDRESS":   "192.168.18.88:27017",
		"REDIS_ADDRESS":   "192.168.18.102:6379",
		"MAX_POOL_SIZE":   "5",
		"REDIS_PASSWORD":  "111111",
		"BAD_TOKEN_QUEUE": "push:ios:baddevicetoken:queue:com.jxyr.ygmobile",
	},
}
