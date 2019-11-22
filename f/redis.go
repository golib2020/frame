package f

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	redisInstancesPrefix = `redis`
	redisDefaultGroup    = `default`
)

//Redis redis实例
func Redis(name ...string) *redis.Pool {
	group := redisDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", redisInstancesPrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return redisInit(key)
	}).(*redis.Pool)
}

func redisInit(key string) *redis.Pool {
	conf := Config().Sub(key)
	db := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.GetString("addr"), redis.DialPassword(conf.GetString("pass")))
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", 0)
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     conf.GetInt("max.idle"),
		MaxActive:   conf.GetInt("max.active"),
		IdleTimeout: time.Minute,
		Wait:        true,
	}
	return db
}
