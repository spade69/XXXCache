package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash function ,using dependecy injection,
// allow sustitute to self define Hash;
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	// Hash function, passing by outer source
	hash Hash
	// 虚拟节点倍数
	replicas int
	// hash ring
	keys []int
	// vNode --> Real Node
	hashMap map[int]string
}

// New Creates a map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		// default is ChecksumIEEE
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// add some keys to hash, allow passing 0 or multiple key
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// one key --> m.replicas virtual node.
		for i := 0; i < m.replicas; i++ {
			// vnode = i+ key to means nvnode keys,
			// m.hash() to calculate vnode's val
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			// add vnode hash to hashmap
			m.hashMap[hash] = key
		}

	}
	// sort ring
	sort.Ints(m.keys)
}

// 1. calcuate hash of key
// 2. using binary search to find a vnode -->idx
// 3. find a node using m.hashmap
// 4. m
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
