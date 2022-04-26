package main

import "github.com/hyahm/cache"

func main() {
	c := cache.NewCache(100, cache.LFU)
	c.Add(1, 2)
	c.Add(4, 12)
	c.OrderPrint(0)
}
