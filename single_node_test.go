package main

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	want := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, want) {
		t.Errorf("Get(%q) = %v, want %v", "key", v, want)
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db)) // 统计命中次数
	group := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key: ", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok { // 首次从db中取值
					loadCounts[key] = 0
				}
				loadCounts[key] += 1 // db命中次数+1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		// load from callback function
		if view, err := group.Get(k); err != nil || view.String() != v { // 第一次从db中取值，会存储在缓存中
			t.Fatalf("group.Get(%q) failed:%v", k, err)
		}
		if _, err := group.Get(k); err != nil || loadCounts[k] > 1 { // 从缓存中取值，故只命中db一次
			t.Fatalf("cache %q miss", k)
		}
	}

	if view, err := group.Get("unknown"); err == nil { // db中不存在
		t.Fatalf("expected to get error, but didn't, value:%s", view)
	}
}
