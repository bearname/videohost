package caching

import (
	"github.com/bearname/videohost/pkg/common/model"
	"github.com/go-redis/redis"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(dsn model.DSN) *RedisCache {
	r := new(RedisCache)

	r.client = redis.NewClient(&redis.Options{
		Addr:     dsn.Address,
		Password: dsn.Password, // no password set
		DB:       0,            // use default DB
	})

	return r
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.client.Get(key).Result()
}

func (c *RedisCache) Set(key string, value string) error {
	return c.client.Set(key, value, time.Hour*24*365).Err()
}

func (c *RedisCache) Del(key string) error {
	return c.client.Del(key).Err()
}
