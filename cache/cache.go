package cache

import (
	"time"
)

type Cache interface {
	Has(key string) bool
	Get(key string) (data string, err error)
	Set(key string, data string, ex ...time.Duration) error
	Del(key string) error
}
