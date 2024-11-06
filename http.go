package main

import (
	"SleepCache/consistent_hash"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultBaseURL  = "/api/"
	defaultReplicas = 50
)

type Data struct {
	Message string `json:"message"`
}

type HttpPool struct {
	self        string // 自己的地址，包括主机名/IP 和端口
	basePath    string // 节点间通讯地址的前缀
	mu          sync.Mutex
	peers       *consistent_hash.Map   // 存储节点，使用一致性哈希算法选择节点
	httpGetters map[string]*httpGetter // 每一个远程节点对应一个 httpGetter, key值示例 "http://10.0.0.2:8008"
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBaseURL,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s\n", p.self, fmt.Sprintf(format, v...))
}

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("unexpected path: " + r.URL.Path)
	}
	p.Log("HTTP request %s - %s", r.URL.Path, p.basePath)
	// /<base_path>/<group_name>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2) // 分隔
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0] // 取出group name
	key := parts[1]       // 取出key

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := Data{Message: string(view.ByteSlice())}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonData)
}

// Set 更新节点列表
func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 初始化map，添加节点
	p.peers = consistent_hash.New(defaultReplicas, nil)
	p.peers.Add(peers...)

	// 初始化HTTP客户端
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 根据key选择一个节点，返回HTTP客户端
func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	peer := p.peers.Get(key)
	if peer != "" && peer != p.self {
		return p.httpGetters[peer], true
	}
	return nil, false
}

// 用于确保类型 HttpPool 实现了 PeerPicker 接口，编译时检查
var _ PeerPicker = (*HttpPool)(nil)
