package cache

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
)
type htRedis struct {
	pool *redis.Pool
}
var redisInstanceTool *htRedis
var isOpen string

func TransmitRedisPool(redisPool *htRedis,cacheIsOpen string)  {
	redisInstanceTool = redisPool
	isOpen = cacheIsOpen
}

func Redis() *htRedis {
	return redisInstanceTool
}

func (htr *htRedis) Do(command string, args ...interface{}) (interface{}, error) {
	if isOpen == "false" {
		return nil, errors.New("不走缓存")
	}
	var conn = htr.pool.Get()
	defer conn.Close()
	return conn.Do(command, args...)
}

func (htr *htRedis) TryLockWithTimeout(key string, value string, timeout int) (ok bool, err error) {
	_, err = redis.String(Redis().Do("SET", key, value, "EX", timeout, "NX"))
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (htr *htRedis) GetString(key string) (ok bool, value string) {
	value, err := redis.String(htr.Do("GET", key))
	if err == redis.ErrNil {
		return false, ""
	}
	if err != nil {
		return false, err.Error()
	}
	return true, value
}

func (htr *htRedis) SetString(key string, value string) (ok bool, err error) {
	if _, err := redis.String(htr.Do("SET", key, value)); err != nil {
		return false, err
	}
	return true, nil
}

func (htr *htRedis) SetStringEX(key, value string, ex int) (ok bool, err error) {
	if _, err := redis.String(htr.Do("SET", key, value, "EX", ex)); err != nil {
		return false, err
	}
	return true, nil
}

func (htr *htRedis) SetEX(key string, value interface{}, ex int) (ok bool, err error) {
	var valueStr []byte
	if valueStr, err = json.Marshal(value); err != nil {
		return false, err
	}
	if _, err := redis.String(htr.Do("SET", key, string(valueStr), "EX", ex)); err != nil {
		return false, err
	}
	return true, nil
}

func (htr *htRedis) Get(key string) (ok bool, value string) {
	value, err := redis.String(htr.Do("GET", key))
	if err == redis.ErrNil {
		return false, ""
	}
	if err != nil {
		return false, err.Error()
	}
	return true, value
}

func (htr *htRedis) Set(key string, value interface{}) (ok bool, err error) {
	var valueStr []byte
	if valueStr, err = json.Marshal(value); err != nil {
		return false, err
	}
	if _, err := redis.String(htr.Do("SET", key, string(valueStr))); err != nil {
		return false, err
	}
	return true, nil
}
