package test

import (
	"strconv"
	"testing"

	"github.com/spade69/xxxcache/consistenthash"
)

func TestHashing(t *testing.T) {
	// replicas 3 means : 1 key --> 3 virtual node
	hash := consistenthash.New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// Given the above hash function, this will give replicas with "hashs":
	// 2, 4, 6, 12, 14, 16,22,24,26..
	// 1. frist add 6, 16, 26
	// 2. second add 4, 14, 24
	// 3. third add 2, 12 ,22
	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":   "2",
		"11":  "2",
		"23":  "4",
		"27":  "2",
		"200": "2",
		// "200": "4",

	}
	t.Log("keys is ", hash.GetKeys())
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Log("hash get is ", hash.Get(k))
			t.Errorf("Asking for %s, should have yield %s", k, v)
		}
	}
	// Adds 8 ,18, 28
	hash.Add("8")
	// 27 should now map to 8
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yield %s", k, v)
		}
	}
}
