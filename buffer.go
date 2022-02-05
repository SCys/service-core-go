package core

import (
	"io"

	"github.com/valyala/bytebufferpool"
)

func Copy(w io.Writer, r io.Reader) (int64, error) {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)

	return io.CopyBuffer(w, r, bb.Bytes())
}
