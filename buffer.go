package core

import (
	"bytes"
	"io"
	"sync"
)

// CopyBuffers cache pool
var CopyBuffers = sync.Pool{
	New: func() interface{} {
		// generate 1MB buffer
		return bytes.NewBuffer(make([]byte, 0, 2<<19))
	},
}

func Copy(w io.Writer, r io.Reader) (int64, error) {
	buf := CopyBuffers.Get().(*bytes.Buffer)
	defer CopyBuffers.Put(buf)

	buf.Grow(2 << 20) // 2MB
	b := buf.Bytes()

	return io.CopyBuffer(w, r, b[:buf.Cap()])
}
