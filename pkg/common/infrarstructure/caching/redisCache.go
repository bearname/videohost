package caching

import (
	"github.com/bearname/videohost/pkg/common/caching"
	"github.com/bearname/videohost/pkg/common/model"
	"github.com/go-redis/redis"
	"time"
)
const defaultDB = 0

type RedisCache struct {
	client *redis.Client
	isOk   bool
}

func NewRedisCache(dsn model.DSN) *RedisCache {
	r := new(RedisCache)

	r.client = redis.NewClient(&redis.Options{
		Addr:     dsn.Address,
		Password: dsn.Password,
		DB:       defaultDB,
	})

	ping := r.client.Ping()
	if ping.Err() != nil {
		r.isOk = false
	} else {
		r.isOk = true
	}

	return r
}

func (c *RedisCache) IsOk() bool {
	return c.isOk
}

func (c *RedisCache) Get(key string) (string, error) {
	if !c.isOk {
		return "", caching.ErrCacheUnavailable
	}
	return c.client.Get(key).Result()
}

func (c *RedisCache) Set(key string, value string) error {
	if !c.isOk {
		return caching.ErrCacheUnavailable
	}
	return c.client.Set(key, value, time.Hour*24*365).Err()
}

func (c *RedisCache) Del(key string) error {
	if !c.isOk {
		return caching.ErrCacheUnavailable
	}
	return c.client.Del(key).Err()
}
