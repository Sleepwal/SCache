package main

// PeerPicker 节点选择器
type PeerPicker interface {
	// PickPeer 根据key选择相应的节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 从节点中获取缓存值，HTTP客户端
type PeerGetter interface {
	// Get 从group中查找缓存值
	Get(group string, key string) ([]byte, error)
}
