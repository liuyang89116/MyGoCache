package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	hash     Hash
	replicas int
	keys     []int          // sorted keys; keys represent different nodes
	hashmap  map[int]string // <hash - real node> pair
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashmap:  make(map[int]string),
	}
	if m.hash == nil { // 默认一个 hash 函数
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add adds some keys to the hash; each key represents a real node
// also for each key, creates a fixed number of virtual replicas
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashmap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// binary search for a replica
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashmap[m.keys[idx%len(m.keys)]] // 环形结构。如果 idx == len(m.keys), 那么应该选择 m.keys[0]
}
