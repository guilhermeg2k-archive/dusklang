package dusk

import (
	"bytes"
	"encoding/binary"
)

func ICmpEquals(a, b []byte) []byte {
	var left int64
	var right int64
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	if left == right {
		return []byte{1}
	}
	return []byte{0}
}

func ICmpLessEquals(a, b []byte) []byte {
	var left int64
	var right int64
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	if left <= right {
		return []byte{1}
	}
	return []byte{0}
}

func ICmpGreaterEquals(a, b []byte) []byte {
	var left int64
	var right int64
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	if left >= right {
		return []byte{1}
	}
	return []byte{0}
}

func ICmpLessThen(a, b []byte) []byte {
	var left int64
	var right int64
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	if left < right {
		return []byte{1}
	}
	return []byte{0}
}

func ICmpGreaterThen(a, b []byte) []byte {
	var left int64
	var right int64
	left, _ = binary.Varint(a)
	right, _ = binary.Varint(b)
	if left > right {
		return []byte{1}
	}
	return []byte{0}
}

func FCmpEquals(a, b []byte) []byte {
	var left float64
	var right float64
	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)
	if left == right {
		return []byte{1}
	}
	return []byte{0}
}

func FCmpLessEquals(a, b []byte) []byte {
	var left float64
	var right float64
	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)
	if left <= right {
		return []byte{1}
	}
	return []byte{0}
}

func FCmpGreaterEquals(a, b []byte) []byte {
	var left float64
	var right float64
	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)
	if left >= right {
		return []byte{1}
	}
	return []byte{0}

}

func FCmpLessThen(a, b []byte) []byte {
	var left float64
	var right float64
	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)
	if left < right {
		return []byte{1}
	}
	return []byte{0}

}

func FCmpGreaterThen(a, b []byte) []byte {
	var left float64
	var right float64
	buf := bytes.NewReader(a)
	binary.Read(buf, binary.LittleEndian, &left)
	buf = bytes.NewReader(b)
	binary.Read(buf, binary.LittleEndian, &right)
	if left > right {
		return []byte{1}
	}
	return []byte{0}
}
