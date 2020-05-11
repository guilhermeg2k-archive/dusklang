package glang

import (
	"encoding/binary"
	"fmt"
)

func Print(bytes []byte) {
	value, _ := binary.Varint(bytes)
	fmt.Println(value)
}
