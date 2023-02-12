package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spade69/xxxcache/group"
)

func TestGetter(t *testing.T) {
	var f group.Getter = group.GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("call back failed")
	}
}

func TestGetFromGroup(t *testing.T) {
	var mock_db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
	// if cache is empty , use callback to get data source
	// if cache exist or hit, just retrieve from cache,
	loadCounts := make(map[string]int, len(mock_db))
	// 2<<10 --> 2048
	g := group.NewGroup("scores", 2<<10, group.GetterFunc(
		func(key string) ([]byte, error) {
			t.Logf("slow db search key %s", key)
			if v, ok := mock_db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	for k, v := range mock_db {
		if view, err := g.Get(k); err != nil || (*view).String() != v {
			t.Fatalf("failed to get value of Tom")
		} // load from callback function
		fmt.Println("print counts", loadCounts)
		if _, err := g.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("value of unknown should be empty but %s got", view)
	}
}
