package main

import (
	"SleepCache/singleflight"
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error) // 回调函数
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 负责与用户的交互，并且控制缓存值存储和获取的流程。
type Group struct {
	name       string
	getter     Getter
	mainCache  cache
	peerPicker PeerPicker
	loader     *singleflight.Group // 确保 同一时间一个key只被一个请求获取
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup
// @Description: 创建一个group
// @param name
// @param cacheBytes
// @param getter
// @return *Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
		loader: &singleflight.Group{},
	}

	groups[name] = g
	return g
}

// GetGroup
// @Description: 根据name查找group
// @param name
// @return *Group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 从 mainCache 中查找缓存
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[SCache] hit")
		return v, nil
	}
	// 没有命中，调用 load 方法
	return g.load(key)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peerPicker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peerPicker = peers
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 样确保了并发场景下针对相同的 key，load 过程只会调用一次
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peerPicker != nil {
			if peer, ok := g.peerPicker.PickPeer(key); ok { // 根据key选择一个节点
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[SCache] Failed to get form peer, ", err)
			}
		}
		return g.getLocally(key)
	})

	if err == nil {
		return view.(ByteView), nil
	}
	return
}

// 从节点中获取值
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用用户回调函数 g.getter.Get() 获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	// 将源数据存储到 mainCache 中
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
