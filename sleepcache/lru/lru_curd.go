package lru

// Get
// @Description: 根据key查找value
// @receiver c
// @param key
// @return value
// @return ok
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 1.从字典中找到对应的双向链表的节点。
	if element, ok := c.cache[key]; ok {
		// 2.将该节点移动到队尾，在这里约定 front 为队尾
		c.linkList.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest
// @Description: 缓存淘汰。即移除最近最少访问的节点（队首）
// @receiver c
func (c *Cache) RemoveOldest() {
	// 1.取出队首
	element := c.linkList.Back()
	// 2.删除队尾
	if element != nil {
		c.linkList.Remove(element) // 从链表中删除
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)                                // 从map中删除
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新已使用内存，减去key和value所占内存

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) // 调用回调函数
		}
	}
}

// Add
// @Description: 添加或更新一个键值对
// @receiver c
// @param key
// @param value
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok { // 已存在，更新
		// 1. 移动到队尾
		c.linkList.MoveToFront(element)
		// 2. 更新value
		kv := element.Value.(*entry)
		kv.value = value
		// 3.更新内存
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
	} else { // 不存在，添加
		// 1. 创建新节点
		element := c.linkList.PushFront(&entry{key, value})
		// 2. 添加到字典
		c.cache[key] = element
		// 3. 更新内存
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	// 移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.linkList.Len()
}
