package database

import (
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"ssh_manage/config"
	"time"
)

var redis_conf = config.Config.Redis

var Cache *redigo.Pool

func init() {
	var addr = fmt.Sprintf("%s:%d",redis_conf.Host,redis_conf.Port)
	var password = redis_conf.Password
	Cache = poolInitRedis(addr, password)
}

func poolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   redis_conf.Poolsize, //最大数
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
