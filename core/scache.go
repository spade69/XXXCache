package core

import (
	"sync"

	"github.com/spade69/xxxcache/lru"
)

type Scache struct {
	mu         sync.Mutex
	cache      *lru.Cache
	CacheBytes int64
}

func (s *Scache) Set(key string, value ByteView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// lazy init lru
	if s.cache == nil {
		s.cache = lru.New(s.CacheBytes, nil)
	}
	s.cache.Set(key, value)
}

func (s *Scache) Get(key string) (value *ByteView, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cache == nil {
		return nil, false
	}
	if v, ok := s.cache.Get(key); ok {
		bv := v.(ByteView)
		return &bv, true
	}
	return nil, false
}
