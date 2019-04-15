package framework

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func InitRedisPool(connString string) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxActive: 500,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", connString,
				redis.DialConnectTimeout(60*time.Second),
				redis.DialReadTimeout(60*time.Second),
				redis.DialWriteTimeout(60*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
	return pool, nil
}

func NewRedisConn(p *redis.Pool) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		return p.Get(), nil
	}
}
