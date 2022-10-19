package cache

import (
	"fmt"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func Test_Get(t *testing.T) {
	cache := New(int64(0), nil)
	cache.Add("wu", String("ff"))
	if v, ok := cache.Get("wu"); !ok || string(v.(String)) != "ff" {
		t.Fatalf("cache Get failed")
	}
}

func Test_Missing(t *testing.T) {
	cache := New(int64(0), nil)
	if _, ok := cache.Get(("he")); ok {
		t.Fatalf("cache missing failed")
	}
}

func Test_TriggerRemove(t *testing.T) {
	callback := func(key string, value Value) {
		t.Log(key, value)
	}
	cache := New(int64(40), callback)

	for i := 1; i <= 10; i++ {
		cache.Add(fmt.Sprintf("key%d", i), String("value"))
	}
}
