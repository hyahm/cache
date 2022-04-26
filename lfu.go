package cache

import (
	"fmt"
	"sync"
)

// 以lru为基础

type Lfu[T1 comparable, T2 any] struct {
	frequent map[int]*Lru[T1, T2]

	// 这里是根据key来查询在那一层
	cache map[T1]int
	min   int // 记录当前最小层的值
	mu    sync.RWMutex
	size  int // 大小
}

func (lfu *Lfu[T1, T2]) OrderPrint(frequent int) {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	for frequent, lru := range lfu.frequent {
		fmt.Printf("%#v\n", lru)
		lru.OrderPrint(frequent)
	}

}

// 为了方便修改， 一样也需要一个双向链表

func (lfu *Lfu[T1, T2]) add(index int, key T1, value T2) {
	if _, ok := lfu.frequent[index]; !ok {
		lfu.frequent[index] = &Lru[T1, T2]{
			lru:  make(map[T1]*element[T1, T2], 0),
			size: lfu.size,
			lock: sync.RWMutex{},
			root: &element[T1, T2]{},
			last: &element[T1, T2]{},
		}
	}

	lfu.frequent[index].Add(key, value)
}

func (lfu *Lfu[T1, T2]) getMin(start int) int {
	if lfu.frequent[start].Len() > 0 {
		return start
	} else {
		return lfu.getMin(start + 1)
	}
}

func (lfu *Lfu[T1, T2]) Len() int {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	return len(lfu.cache)
}

// get lastKey
func (lfu *Lfu[T1, T2]) LastKey() T1 {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	return lfu.frequent[lfu.min].LastKey()
}

func (lfu *Lfu[T1, T2]) Remove(key T1) {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()
	// 先找到这个key
	if index, ok := lfu.cache[key]; ok {
		if _, ok := lfu.frequent[index]; ok {
			lfu.frequent[index].Remove(key)
		}
		delete(lfu.cache, key)
	}
}

func (lfu *Lfu[T1, T2]) Add(key T1, value T2) {
	// 添加一个key
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	if li, ok := lfu.cache[key]; ok {
		// 如果存在的话，删除此层的值，
		lfu.frequent[li].Remove(key)
		//添加到新层中
		// 判断是否存在新层， 不存在就新建
		lfu.cache[key] = li + 1
		lfu.add(li+1, key, value)
	} else {
		lfu.cache[key] = 1
		lfu.min = 1
		lfu.add(1, key, value)
		// 判断是否超过了缓存值
		if len(lfu.cache) >= int(lfu.size) {
			// 删除最后一个
			removeKey := lfu.frequent[lfu.min].RemoveLast()
			// 删除总缓存
			delete(lfu.cache, removeKey)
			if lfu.frequent[lfu.min].Len() == 0 {
				// 如果长度为空， 我们就要重新获取最小层
				// delete(frequent, min)
				// 继续取最小层数
				lfu.min = lfu.getMin(lfu.min + 1)
			}
		}
	}

}

//
func (lfu *Lfu[T1, T2]) Get(key T1) any {
	lfu.mu.RLock()
	defer lfu.mu.RUnlock()
	if index, ok := lfu.cache[key]; ok {
		if v, ok := lfu.frequent[index]; ok {
			return v.Get(key)
		}
	}
	return nil

}
