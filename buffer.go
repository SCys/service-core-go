package core

import (
	"bytes"
	"io"
	"sync"
)

// CopyBufPool cache pool
var CopyBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// CopyZeroAlloc copy from fasthttp.http.copyZeroAlloc
func CopyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	buf := CopyBufPool.Get().(*bytes.Buffer)
	buf.Reset()

	n, err := io.CopyBuffer(w, r, buf.Bytes())
	CopyBufPool.Put(buf)
	return n, err
}
