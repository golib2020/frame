package cache

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisCache struct {
	prefix string
	pool   *redis.Pool
}

func NewRedis(prefix string, pool *redis.Pool) Cache {
	return &redisCache{
		prefix: prefix,
		pool:   pool,
	}
}

func (d *redisCache) Has(key string) bool {
	conn := d.pool.Get()
	defer conn.Close()
	b, err := redis.Bool(conn.Do("EXISTS", d.getKey(key)))
	if err != nil {
		return false
	}
	return b
}

func (d *redisCache) Get(key string) (string, error) {
	conn := d.pool.Get()
	defer conn.Close()
	content, err := redis.String(conn.Do("GET", d.getKey(key)))
	if err != nil {
		return "", err
	}
	return content, nil
}

func (d *redisCache) Set(key string, data string, ex ...time.Duration) error {
	conn := d.pool.Get()
	defer conn.Close()

	if len(ex) > 0 {
		_, err := conn.Do("SET", d.getKey(key), data, "EX", int(ex[0].Seconds()))
		return err
	}
	_, err := conn.Do("SET", d.getKey(key), data)
	return err
}

func (d *redisCache) Del(key string) error {
	conn := d.pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", d.getKey(key))
	return err
}

func (d *redisCache) getKey(key string) string {
	return fmt.Sprintf("%s.%s", d.prefix, key)
}
