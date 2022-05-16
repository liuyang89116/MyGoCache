package geecache

// A ByteView holds an immutable view of bytes
type ByteView struct {
	b []byte
}

// Len returns the length of a ByteView
func (bv ByteView) Len() int {
	return len(bv.b)
}

// ByteSlice returns a copy of the data as a byte slice
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

func (bv ByteView) String() string {
	return string(bv.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b) // copy slice: func copy(dst, src []Type) int
	return c
}
