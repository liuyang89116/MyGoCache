package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Bob":  "500",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("Callback failed!")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	geeCache := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[callback - SlowDB] search key:", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s doesn't exist", key)
		},
	))

	for k, v := range db {
		if bv, err := geeCache.Get(k); err != nil || bv.String() != v {
			t.Fatalf("Failed to get value of key: %s\n", k)
		}
		if _, err := geeCache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss\n", k)
		}
	}

	if bv, err := geeCache.Get("unknown"); err == nil {
		t.Fatalf("The value of unknow should be empty, but got %s\n", bv.String())
	}
}

func TestGetGroup(t *testing.T) {
	groupName := "TestGroup"
	NewGroup(groupName, 2<<10, GetterFunc(func(key string) (bytes []byte, err error) {
		return
	}))

	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("group %s doesn't exist\n", groupName)
	}
	if group := GetGroup("unknown_group"); group != nil {
		t.Fatalf("expect empty, but got %s\n", group.name)
	}
}
