package core

// A ByteView holds an immutable view of bytes.
type ByteView struct {
	///	// must support arbitrary data store
	b []byte
}

func NewByteView(b []byte) ByteView {
	return ByteView{
		b: b,
	}
}

func (v ByteView) Len() int {
	// impolemtn Len() , because lru.Cache require cache object
	// must impolemtn Value interface,
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string,
// making a copy if necessary., b is read-only
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
