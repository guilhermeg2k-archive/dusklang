package dusk

import (
	"bytes"
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

func FAdd(a, b []byte) []byte {
	var left float64
	var right float64

	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)

	var res bytes.Buffer
	binary.Write(&res, binary.LittleEndian, left+right)
	return res.Bytes()
}

func FSub(a, b []byte) []byte {
	var left float64
	var right float64

	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)

	var res bytes.Buffer
	binary.Write(&res, binary.LittleEndian, left-right)

	return res.Bytes()
}

func FMult(a, b []byte) []byte {
	var left float64
	var right float64
	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)

	var res bytes.Buffer
	binary.Write(&res, binary.LittleEndian, left*right)

	return res.Bytes()
}

func FDiv(a, b []byte) []byte {
	var left float64
	var right float64
	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)

	var res bytes.Buffer
	binary.Write(&res, binary.LittleEndian, left/right)

	return res.Bytes()
}
