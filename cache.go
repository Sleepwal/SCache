package main

import (
	"SleepCache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

// add
// @Description: 封装 add 方法，并添加互斥锁 mu。
// @receiver c
// @param key
// @param value
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock() // 锁
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// get
// @Description: 封装 get 方法，并添加互斥锁 mu。
// @receiver c
// @param key
// @return value
// @return ok
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
