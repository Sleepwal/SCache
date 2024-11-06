package main

// ByteView 抽象的只读数据，表示缓存值
type ByteView struct {
	b []byte // 存储真实的缓存值，byte可以支持任意类型。
}

// Len 返回其所占的内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b) // 返回拷贝，防止被修改
}

func (v ByteView) String() string {
	return string(v.b)
}

// 返回拷贝，防止被修改
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
