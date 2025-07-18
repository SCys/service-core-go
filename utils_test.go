package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func IsZeroA[T comparable](v T) bool {
	var zero T
	return v == zero
}

func IsZeroB[T uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64](v T) bool {
	fmt.Printf("%T %v %T %v\n", v, v, 0, 0)
	return v == 0
}

func IsZeroD(v interface{}) bool {
	if v == nil {
		return true
	}
	switch k := v.(type) {
	case int, int8, int32, int16, int64:

		if k == 0 {
			return true
		}

		fmt.Printf("%T %T\n", k, 0)
	case uint, uint8, uint32, uint16, uint64:

		if k == 0 {
			return true
		}

		fmt.Printf("%T %T\n", k, 0)
	default:
		fmt.Printf("unknown type ?? %T %v\n", v, v)
	}

	return false
}

func IsZeroE(v interface{}) bool {
	return v == 0
}

func IsZeroF(v interface{}) bool {
	if v == nil {
		return true
	}

	switch v := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v == 0
	}

	return false
}

func TestZeroMethods(t *testing.T) {
	// 测试泛型IsZero函数
	assert.True(t, IsZero(uint8(0)), "uint8 eq method A?")
	assert.True(t, IsZero(uint16(0)), "uint16 eq method A?")
	assert.True(t, IsZero(uint32(0)), "uint32 eq method A?")
	assert.True(t, IsZero(uint64(0)), "uint64 eq method A?")
	assert.True(t, IsZero(int8(0)), "int8 eq method A?")
	assert.True(t, IsZero(int16(0)), "int16 eq method A?")
	assert.True(t, IsZero(int32(0)), "int32 eq method A?")
	assert.True(t, IsZero(int64(0)), "int64 eq method A?")

	// 测试IsZeroValue函数
	assert.True(t, IsZeroValue(uint8(0)), "uint8 eq method IsZeroValue?")
	assert.True(t, IsZeroValue(uint16(0)), "uint16 eq method IsZeroValue?")
	assert.True(t, IsZeroValue(uint32(0)), "uint32 eq method IsZeroValue?")
	assert.True(t, IsZeroValue(uint64(0)), "uint64 eq method IsZeroValue?")
	assert.True(t, IsZeroValue(int8(0)), "int8 eq method IsZeroValue?")
	assert.True(t, IsZeroValue(int16(0)), "int16 eq method IsZeroValue?")
	assert.True(t, IsZeroValue(int32(0)), "int32 eq method IsZeroValue?")
	assert.True(t, IsZeroValue(int64(0)), "int64 eq method IsZeroValue?")

	// 测试其他类型的零值
	assert.True(t, IsZeroValue(""), "string empty eq method IsZeroValue?")
	assert.True(t, IsZeroValue(false), "bool false eq method IsZeroValue?")
	assert.True(t, IsZeroValue(0.0), "float64 zero eq method IsZeroValue?")
	assert.True(t, IsZeroValue(nil), "nil eq method IsZeroValue?")

	// 测试非零值
	assert.False(t, IsZeroValue(uint8(1)), "uint8 non-zero eq method IsZeroValue?")
	assert.False(t, IsZeroValue("hello"), "string non-empty eq method IsZeroValue?")
	assert.False(t, IsZeroValue(true), "bool true eq method IsZeroValue?")
}
