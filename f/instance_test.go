package f

import (
	"testing"
)

func TestStore_GetOrSet(t *testing.T) {
	es := []struct {
		key    string
		data   interface{}
		answer interface{}
	}{
		{
			key:    "key",
			data:   "data",
			answer: "data",
		},
		{
			key:    "key",
			data:   "data2",
			answer: "data",
		},
		{
			key:    "key",
			data:   "data3",
			answer: "data",
		},
	}
	// default
	m := NewStore()
	for _, e := range es {
		res := m.GetOrSet(e.key, e.data)
		if res != e.answer {
			t.Errorf("key: %s, data: %v, answer: %v, put: %v", e.key, e.data, e.answer, res)
		}
	}
}

func TestStore_GetOrSetFunc(t *testing.T) {
	es := []struct {
		key    string
		data   func() interface{}
		answer interface{}
	}{
		{
			key: "key",
			data: func() interface{} {
				return "data"
			},
			answer: "data",
		},
		{
			key: "key",
			data: func() interface{} {
				return "data2"
			},
			answer: "data",
		},
		{
			key: "key",
			data: func() interface{} {
				return "data3"
			},
			answer: "data",
		},
	}
	//with func
	m := NewStore()
	for _, e := range es {
		res := m.GetOrSetFunc(e.key, e.data)
		if res != e.answer {
			t.Errorf("key: %s, data: %T, answer: %v, put: %v", e.key, e.data, e.answer, res)
		}
	}
}
