package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	// protects m
	mu sync.Mutex
	m  map[string]*call
}

// key :request key in cache, fn : only call once and return no matter ho many time do
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// g.mu is protecting Group member m -> is map for every key
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// if request  processing then wait
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	// request and lock
	c.wg.Add(1)
	// mapping for key --> c, means already get request to process
	g.m[key] = c
	g.mu.Unlock()
	// assign to call obj, call fn
	c.val, c.err = fn()
	c.wg.Done()
	// lock
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
