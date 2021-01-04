package core

import (
	"github.com/valyala/fastjson"
	"reflect"
	"syscall"
	"time"
	"unsafe"
)

// String copy from strings.Builder
func String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

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
	var tv syscall.Timeval
	err := syscall.Gettimeofday(&tv)
	if err != nil {
		return time.Now().In(time.Local)
	}
	return time.Unix(0, syscall.TimevalToNsec(tv)).In(time.Local)
}

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
