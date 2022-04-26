package cache

import (
	"fmt"
	"sync"
	"time"
)

type element[T1 comparable, T2 any] struct {
	// 上一个元素和下一个元素
	next, prev *element[T1, T2]
	// The list to which this element belongs.
	//元素的key
	key T1
	// 这个元素的值
	value  T2
	update time.Time
}

type Lru[T1 comparable, T2 any] struct {
	lru map[T1]*element[T1, T2] //  这里存key 和 元素
	//保存第一个元素
	lock      sync.RWMutex
	root      *element[T1, T2] // sentinel list element, only &root, root.prev, and root.next are used
	last      *element[T1, T2] // 最后一个元素
	len       int              // 元素长度
	size      int              // 缓存多少元素
	PrintFunc func(key, value T1, update time.Time)
}

func defaultLru[T1 comparable, T2 any]() *Lru[T1, T2] {
	return &Lru[T1, T2]{
		lru:  make(map[T1]*element[T1, T2], 0),
		size: DEFAULTCOUNT,
		lock: sync.RWMutex{},
		root: &element[T1, T2]{},
		last: &element[T1, T2]{},
	}
}

//开始存一个值
func (l *Lru[T1, T2]) Add(key T1, value T2) {
	if l.lru == nil {
		l = defaultLru[T1, T2]()
	}

	l.add(key, value)
}

// 获取值
func (l *Lru[T1, T2]) Get(key T1) T2 {
	if l.lru == nil {
		l = defaultLru[T1, T2]()
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.lru[key].value
}

func (l *Lru[T1, T2]) GetLastKeyUpdateTime() (T1, T2, time.Time) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.last.key, l.last.value, l.last.update
}

func (l *Lru[T1, T2]) NextKey(key T1) any {
	if l.lru == nil {
		l = defaultLru[T1, T2]()
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	if value, ok := l.lru[key]; ok {
		if value.next == nil {
			return nil
		}
		return value.next.key
	}
	return nil
}

func (l *Lru[T1, T2]) PrevKey(key T1) any {
	if l.lru == nil {
		l = defaultLru[T1, T2]()
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	if value, ok := l.lru[key]; ok {
		if value.prev == nil {
			return nil
		}
		return value.prev.key
	}
	return nil
}

func (l *Lru[T1, T2]) Remove(key T1) {
	if l.lru == nil {
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	// 不存在就直接返回
	if _, ok := l.lru[key]; !ok {
		return
	}
	this := l.lru[key]
	//如果是第一个元素
	if this == l.root {
		l.root = l.root.next
		l.root.prev = nil
		// 更新第二个元素的值
		l.lru[l.root.key] = l.root
		delete(l.lru, key)
		l.len--
		return
	}
	//如果是最后一个
	if this == l.last {
		l.last = l.last.prev
		l.last.next = nil
		l.lru[l.last.key] = l.last
		delete(l.lru, key)
		l.len--
		return
	}

	// 更改上一个元素的下一个值

	l.lru[key].prev.next = l.lru[key].next
	//更新下一个元素的上一个值
	l.lru[key].next.prev = l.lru[key].prev
	l.lru[l.lru[key].prev.key] = l.lru[key].prev
	l.lru[l.lru[key].next.key] = l.lru[key].next
	//删除
	delete(l.lru, key)
	l.len--
}

func (l *Lru[T1, T2]) OrderPrint(frequent int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	li := l.root
	for li != nil {
		fmt.Printf("key: %v, value: %v, frequent: %d, update_time: %s\n", li.key, li.value, frequent, li.update.String())
		li = li.next
	}

}

func (l *Lru[T1, T2]) Len() int {
	return l.len
}

func (l *Lru[T1, T2]) Resize(n int) {
	//如果缩小了缓存, 那么可能需要删除后面多余的索引
	l.size = n
	l.lock.Lock()
	defer l.lock.Unlock()
	if n < l.size {
		for l.len > n {
			l.removeLast()
		}
	}
}

// 返回被删除的key, 如果没删除返回nil
func (l *Lru[T1, T2]) add(key T1, value T2) any {
	//先要判断是否存在这个key, 存在的话，就将元素移动最开始的位置,
	l.lock.Lock()
	defer l.lock.Unlock()

	if _, ok := l.lru[key]; ok {
		//如果是第一个元素的话, 什么也不用操作
		if l.root.key == key {
			l.root.value = value
			l.root.update = time.Now()
			return nil
		} else {

			// 否则就插入到开头, 开头的元素后移
			l.moveToPrev(key, value)
		}

	} else {

		el := &element[T1, T2]{
			prev:   nil,
			next:   nil,
			update: time.Now(),
			key:    key,
			value:  value,
		}

		if l.len == 0 {
			// 更新第一个元素
			l.root = el
			// 更新最后一个元素
			l.last = el
			// 更新长度
			l.len = 1
			// 更新lru
			l.lru[key] = el
			return nil
		}
		//如果不存在的话, 直接添加到开头
		// 第二个元素抽取出来， 也就是当前的root
		second := l.root
		second.prev = el

		el.next = l.root
		//将开头的元素修改成新的元素
		l.root = el
		l.root.next = second
		l.lru[key] = l.root
		l.lru[l.root.next.key] = l.root.next
		//判断长度是否超过了缓存
		if l.Len() >= l.size {
			key := l.removeLast()

			return key
		}
		l.len++

	}
	return nil
}

// 移除最后一个, 返回key
func (l *Lru[T1, T2]) RemoveLast() T1 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.removeLast()
}

func (l *Lru[T1, T2]) removeLast() T1 {
	lastKey := l.last.key

	if l.last.prev != nil {
		l.last.prev.next = nil
		l.lru[l.last.prev.key] = l.last.prev
	} else {
		// 如果上个元素是空的， 那么说明这是开头的一个元素
		if l.len == 1 {
			// 如果只有这一个元素， 那么全部清空
			l.root = nil
			l.last = nil
			l.len = 0
			l.lru = make(map[T1]*element[T1, T2])
			return lastKey
		}
		delete(l.lru, lastKey)
		l.len--
	}

	return lastKey
}

func (l *Lru[T1, T2]) moveToPrev(key T1, value T2) {
	// 这里面的元素至少有2个, 否则进不来这里
	// 否则就插入到开头, 开头的元素后移
	//把当前位置元素的上一个元素的下一个元素指向本元素的下一个元素
	//el := &element{}

	if l.len == 2 {
		//如果是2个元素
		//也就是更换元素的值就好了
		//把第一个元素换到第二去

		l.root, l.last = l.last, l.root
		l.last.prev = l.root
		l.last.next = nil
		l.root.next = l.last
		l.root.prev = nil
		l.root.key = key
		l.root.value = value
		l.root.update = time.Now()
		l.lru[key] = l.root
		return
	}
	if l.len > 2 {
		// 拿到这个key的值

		if l.lru[key] == l.last {

			l.last.prev.next = nil
			// 最后一个元素 是最后一个元素
			l.last = l.last.prev
			l.lru[l.last.key] = l.last
		}
		//如果不是, 更新这个元素 上一个和下一个元素的值
		l.lru[key].prev.next = l.lru[key].next
		l.lru[key].next.prev = l.lru[key].prev
		//抽出来这个值到开头
		l.lru[key].prev = nil
		l.lru[key].update = time.Now()
		l.lru[key].value = value
		l.lru[key].next = l.root
		// tmp 是第二个元素
		tmp := l.root
		l.root = l.lru[key]

		// 更新 第二个元素
		tmp.prev = l.root
		//更新第二个元素的Lru
		l.lru[tmp.key] = tmp

	}
}

func (l *Lru[T1, T2]) FirstKey() any {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.root.key
}

func (l *Lru[T1, T2]) LastKey() T1 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.last.key
}

func (l *Lru[T1, T2]) Clean(n int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l = nil
	l.lru = nil
	l = &Lru[T1, T2]{
		lru:  make(map[T1]*element[T1, T2], 0),
		len:  0,
		size: n,
		lock: sync.RWMutex{},
		root: &element[T1, T2]{},
		last: &element[T1, T2]{},
	}
}

func (l *Lru[T1, T2]) Exsit(key T1) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.lru[key]; ok {
		return true
	}
	return false
}
