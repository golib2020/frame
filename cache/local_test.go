package cache

import (
	"os"
	"testing"
	"time"
)

func TestLocal(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{
			"key1",
			"value1",
		},
		{
			"key2",
			"value2",
		},
	}

	c := NewLocal("", "")

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if err := c.Set(tt.key, tt.value, time.Minute); err != nil {
				t.Errorf("Set() %v", err)
			}
			if !c.Has(tt.key) {
				t.Errorf("Has() error")
			}
			var s string
			if err := c.Get(tt.key, &s); err != nil {
				t.Errorf("Get() %v", err)
			}
			if s != tt.value {
				t.Errorf("Get() value %s", s)
			}
			if err := c.Del(tt.key); err != nil {
				t.Errorf("Del() %v", err)
			}
		})
	}

	os.RemoveAll("./runtime")

}
