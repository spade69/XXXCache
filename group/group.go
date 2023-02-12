package group

import (
	"fmt"
	"sync"

	"github.com/spade69/xxxcache/core"
)

// A Group means a cache namespace,every group own a unique name
// eg: StudentGroup,InfoGroup, ....
type Group struct {
	name   string
	getter Getter
	// concurrent safe cache
	mainCache core.Scache
}

// Getter is a interface used to get data from different datasource
// eg: mysql, redis, pgsql
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function
type GetterFunc func(key string) ([]byte, error)

// Implement Getter method Get, as a Getter function?
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	// assign to group map, store at mainCache
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: core.Scache{
			CacheBytes: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

// GetGroup用来特定名称的 Group，
// 这里使用了只读锁 RLock()，因为不涉及任何冲突变量的写操作。
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (*core.ByteView, error) {
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.Get(key); ok {
		fmt.Println("cache hit!")
		return v, nil
	}
	// load key from local if fail
	return g.Load(key)
}

func (g *Group) Load(key string) (*core.ByteView, error) {
	return g.GetLocally(key)
}

// GetLocally Get data from user define data source
// and set data into mainCache(by populateCache method)
func (g *Group) GetLocally(key string) (*core.ByteView, error) {
	// Get data from local datasource
	bytes, err := g.getter.Get(key)
	if err != nil {
		fmt.Println("Get data from datasource fail", err)
		return nil, err
	}
	value := core.NewByteView(bytes)
	g.PopulateCache(key, value)
	return &value, nil
}

func (g *Group) PopulateCache(key string, value core.ByteView) {
	g.mainCache.Set(key, value)
}
