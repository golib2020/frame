package f

import (
	"fmt"

	"github.com/golib2020/frame/cache"
)

const (
	cacheInstancesPrefix = `cache`
	cacheDefaultGroup    = `default`
)

//Cache 缓存
func Cache(name ...string) cache.Cache {
	group := cacheDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", cacheInstancesPrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return cacheInit(key)
	}).(cache.Cache)
}

func cacheInit(key string) cache.Cache {
	conf := Config().Sub(key)
	var s cache.Cache
	switch conf.GetString("driver") {
	case "local":
		conf.SetDefault("root", "./")
		s = cache.NewLocal(conf.GetString("prefix"), conf.GetString("root"))
	case "redis":
		s = cache.NewRedis(conf.GetString("prefix"), Redis())
	}
	return s
}
