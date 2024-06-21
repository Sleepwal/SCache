package lru

import "container/list"

type Cache struct {
	maxBytes  int64                    // 最大内存
	nBytes    int64                    // 已使用的内存
	linkList  *list.List               // 双向链表
	cache     map[string]*list.Element // 双向链表中对应节点的指针
	OnEvicted func(key string, value Value)
}

// 双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int // 返回值所占用的内存大小
}

func New(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		linkList:  list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}
