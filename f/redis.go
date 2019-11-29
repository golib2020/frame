package f

import (
	"fmt"
	//"github.com/gomodule/redigo/redis"

	"github.com/golib2020/frame/redis"
)

const (
	redisInstancesPrefix = `redis`
	redisDefaultGroup    = `default`
)

//Redis redis实例
func Redis(name ...string) redis.Redis {
	group := redisDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", redisInstancesPrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return redisInit(key)
	}).(redis.Redis)
}

func redisInit(key string) redis.Redis {
	conf := Config().Sub(key)
	var c redis.Redis
	switch conf.GetString("driver") {
	case "radix":
		c = redis.NewRadixRedis(
			redis.WithAddr(conf.GetString("addr")),
			redis.WithPass(conf.GetString("pass")),
			redis.WithSize(conf.GetInt("size")),
			redis.WithSelectDB(conf.GetInt("db")),
		)
	}
	return c
}
