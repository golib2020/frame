package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golib2020/frame/redis"
)

type redisCache struct {
	prefix string
	pool   redis.Redis
}

func NewRedis(prefix string, pool redis.Redis) Cache {
	return &redisCache{
		prefix: prefix,
		pool:   pool,
	}
}

func (d *redisCache) Has(key string) bool {
	var b bool
	if err := d.pool.Do(&b, "EXISTS", d.getKey(key)); err != nil {
		return false
	}
	return b
}

func (d *redisCache) Get(key string, res interface{}) error {
	var bts []byte
	if err := d.pool.Do(&bts, "GET", d.getKey(key)); err != nil {
		return err
	}
	return json.Unmarshal(bts, res)
}

func (d *redisCache) Set(key string, data interface{}, ex ...time.Duration) error {
	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(ex) > 0 {
		exs := fmt.Sprintf("%f", ex[0].Seconds())
		return d.pool.Do(nil, "SET", d.getKey(key), string(bts), "EX", exs)
	}
	return d.pool.Do(nil, "SET", d.getKey(key), string(bts))
}

func (d *redisCache) Del(key string) error {
	return d.pool.Do(nil, "DEL", d.getKey(key))
}

func (d *redisCache) getKey(key string) string {
	return fmt.Sprintf("%s.%s", d.prefix, key)
}
