package test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var mutex sync.Mutex
var set = make(map[int]bool, 0)

func printNum(num int) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true // 已打印过
}

func TestMutex(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printNum(100) //相同的数字只会被打印一次
	}
	time.Sleep(1 * time.Second)
}
