package cache

import (
	"testing"
	"time"
)

func Test_Add(t *testing.T) {
	l := NewCache(3, LRU)
	l.Add("apple", 1)
	// time.Sleep(time.Second)
	l.Add("orange", 2)
	t.Log(l.Len())
	l.Add("apple", 3)
	l.Add("orange", 378)
	l.Add("orange", 313)
	l.Add("apple", 262)
	t.Log(l.Len())
	l.Remove("apple")
	t.Log(l.Len())
	l.(*Lru).PrintFunc = func(key, value interface{}, update time.Time) {
		t.Logf("key: %v, value: %v, update: %d\n", key, value, update.UnixNano())
	}
	l.(*Lru).OrderPrint(0)
	t.Log(l.Len())
	// x := l.Keys()
	// fmt.Println(x)
	// fmt.Println(l.Get("orange"))
}
