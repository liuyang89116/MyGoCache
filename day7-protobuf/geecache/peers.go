package geecache

// PeerPicker is the interface that must be implemented to locate
// the peer that owns the specific key.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
// httpGetter under http implemented it.
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
