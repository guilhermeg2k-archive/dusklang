package dusk

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//Temporary
func Print(_bytes []byte) {
	var value int64
	buffer := bytes.NewBuffer(_bytes)
	binary.Read(buffer, binary.LittleEndian, &value)
	fmt.Println(value)
}
