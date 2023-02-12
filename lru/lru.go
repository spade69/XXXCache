package lru

import (
	"container/list"
)

// LRU Cache, not safe for concurrent access.
// LRU = hashmap + dual
type Cache struct {
	// max memory in bytes
	maxBytes int64
	// used memory
	nBytes int64
	// link list fron golang container
	ll *list.List
	// hash map key is string,  value is list element
	cache map[string]*list.Element
	// optiinal , evict key callback func
	OnEvicted func(key string, value Value)
}

// Entry is dual link list's element,
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		// initializa list
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Set key if exist
func (c *Cache) Set(key string, value Value) {
	// 1. if key exist, update
	if ele, ok := c.cache[key]; ok {
		// move ele to front which means it is latest
		c.ll.MoveToFront(ele)
		// get origin entry ele
		originEle := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(originEle.value.Len())
		// update value of ele
		originEle.value = value
	} else {
		// 2. else if key not exist
		// first push into linklist, then return ele
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key) + value.Len())
	}
	// 3.  remove old keys if excceed max bytes
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.EvictKey()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// access more frequently
		c.ll.MoveToFront(ele)
		// list elemtn is entry, entry's value is result
		en := ele.Value.(*entry)
		return en.value, true
	}
	return nil, false
}

// Evict key which access less than usual
func (c *Cache) EvictKey() {
	// last element need to be evict
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		en := ele.Value.(*entry)
		delete(c.cache, en.key)
		// release bytes
		c.nBytes -= int64(len(en.key) + en.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(en.key, en.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
