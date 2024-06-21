package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type TestString string

func (t TestString) Len() int {
	return len(t)
}

// TestGet
// @Description: 测试LRU缓存是否能够正确获取数据
// @param t
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", TestString("value1"))
	lru.Add("key2", TestString("value2"))
	if v, ok := lru.Get("key1"); ok {
		if string(v.(TestString)) != "value1" {
			t.Fatalf("got %s, want value1", v)
		}
	} else {
		t.Fatalf("key1 not found")
	}

	if v, ok := lru.Get("key2"); ok {
		if string(v.(TestString)) != "value2" {
			t.Fatalf("got %s, want value2", v)
		}
	} else {
		t.Fatalf("key2 not found")
	}
}

// TestRemoveOldest
// @Description: 使用内存超过了设定值时，是否会触发“无用”节点的移除
// @param t
func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	maxCap := len(k1 + k2 + v1 + v2)
	fmt.Println("max capity:", maxCap)
	lru := New(int64(maxCap), nil)
	lru.Add(k1, TestString(v1))
	lru.Add(k2, TestString(v2))
	lru.Add(k3, TestString(v3))

	if v, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("key1 should be evicted")
	} else {
		fmt.Println("value: ", v, " len: ", lru.Len())
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) { // 回调函数
		keys = append(keys, key) // 添加在keys中
	}

	lru := New(int64(10), callback)
	lru.Add("key1", TestString("value1")) // 加10字节，剩0字节
	lru.Add("k2", TestString("v2"))       // 加4字节，淘汰key1，剩6字节
	lru.Add("k3", TestString("v3"))       // 加4字节，剩2字节，不淘汰
	lru.Add("k4", TestString("v4"))       // 加4字节，共12字节，淘汰k2，省8字节

	want := []string{"key1", "k2"} // key1、k2被淘汰

	if !reflect.DeepEqual(want, keys) {
		fmt.Println("keys: ", keys)
		t.Fatalf("Call OnEvicted function, but got different keys")
	}
}
