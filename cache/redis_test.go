package cache

import (
	"testing"
)

func TestRedis(t *testing.T) {
	/*
		radixRedis := redis.NewRadixRedis("127.0.0.1:6379", "", 0, 10)
		cache := NewRedis("test.cache", radixRedis)
		//set
		if err := cache.Set("key", 1); err != nil {
			t.Errorf("set %s", err)
		}
		//has
		if !cache.Has("key") {
			t.Errorf("has err")
		}
		//get
		var res int
		if err := cache.Get("key", &res); err != nil {
			t.Errorf("get %s", err)
		}
		if res != 1 {
			t.Errorf("get %d", res)
		}
		//del
		if err := cache.Del("key"); err != nil {
			t.Errorf("del %s", err)
		}
	*/
}
