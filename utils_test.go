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
	assert.True(t, IsZeroA(uint8(0)), "uint8 eq method A?")
	assert.True(t, IsZeroA(uint16(0)), "uint16 eq method A?")
	assert.True(t, IsZeroA(uint32(0)), "uint32 eq method A?")
	assert.True(t, IsZeroA(uint64(0)), "uint64 eq method A?")
	assert.True(t, IsZeroA(int8(0)), "int8 eq method A?")
	assert.True(t, IsZeroA(int16(0)), "int16 eq method A?")
	assert.True(t, IsZeroA(int32(0)), "int32 eq method A?")
	assert.True(t, IsZeroA(int64(0)), "int64 eq method A?")

	assert.True(t, IsZeroB(uint8(0)), "uint8 eq method B?")
	assert.True(t, IsZeroB(uint16(0)), "uint16 eq method B?")
	assert.True(t, IsZeroB(uint32(0)), "uint32 eq method B?")
	assert.True(t, IsZeroB(uint64(0)), "uint64 eq method B?")
	assert.True(t, IsZeroB(int8(0)), "int8 eq method B?")
	assert.True(t, IsZeroB(int16(0)), "int16 eq method B?")
	assert.True(t, IsZeroB(int32(0)), "int32 eq method B?")
	assert.True(t, IsZeroB(int64(0)), "int64 eq method B?")

	assert.True(t, IsZeroD(int(0)), "int eq method D?")
	assert.True(t, IsZeroD(uint8(0)), "uint8 eq method D?")
	assert.True(t, IsZeroD(uint16(0)), "uint16 eq method D?")
	assert.True(t, IsZeroD(uint32(0)), "uint32 eq method D?")
	assert.True(t, IsZeroD(uint64(0)), "uint64 eq method D?")
	assert.True(t, IsZeroD(int8(0)), "int8 eq method D?")
	assert.True(t, IsZeroD(int16(0)), "int16 eq method D?")
	assert.True(t, IsZeroD(int32(0)), "int32 eq method D?")
	assert.True(t, IsZeroD(int64(0)), "int64 eq method D?")

	assert.True(t, IsZeroE(int(0)), "int eq method E?")
	assert.True(t, IsZeroE(uint8(0)), "uint8 eq method E?")
	assert.True(t, IsZeroE(uint16(0)), "uint16 eq method E?")
	assert.True(t, IsZeroE(uint32(0)), "uint32 eq method E?")
	assert.True(t, IsZeroE(uint64(0)), "uint64 eq method E?")
	assert.True(t, IsZeroE(int8(0)), "int8 eq method E?")
	assert.True(t, IsZeroE(int16(0)), "int16 eq method E?")
	assert.True(t, IsZeroE(int32(0)), "int32 eq method E?")
	assert.True(t, IsZeroE(int64(0)), "int64 eq method E?")

	assert.True(t, IsZeroF(int(0)), "int eq method F?")
	assert.True(t, IsZeroF(uint8(0)), "uint8 eq method F?")
	assert.True(t, IsZeroF(uint16(0)), "uint16 eq method F?")
	assert.True(t, IsZeroF(uint32(0)), "uint32 eq method F?")
	assert.True(t, IsZeroF(uint64(0)), "uint64 eq method F?")
	assert.True(t, IsZeroF(int8(0)), "int8 eq method F?")
	assert.True(t, IsZeroF(int16(0)), "int16 eq method F?")
	assert.True(t, IsZeroF(int32(0)), "int32 eq method F?")
	assert.True(t, IsZeroF(int64(0)), "int64 eq method F?")
}
