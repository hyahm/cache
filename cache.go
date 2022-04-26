package cache

import "sync"

type Cacher interface {
	Add(key, value interface{})
	Remove(key interface{})
	Len() int
	OrderPrint(int)
	Get(key interface{}) interface{}
	LastKey() interface{}
}

type Algorithm int

const (
	LRU Algorithm = iota
	LFU
	ALFU
)

func NewCache(n int, t Algorithm) Cacher {
	// 内存足够的话, 可以设置很大, 所有计算都是O(1)
	if n <= 0 {
		n = 2 << 10
	}
	switch t {
	case LRU:
		return &Lru{
			lru:  make(map[interface{}]*element, 0),
			size: n,
			lock: sync.RWMutex{},
			root: &element{},
			last: &element{},
		}
	case LFU:
		return &Lfu{
			frequent: make(map[int]*Lru),

			// 这里是根据key来查询在那一层
			cache: make(map[interface{}]int),
			mu:    sync.RWMutex{},
			size:  n,
		}
	case ALFU:
		alfu := &Alfu{
			frequent: make(map[int]*Lru),

			// 这里是根据key来查询在那一层
			cache: make(map[interface{}]int),
			mu:    sync.RWMutex{},
			size:  DEFAULTCOUNT,
		}
		go alfu.auto()
		return alfu
	default:
		return nil
	}

}

const DEFAULTCOUNT = 2 << 10
