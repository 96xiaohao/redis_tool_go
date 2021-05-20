package cache

import (
	"github.com/sirupsen/logrus"
	"gitlab.yc345.tv/backend/onion_util"
)

type cacheConfig struct {
	Server      string `json:"server"`
	Password    string `json:"password"`
	Database    int    `json:"database"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	IdleTimeout int    `json:"idleTimeout"`
	*cacheLife
}

type cacheLife struct {
	RedisKeyLifespan   int `json:"redisKeyLifespan"`
	Lifespan           int `json:"lifespan"`
	CachePurgeInterval int `json:"cachePurgeInterval"`
}

var cacheInitFuncs = map[string]func(c *cacheConfig){}

var cacheLifeDetails = map[string]*cacheLife{}

var log *logrus.Entry

var cacheIsOpen  string

func InitCache(server, password, isOpenCache string, dataBase, maxIdle, maxActive, idleTimeout, redisKeyLifespan, lifespan, cachePurgeInterval int) {
	cacheIsOpen = isOpenCache
	log = onion_util.AppLog(4, "cache")
	var cc = &cacheConfig{
		Server:      server,
		Password:    password,
		Database:    dataBase,
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		cacheLife: &cacheLife{
			RedisKeyLifespan:   redisKeyLifespan,
			CachePurgeInterval: cachePurgeInterval,
			Lifespan:           lifespan,
		},
	}
	RedisInit(cc)
	log.Info("Redis connect: ", *cc)
	for k, f := range cacheInitFuncs {
		ccCopy := *cc
		if v, exists := cacheLifeDetails[k]; exists && v != nil {
			cacheLifeInfo := *ccCopy.cacheLife
			if v.RedisKeyLifespan != 0 {
				cacheLifeInfo.RedisKeyLifespan = v.RedisKeyLifespan
			}
			if v.Lifespan != 0 {
				cacheLifeInfo.Lifespan = v.Lifespan
			}
			if v.CachePurgeInterval != 0 {
				cacheLifeInfo.CachePurgeInterval = v.CachePurgeInterval
			}
			ccCopy.cacheLife = &cacheLifeInfo
		}
		f(&ccCopy)
	}
}
