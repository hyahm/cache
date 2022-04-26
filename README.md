# 缓存算法（lru, lfu, alfu）
 线程安全的go语言通用lru, lfu算法
### 安装
```
go get github.com/hyahm/lru
```
### 使用

在使用前, 先要初始化缓存个数, 内存足够大的话, 可以设置很大, 不会因为长度影响效率, 
超过设定值会自动删除末尾的值, 如果存在的话会自动更新此值到开头, 更新值
 > 初始化(初始化完成后, 可以在任何地方调用方法)
  ```
  lru.Init(10)
  ```
 > 添加 key和value, 以下为doc中的例子
  ```
package main

import (
	"fmt"
	"lru"
)

type el struct {
	Id int
	Name string
}

func main() {
	cache = cache.New(10, cache.LRU)

	cache.Add("adsf", "bbbbb")
	cache.Add("cccc", "111111")
	e := &el{
		Id: 1,
		Name: "68",
	}
	cache.Add("adsf", e)
	fmt.Println(lru.Len())
	cache.OrderPrint(0)
}
```
> 万能的add方法, 只要是添加值都可以使用此方法, 存在就会更新, 不存在就会插入, 返回删除的key, 没删除返回nil
```
cache.Add(key, value interface{}) interface{}
```
> 顺序打印(有读写锁, 会阻碍读写操作,正式环境建议别使用)
```
cache.OrderPrint()
```
> 无序打印(查看缓存, 推荐使用)
```
cache.Print()
```
> 删除key
```
cache.Remove(key interface{})
```
> 获取所有的key, 没有就返回空, 返回的key因为执行时间的问题, 可能导致有些key被删除了
```
cache.Keys(key interface{}) []interface{}
```
> 获取缓存长度 
```
cache.Len() uint64
```
> 根据key获取值
```
cache.Get(key interface{}) interface{}
```
> 根据key获取上一个key
```
cache.Prev(key interface{}) interface{}
```
> 根据key获取下一个key
```
cache.Next(key interface{}) interface{}
```
> 判断是否存在key
```
cache.Exsit(key interface{}) bool
```
> 重新设置缓存的长度
```
cache.Resize(n uint64)
```
> 清空缓存(不推荐使用, 也未测试)
```
cache.Clean(n)
```
基本上这些就能满足需求
