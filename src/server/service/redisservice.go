package service

import (
	"github.com/garyburd/redigo/redis"
	"server/utils/logger"
)

var REDIS_MAX_POOL_SIZE = 5
var REDIS_ADDRESS = ""
var REDIS_PASSWORD = "111111"
var redisPoll chan redis.Conn
var redisInit = false

func putRedis(conn redis.Conn) {
	if redisPoll == nil {
		redisPoll = make(chan redis.Conn, REDIS_MAX_POOL_SIZE)
	}
	if len(redisPoll) >= REDIS_MAX_POOL_SIZE {
		conn.Close()
		return
	}
	redisPoll <- conn
}

func initRedis(network, address string) redis.Conn {
	if len(redisPoll) == 0 && !redisInit {
		redisPoll = make(chan redis.Conn, REDIS_MAX_POOL_SIZE)
		go func() {
			for i := 0; i < REDIS_MAX_POOL_SIZE; i++ {
				initOneConn(network, address, i)
			}
			redisInit = true
		}()
	}
	return <-redisPoll
}

func initOneConn(network, address string, i int) {
	c, err := redis.Dial(network, address, redis.DialPassword(REDIS_PASSWORD))
	if err != nil {
		logger.Error(logger.ExtraFileds{"address": address, "index": i}, "init redis connection error, %s", err)
	} else {
		logger.Info(logger.ExtraFileds{"address": address, "index": i}, "init redis connection success")
		putRedis(c)
	}
}

func ExecuteSimple(cmd string, key string, value interface{}, resultChan chan []byte) (interface{}, error) {
	defer func() {
		if v := recover(); v != nil {
			logger.Info(nil, "ExecuteSimple error occur， %s", v)
		}
	}()
	c := initRedis("tcp", REDIS_ADDRESS)
	var result interface{}
	var err error
	switch cmd {
	case "LPUSH":
		result, err = redis.Bool(c.Do(cmd, key, value))
	case "BRPOP":
		var objs []string
		objs, err = redis.Strings(c.Do(cmd, key, value))
		if err == nil {
			result = objs[1]
		}
	case "SET":
		_, err = c.Do(cmd, key, value)
		if err == nil {
			result = 1
		}
	case "GET":
		//var objs interface{}
		objs, err := c.Do(cmd, key)
		if err == nil {
			result = objs
		}
	case "DECR":
		//var objs interface{}
		objs, err := c.Do(cmd, key)
		if err == nil {
			result = objs
		}
	case "INCR":
		//var objs interface{}
		objs, err := c.Do(cmd, key)
		if err == nil {
			result = objs
		}
	case "DEL":
		//var objs interface{}
		objs, err := c.Do(cmd, key)
		if err == nil {
			result = objs
		}
	case "EXPIRE":
		//var objs interface{}
		objs, err := c.Do(cmd, key, value)
		if err == nil {
			result = objs
		}
	case "KEYS":
		//var objs interface{}
		objs, err := c.Do(cmd, key)
		if err == nil {
			result = objs
		}
	case "PUBLISH":
		//var objs interface{}
		objs, err := c.Do(cmd, key, value)
		if err == nil {
			result = objs
		}
	case "SUBSCRIBE":
		psc := redis.PubSubConn{c}
		psc.Subscribe(key)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				logger.Debug(nil, "SUBSCRIBE %s: message: %s\n", v.Channel, v.Data)
				resultChan <- v.Data
			case redis.Subscription:
				logger.Debug(nil, "SUBSCRIBE %s: %s %d", v.Channel, v.Kind, v.Count)
			case error:
				logger.Error(logger.ExtraFileds{"key": key, "cmd": cmd}, "SUBSCRIBE error %s", v)
				err = v
			}
		}
	default:
		//var objs interface{}
		var objs interface{}
		if value == nil {
			objs, err = c.Do(cmd, key)
		} else {
			objs, err = c.Do(cmd, key, value)
		}
		if err == nil {
			result = objs
		}
	}
	checkConnection(c, err)
	return result, err
}

func Lpush(key string, value interface{}) (bool, error) {
	_, err := ExecuteSimple("LPUSH", key, value, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

func Brpop(key string, timeout int, valueChan chan interface{}) {
	value, err := ExecuteSimple("BRPOP", key, timeout, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "BRPOP", "key": key, "timeout": timeout}, "brpop command execute error %s", err)
	} else {
		valueChan <- value
	}
}

func Set(key string, value interface{}) {
	value, err := ExecuteSimple("SET", key, value, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "SET", "key": key, "value": value}, "set command execute error %s", err)
	}
}

func Get(key string) interface{} {
	value, err := ExecuteSimple("GET", key, nil, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "GET", "key": key, "value": value}, "get command execute error %s", err)
		return nil
	}
	return value
}

func GetString(key string) interface{} {
	r := Get(key)
	if r == nil {
		return nil
	}
	value, _ := redis.String(r, nil)
	return value
}

func GetInt(key string) int {
	r := Get(key)
	if r == nil {
		return 0
	}
	value, _ := redis.Int(r, nil)
	return value
}

func DECR(key string) {
	value, err := ExecuteSimple("DECR", key, nil, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "DECR", "key": key, "value": value}, "DECR command execute error %s", err)
	}
}

func INCR(key string) {
	value, err := ExecuteSimple("INCR", key, nil, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "INCR", "key": key, "value": value}, "INCR command execute error %s", err)
	}
}

func DEL(key string) {
	value, err := ExecuteSimple("DEL", key, nil, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "DEL", "key": key, "value": value}, "DEL command execute error %s", err)
	}
}

func EXPIRE(key string, seconds int) {
	value, err := ExecuteSimple("EXPIRE", key, seconds, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "EXPIRE", "key": key, "value": value}, "EXPIRE command execute error %s", err)
	}
}

func KEYS(key string) []string {
	value, err := ExecuteSimple("KEYS", key, nil, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "KEYS", "key": key, "value": value}, "KEYS command execute error %s", err)
	}
	if value == nil {
		return make([]string, 0)
	}
	result, e := redis.Strings(value, err)
	if e == nil {
		return result
	} else {
		logger.Error(logger.ExtraFileds{"cmd": "KEYS", "key": key, "value": value}, "convert result to []string error [KEYS] %s", e)
		return make([]string, 0)
	}
}

func PUBLISH(key string, value string) {
	_, err := ExecuteSimple("PUBLISH", key, value, nil)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "PUBLISH", "key": key, "value": value}, "PUBLISH command execute error %s", err)
	}
}

func SUBSCRIBE(key string, subchan chan []byte) {
	value, err := ExecuteSimple("SUBSCRIBE", key, nil, subchan)
	if err != nil {
		logger.Error(logger.ExtraFileds{"cmd": "SUBSCRIBE", "key": key, "value": value}, "SUBSCRIBE command execute error %s", err)
	}
}

func ZADD(key string, score int32, value string) (interface{}, error) {
	defer func() {
		if v := recover(); v != nil {
			logger.Info(nil, "ZADD execute error occur， %s", v)
		}
	}()
	c := initRedis("tcp", REDIS_ADDRESS)
	var result interface{}
	var err error
	objs, err := c.Do("ZADD", key, score, value)
	if err == nil {
		result = objs
	}
	checkConnection(c, err)
	return result, err
}

func ZRANGE(key string, start int32, end int32, withscores bool) ([]interface{}, error) {
	defer func() {
		if v := recover(); v != nil {
			logger.Info(nil, "ZRANGE execute error occur， %s", v)
		}
	}()
	c := initRedis("tcp", REDIS_ADDRESS)
	var result []interface{}
	var err error
	var objs interface{}
	if withscores {
		objs, err = c.Do("ZRANGE", key, start, end, "WITHSCORES")
	} else {
		objs, err = c.Do("ZRANGE", key, start, end)

	}
	if err == nil {
		result = objs.([]interface{})
	}
	checkConnection(c, err)
	return result, err
}

func ZCARD(key string) (int32, error) {
	defer func() {
		if v := recover(); v != nil {
			logger.Info(nil, "ZCARD execute error occur， %s", v)
		}
	}()
	c := initRedis("tcp", REDIS_ADDRESS)
	var result int32
	var err error
	objs, err := c.Do("ZCARD", key)
	if err == nil {
		result = int32(objs.(int64))
	}
	checkConnection(c, err)
	return result, err
}

func ZREMRANGEBYRANK(key string, start int32, end int32) {
	defer func() {
		if v := recover(); v != nil {
			logger.Info(nil, "ZREMRANGEBYRANK execute error occur， %s", v)
		}
	}()
	c := initRedis("tcp", REDIS_ADDRESS)
	var err error
	_, err = c.Do("ZREMRANGEBYRANK", key, start, end)
	checkConnection(c, err)
}

/**
  验证链接
*/
func checkConnection(c redis.Conn, err error) {
	if err != nil {
		obj, e := redis.String(c.Do("PING"))
		if err == nil && obj == "PONG" {
			redisPoll <- c
		} else {
			//出现异常丢弃链接
			logger.Error(nil, "execute redis cmd error, redis connection error %s", e)
			c.Close()
			initOneConn("tcp", REDIS_ADDRESS, 0) //补充一个
		}
	} else {
		redisPoll <- c
	}
}

//
//func main()  {
//	ok,err := Lpush("test:key2","hi")
//	if ok{
//		fmt.Println("success")
//	}else{
//		fmt.Println(err)
//	}
//
//	valueChan := make(chan interface{},10)
//	go func(){
//		for{
//			select {
//			case s := <-valueChan:
//				fmt.Println(s)
//			}
//		}
//	}()
//	for {
//		Brpop("test:key2", 10, valueChan)
//	}
//}

func main() {
	REDIS_ADDRESS = "192.168.19.141:6379"
	ZADD("test", 1, "just a play")
}
