package cache

import (
	"time"
)

type Cache interface {
	Has(key string) bool
	Get(key string, data interface{}) error
	Set(key string, data interface{}, ex ...time.Duration) error
	Del(key string) error
}
