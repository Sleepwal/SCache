package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBaseURL = "/api/"

type Data struct {
	Message string `json:"message"`
}

type HttpPool struct {
	self     string // 地址
	basePath string // 前缀
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
