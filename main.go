package main

import (
	"fmt"
	"hash/crc32"
)

func main() {
	fmt.Println("print cache here")
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
