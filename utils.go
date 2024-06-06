package core

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/valyala/fastjson"
)

// String copy from strings.Builder
func String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// Bytes convert str to bytes
func Bytes(str string) []byte {
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
}

// Now from fasttime
func Now() time.Time {
	return time.Now().In(time.Local)
}

// OffsetAndLimit get offset,limit in json params
func OffsetAndLimit(values *fastjson.Value) (int64, int) {
	offset := values.GetInt64("params", "offset")
	limit := values.GetInt("params", "limit")

	if limit < 0 {
		limit = 0
	} else if limit > 1000 {
		limit = 1000
	}

	if offset < 0 {
		offset = 0
	}

	return offset, limit
}

// IsZero compare value with zero
func IsZero[T comparable](v T) bool {
	var zero T // init with value (0?)
	return v == zero
}
