package group

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/spade69/xxxcache/communication"
	"github.com/spade69/xxxcache/consistenthash"
)

type server int

// HTTPPool implements PeerPicker for a pool of HTTP peers.

type HTTPPool struct {
	// this peer's base url , eg: http://example.com:8000
	self     string
	basePath string
	// guards peers and httpGetters
	mu sync.Mutex
	// peers is consistenthash map ,used for specific key to select node
	peers *consistenthash.Map
	//  keyed by e.g. "http://10.0.0.2:8008", mapping remote node and corresponse httpGetter
	// one remote endpoint <--> one httpGetter
	httpGetters map[string]*HTTPGetter
}

type HTTPGetter struct {
	baseURL string
}

const (
	defaultBasePath = "/xxxcache/"
	defaultReplicas = 50
)

var _ communication.PeerPicker = (*HTTPPool)(nil)

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self, // own address of itself
		// prefix of node
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 创建任意类型 server，并实现 ServeHTTP 方法。
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.URL.Path)
	// w.Write([]byte("heeloworld"))
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTP pool serving unexpected path:")
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	//
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// groupNmae means api gourp
	groupName := parts[0]
	// get key of this group, group means a cache
	key := parts[1]
	g := GetGroup(groupName)
	if g == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	view, err := g.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

// Set Updates the pool's list of peers
// Set() 方法实例化了一致性哈希算法，并且添加了传入的节点。
//并为每一个节点创建了一个 HTTP 客户端 httpGetter。
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	// make a map for peer <--> httpGetter
	p.httpGetters = make(map[string]*HTTPGetter, len(peers))
	// here we assign map ,write to map ,so using mutex to lock it
	for _, peer := range peers {
		p.httpGetters[peer] = &HTTPGetter{
			baseURL: peer + p.basePath,
		}
	}
}

// 包装了一致性哈希算法的 Get() 方法，根据具体的 key,选择节点，
// 返回节点对应的 HTTP 客户端。
func (p *HTTPPool) PickPeer(key string) (communication.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		log.Printf("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// baseURL --> remote endpoint
// use http.Get to retrieve return value
func (h *HTTPGetter) Get(group, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	// request using http client
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return : %v", res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}
