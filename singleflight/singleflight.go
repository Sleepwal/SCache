package singleflight

import "sync"

// 代表正在进行中，或已经结束的请求
type call struct {
	wg  sync.WaitGroup // 使用 sync.WaitGroup 锁避免重入
	val interface{}
	err error
}

// Group 管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex // 保护 Group 的成员变量 m 不被并发读写而加上的锁
	m  map[string]*call
}

// Do 针对相同的 key，无论 Do 被调用多少次，函数 fn 都只会被调用一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果请求正在进行中，则等待
		return c.val, c.err // 请求结束，返回结果
	}

	c := &call{}
	c.wg.Add(1)  // 发起请求前加锁
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 删除请求
	g.mu.Unlock()

	return c.val, c.err
}
