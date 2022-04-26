package cache

import "sync"

type Cacher[T1 comparable, T2 any] interface {
	Add(key T1, value T2)
	Remove(key T1)
	Len() int
	OrderPrint(int)
	Get(key T1) T2
	LastKey() T1
}

type Algorithm int

const (
	LRU Algorithm = iota
	LFU
	ALFU
)

func NewCache[T1 comparable, T2 any](n int, t Algorithm) Cacher[T1, T2] {
	// 内存足够的话, 可以设置很大, 所有计算都是O(1)
	if n <= 0 {
		n = 2 << 10
	}
	switch t {
	case LRU:
		return &Lru[T1, T2]{
			lru:  make(map[T1]*element[T1, T2]),
			size: n,
			lock: sync.RWMutex{},
			root: &element[T1, T2]{},
			last: &element[T1, T2]{},
		}
	case LFU:
		return &Lfu[T1, T2]{
			frequent: make(map[int]*Lru[T1, T2]),

			// 这里是根据key来查询在那一层
			cache: make(map[T1]int),
			mu:    sync.RWMutex{},
			size:  n,
		}
	case ALFU:
		alfu := &Alfu[T1, T2]{
			frequent: make(map[int]*Lru[T1, T2]),

			// 这里是根据key来查询在那一层
			cache: make(map[any]int),
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
