package cache

import (
	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	pool *redis.Pool
}

func NewCache() *Cache {
	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}

	return &Cache{pool: redisPool}
}
