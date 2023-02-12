package test

import (
	"testing"

	"github.com/spade69/xxxcache/lru"
)

//
type String string

func (d String) Len() int {
	return len(d)
}
func TestGet(t *testing.T) {
	c := lru.New(100, nil)
	c.Set("key1", String("1233"))
	v, ok := c.Get("key1")
	if !ok || string(v.(String)) != "1233" {
		t.Fatalf("cache get key1 fail")
	}
	if _, ok := c.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestEvict(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	c := lru.New(int64(cap), nil)
	c.Set(k1, String(v1))
	c.Set(k2, String(v2))
	c.Set(k3, String(v3))
	if _, ok := c.Get(k1); ok || c.Len() != 2 {
		t.Fatalf("Remove oldest key fail")
	}
}
