package composites

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisComposite struct {
	redis *redis.Pool
}

func NewRedisComposite() (*RedisComposite, error) {
	redisPool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(os.Getenv("REDISPROTOCOL"), os.Getenv("REDISHOST")+":"+os.Getenv("REDISPORT"))
		},
		DialContext:     nil,
		TestOnBorrow:    nil,
		MaxIdle:         10,
		MaxActive:       0,
		IdleTimeout:     240 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	return &RedisComposite{redis: &redisPool}, nil
}
