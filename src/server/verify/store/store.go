package store

import (
	"time"
	"github.com/dchest/captcha"
	"server/utils/logger"
	"server/service/redis"
	"fmt"
	"server/utils/convert"
)

var (
	REDIS_KEY_STORE_ID = "sms:key:store:id:%s"
)

type redisStore struct {
	expiration time.Duration
}

//NewRedisStore returns a new store for captchas
func NewRedisStore(expiration time.Duration) captcha.Store{
	store := new(redisStore)
	store.expiration = expiration
	return store
}

//缓存验证码
func (store *redisStore) Set(id string, digits []byte){
	logger.Info(nil,"set captcha to redis id:%s, digits:%s", id, convert.BytesToString(digits))
	key := fmt.Sprintf(REDIS_KEY_STORE_ID, id)
	value := convert.BytesToString(digits)
	redis.Set(key, value)
	redis.EXPIRE(key, int(store.expiration.Seconds()))
}

//获取验证码
func (store *redisStore) Get(id string, clear bool)([]byte) {
	key := fmt.Sprintf(REDIS_KEY_STORE_ID, id)
	value := redis.GetString(key)
	logger.Info(nil,"get captcha to redis id:%s, value:%s", id, value)
	s := convert.StringToBytes(value)
	return s
}