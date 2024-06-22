package main

import (
	"SleepCache/sleepcache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	sleepcache.NewGroup("scores", 2<<10, sleepcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	address := "localhost:9999"
	handler := sleepcache.NewHttpPool(address)
	log.Println("SleepCache is running at", address)
	log.Fatal(http.ListenAndServe(address, handler))
}
