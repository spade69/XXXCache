package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"

	"github.com/spade69/xxxcache/group"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *group.Group {
	return group.NewGroup("scores", 2<<10, group.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// startCacheServer() 用来启动缓存服务器：
// 创建 HTTPPool，添加节点信息，注册到 g 中，启动 HTTP 服务（共3个端口，8001/8002/8003），用户不感知。
func startCacheServer(addr string, addrs []string, g *group.Group) {
	// self
	peers := group.NewHTTPPool(addr)
	// set peers
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("xxxcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// api server 9999, used for interactive with client
func startAPIServer(apiAddr string, g *group.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := g.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	// start server
	var port int
	var api bool
	// command line input
	flag.IntVar(&port, "port", 8001, "XXXCAche server port")
	// command lint input
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	// api
	apiAddr := "http://127.0.0.1:9999"
	addrMap := map[int]string{
		8001: "http://127.0.0.1:8001",
		8002: "http://127.0.0.1:8002",
		8003: "http://127.0.0.1:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	g := createGroup()
	// api is true
	if api {
		go startAPIServer(apiAddr, g)
	}
	startCacheServer(addrMap[port], []string(addrs), g)
	// startServer()

}

func testHash() {
	key := "test hash"
	key2 := "test3 hash"
	node := simplehash(key)
	node2 := simplehash(key2)
	fmt.Println("node is ", node, "node2 is", node2)
}

func simplehash(key string) uint32 {
	hval := crc32.ChecksumIEEE([]byte(key))
	fmt.Println("hval is", hval)
	node := hval % 10
	return node
}

func startLocalServer() {
	addr := "127.0.0.1:9999"
	peers := group.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	// HTTPPool implements ServeHTTP interface --> handler interface
	// func ListenAndServe(addr string, handler Handler) error {
	//	server := &Server{Addr: addr, Handler: handler}
	//return server.ListenAndServe()
	log.Fatal(http.ListenAndServe(addr, peers))
}
