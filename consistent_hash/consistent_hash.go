package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // Hash 函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环，需要排序，放的是虚拟节点
	hashMap  map[int]string // 虚拟与真实节点的映射表
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys { // 每一个真实节点 key，创建 m.replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) // 添加编号i，区分节点
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key))) // 计算hash
	// 找到第一个匹配的虚拟节点的下标 idx
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 从m.keys中获取哈希值，hashMap映射得到真实的节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
