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

func TransmitRedisPool(redisPool *redis.Pool,cacheIsOpen string)  {
	redisInstanceTool = &htRedis{}
	redisInstanceTool.pool = redisPool
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
func (htr *htRedis) ZrangeByScore(k string, n, m, point interface{}) error {
	values, err := redis.Values(htr.Do("ZRANGEBYSCORE", k, n, m, "WITHSCORES"))
	if err != nil {
		return err
	}
	return redis.ScanSlice(values, point)
}

func (htr *htRedis) Zincrby(k, mem string, socre int) error {
	_, err := redis.Int(htr.Do("ZINCRBY", k, socre, mem))
	if err != nil {
		return err
	}
	return nil
}

func (htr *htRedis) Zadd(k, mem string, score int) error {
	_, err := redis.Int(htr.Do("ZADD", k, score, mem))
	if err != nil {
		return err
	}
	return nil
}

func (htr *htRedis) Zrem(k, mem string) error {
	_, err := redis.Int(htr.Do("ZREM", k, mem))
	if err != nil {
		return err
	}
	return nil
}

func (htr *htRedis) Del(k string) error {
	_, err := redis.Int(htr.Do("DEL", k))
	if err != nil {
		return err
	}
	return nil
}

// func (htr *htRedis) ZAddArr(k string, mems []string) error {
// 	if len(mems)%2 != 0 {
// 		return errors.New("mems is fail")
// 	}
// 	setMems := []string{}
// 	setMems = append(setMems, "ZADD")
// 	setMems = append(setMems, "k")
// 	setMems = append(setMems, mems...)
// 	var conn = htr.pool.Get()
// 	defer conn.Close()
// 	_, err := redis.Int(conn.Do(setMems...))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (htr *htRedis) Sadd(k, mem string) error {
	_, err := redis.Int(htr.Do("SADD", k, mem))
	if err != nil {
		return err
	}
	return nil
}

func (htr *htRedis) Smembers(k string) *[]string {
	if values, err := redis.Strings(htr.Do("SMEMBERS", k)); err != nil {
		return &[]string{}
	} else {
		return &values
	}
}

func (htr *htRedis) Srem(key, member string) error {
	if _, err := redis.Int(htr.Do("SREM", key, member)); err != nil {
		return err
	}
	return nil
}

func (htr *htRedis) Rpush(k, mem string) error {
	_, err := redis.Int(htr.Do("RPUSH", k, mem))
	if err != nil {
		return err
	}
	return nil
}

func (htr *htRedis) Lrange(k string, start, end int) []string {
	strs, err := redis.Strings(htr.Do("LRANGE", k, start, end))
	if err != nil {
		return []string{}
	}
	return strs
}

func (htr *htRedis) Hset(keyName, fieldName string, value interface{}) (ok bool, err error) {
	if _, err := htr.Do("HSET", keyName, fieldName, value); err != nil {
		return false, err
	}

	return true, nil
}

func (htr *htRedis) Hmset(keyName string, args ...interface{}) (ok bool, err error) {
	argumnets := []interface{}{keyName}
	argumnets = append(argumnets, args...)
	if _, err := htr.Do("HSET", argumnets...); err != nil {
		return false, err
	}
	return true, nil
}

func (htr *htRedis) HGETALL(keyName string) (value map[string]string, err error) {

	result, err := redis.StringMap(htr.Do("HGETALL", keyName))
	if err != nil {
		return nil, err
	}
	return result, nil
}
