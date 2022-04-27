package cache

import "sync"

type Cacher[K comparable, V any] interface {
	Add(key K, value V)
	Remove(key K)
	Len() int
	OrderPrint(int)
	Get(key K) (V, bool)
	LastKey() K
}

type Algorithm int

const (
	LRU Algorithm = iota
	LFU
	ALFU
)

func NewCache[K comparable, V any](n int, t Algorithm) Cacher[K, V] {
	// 内存足够的话, 可以设置很大, 所有计算都是O(1)
	if n <= 0 {
		n = 2 << 10
	}
	switch t {
	case LRU:
		return &Lru[K, V]{
			lru:  make(map[K]*element[K, V]),
			size: n,
			lock: sync.RWMutex{},
			root: &element[K, V]{},
			last: &element[K, V]{},
		}
	case LFU:
		return &Lfu[K, V]{
			frequent: make(map[int]*Lru[K, V]),

			// 这里是根据key来查询在那一层
			cache: make(map[K]int),
			mu:    sync.RWMutex{},
			size:  n,
		}
	case ALFU:
		alfu := &Alfu[K, V]{
			frequent: make(map[int]*Lru[K, V]),

			// 这里是根据key来查询在那一层
			cache: make(map[K]int),
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
