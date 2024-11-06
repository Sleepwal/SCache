package consistent_hash

import (
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i) // 返回字符串对应的数字
	})

	hash.Add("6", "4", "2") // 2、4、6三个真实节点
	// 2 对应虚拟节点-- 2、12、22
	// 4 -- 2、14、24
	// 6 -- 6、16、26
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s, but get %s", k, hash.Get(k), v)
		}
	}

	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s, but get %s", k, hash.Get(k), v)
		}
	}
}
