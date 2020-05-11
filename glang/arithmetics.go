package glang

import (
	"encoding/binary"
)

//TODO: TREAT OVERFLOW ERRORS

func IAdd(a, b []byte) []byte {
	var left int64
	var right int64
	res := make([]byte, 8)
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	binary.PutVarint(res, left+right)
	return res
}

func ISub(a, b []byte) []byte {
	var left int64
	var right int64
	res := make([]byte, 8)
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	binary.PutVarint(res, left-right)
	return res
}

func IMult(a, b []byte) []byte {
	var left int64
	var right int64
	res := make([]byte, 8)
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	binary.PutVarint(res, left*right)
	return res
}

func IDiv(a, b []byte) []byte {
	var left int64
	var right int64
	res := make([]byte, 8)
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	binary.PutVarint(res, left/right)
	return res
}

func IMod(a, b []byte) []byte {
	var left int64
	var right int64
	res := make([]byte, 8)
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	binary.PutVarint(res, left%right)
	return res
}
