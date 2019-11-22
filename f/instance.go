package f

import (
	"fmt"
	"sync"
)

type Store interface {
	GetOrSet(key string, data interface{}) interface{}
	GetOrSetFunc(key string, data func() interface{}) interface{}
}

type store struct {
	data sync.Map
}

func NewStore() Store {
	return &store{
		data: sync.Map{},
	}
}

func (s *store) GetOrSet(key string, data interface{}) interface{} {
	value, ok := s.data.Load(key)
	if ok {
		return value
	}
	s.data.Store(key, data)
	return data
}

func (s *store) GetOrSetFunc(key string, data func() interface{}) interface{} {
	value, ok := s.data.Load(key)
	if ok {
		return value
	}
	value = data()
	s.data.Store(key, value)
	fmt.Printf("%-20s init success. \n", key)
	return value
}


var instances = NewStore()

func Instance() Store {
	return instances
}