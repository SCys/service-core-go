package core

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/valyala/fastjson"
)

// String 将字节切片转换为字符串
// 注意：这个函数使用了unsafe转换，仅适用于临时使用，不要用于长期存储
func String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// Bytes 将字符串转换为字节切片
// 注意：这个函数使用了unsafe转换，仅适用于临时使用，不要用于长期存储
// 建议使用 []byte(str) 进行安全转换
func Bytes(str string) []byte {
	// 使用更安全的方式
	return []byte(str)
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

// IsZero 比较值与零值
// 支持所有可比较类型，包括基本类型和自定义类型
func IsZero[T comparable](v T) bool {
	var zero T // 零值初始化
	return v == zero
}

// IsZeroValue 通用的零值检查函数，支持interface{}类型
func IsZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case int:
		return val == 0
	case int8:
		return val == 0
	case int16:
		return val == 0
	case int32:
		return val == 0
	case int64:
		return val == 0
	case uint:
		return val == 0
	case uint8:
		return val == 0
	case uint16:
		return val == 0
	case uint32:
		return val == 0
	case uint64:
		return val == 0
	case float32:
		return val == 0
	case float64:
		return val == 0
	case string:
		return val == ""
	case bool:
		return !val
	default:
		// 对于其他类型，使用反射检查
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			return rv.IsNil()
		case reflect.Slice, reflect.Map, reflect.Chan:
			return rv.Len() == 0
		default:
			return reflect.DeepEqual(v, reflect.Zero(rv.Type()).Interface())
		}
	}
}
