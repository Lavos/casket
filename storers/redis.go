package storers

import (
	"time"
	"log"
	"fmt"
	"strings"
	"github.com/garyburd/redigo/redis"
)

type RedisStorer struct {
	Pool      *redis.Pool
	Namespace string
}

func (r *RedisStorer) Init(host, namespace string, port int64) {
	pool := &redis.Pool{
		MaxIdle:     5,
		MaxActive:   1000,
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

			if err != nil {
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	r.Pool = pool
	r.Namespace = namespace
}

func (r *RedisStorer) GetFullKey(key string) string {
	return fmt.Sprintf("%s:%s", r.Namespace, key)
}

func (r *RedisStorer) GetSubKey(key string) string {
	return strings.TrimPrefix(key, fmt.Sprintf("%s:", r.Namespace))
}

