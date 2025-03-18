package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Number interface {
	uint16 | uint32 | uint64
}

func ToLittleEndianGen[T Number](number T) T {
	octetsCount := unsafe.Sizeof(number)
	var result T
	pointer := unsafe.Pointer(&number)
	var idx uintptr

	for idx = 0; idx < octetsCount; idx++ {
		byteValue := *(*uint8)(unsafe.Add(pointer, idx))
		bitOffset := (octetsCount - idx - 1) * 8
		result |= T(byteValue) << bitOffset
	}

	return result
}

func TestСonversionGen32(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGen(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestСonversionGen16(t *testing.T) {
	tests := map[string]struct {
		number uint16
		result uint16
	}{
		"test case #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"test case #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case #3": {
			number: 0x00FF,
			result: 0xFF00,
		},
		"test case #4": {
			number: 0x0FFF,
			result: 0xFF0F,
		},
		"test case #5": {
			number: 0x0102,
			result: 0x0201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGen(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestСonversionGen64(t *testing.T) {
	tests := map[string]struct {
		number uint64
		result uint64
	}{
		"test case #1": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"test case #2": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF00FF00FF,
			result: 0xFF00FF00FF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF0000FFFF,
			result: 0xFFFF0000FFFF0000,
		},
		"test case #5": {
			number: 0x0102030405060708,
			result: 0x0807060504030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianGen(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
