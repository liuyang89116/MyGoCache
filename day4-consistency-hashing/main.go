package main

import (
	"day4-consistency-hashing/geecache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "821",
}

/**
简单测试：

$ curl http://localhost:9999/_geecache/scores/Tom
630

$ curl http://localhost:9999/_geecache/scores/kkk
kkk not exist
*/
func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s doesn't exits\n", key)
		},
	))

	addr := "localhost:9999"
	httpPool := geecache.NewHttpPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, httpPool))
}

/**
运行结果：

2022/05/15 23:56:47 geecache is running at localhost:9999
2022/05/15 23:57:19 [Server localhost:9999] GET /_geecache/scores/Tom/n
2022/05/15 23:57:19 [SlowDB] search key Tom
2022/05/15 23:57:34 [Server localhost:9999] GET /_geecache/scores/TTom/n
2022/05/15 23:57:34 [SlowDB] search key TTom
2022/05/15 23:57:39 [Server localhost:9999] GET /_geecache/scores/Tom/n
2022/05/15 23:57:39 [GeeCache] hit
2022/05/15 23:57:50 [Server localhost:9999] GET /_geecache/scores/Tom1/n
2022/05/15 23:57:50 [SlowDB] search key Tom1
*/
