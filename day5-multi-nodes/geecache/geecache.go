package geecache

import (
	"fmt"
	"log"
	"sync"
)

// A Group is a cache namespace and associates date loading and spreading over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

// A Getter loads data from a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc defines func type to be implemented Get
type GetterFunc func(key string) ([]byte, error)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// Get implements Getter interface
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// NewGroup creates a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter!")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group created previously with NewGroup, or
// nil if there's no such group
func GetGroup(name string) *Group {
	mu.RLock() // 只使用只读锁，因为不涉及任何冲突变量的写操作
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeers called more than once!")
	}

	g.peers = peers
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is empty")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFormPeer(peer, key); err == nil {
				return value, nil
			}
		}
		log.Panicln("[GeeCache] failed to get from peer", err)
	}

	return g.getLocally(key)
}

func (g *Group) getFormPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{b: bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{
		b: cloneBytes(bytes), // 返回的是 cloned bytes，而不是直接把 bytes 返回去
	}
	g.populateCache(key, value)

	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
