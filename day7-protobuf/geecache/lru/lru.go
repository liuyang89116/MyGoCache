package lru

import "container/list"

// Cache is an LRU cache. It's not concurrent safe
type Cache struct {
	maxBytes int64 // max allowed bytes
	nbytes   int64 // already used bytes
	ll       *list.List
	cache    map[string]*list.Element // key - string, value - list pointer
	// optional: executed when an entry is purged
	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len() to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the constructor for Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(), // list 创建出来就是指针形式
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// Get looks up a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry) // 这里的 Value 是 Element 下面的 Value。我们把它转成 entry 的指针类
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest node
func (c *Cache) RemoveOldest() {
	elem := c.ll.Back()
	if elem != nil {
		c.ll.Remove(elem)
		kv := elem.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

// Add new key value pair
func (c *Cache) Add(key string, value Value) {
	if elem, ok := c.cache[key]; ok { // 如果有这个 key
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry) // 类型转换，转成 *entry
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		elem := c.ll.PushFront(&entry{key, value})
		c.cache[key] = elem
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.nbytes > c.maxBytes { // 可能需要 remove 多个 old entries，直到小于 max
		c.RemoveOldest()
	}
}

// Test purposes only: implement len()
func (c *Cache) Len() int {
	return c.ll.Len()
}
