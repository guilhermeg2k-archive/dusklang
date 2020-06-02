package vm

import (
	"fmt"
	"io"
)

type Storage map[uint64][]byte
type Consts map[uint64][]byte
type Labels map[uint64]int
type Stack []byte

type VirtualMachine struct {
	Stack     *Stack
	Functions *[]Function
}

type Function struct {
	Labels        Labels
	Consts        Consts
	Storage       Storage
	Bytecode      []byte
	CurrentOffset int
}

func (f *Function) readByte() (byte, error) {
	bytes, err := f.readBytes(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil
}

func (f *Function) readBytes(n int) ([]byte, error) {
	if f.CurrentOffset == len(f.Bytecode) {
		return []byte{}, io.EOF
	}
	defer func() {
		f.CurrentOffset += n
	}()
	return f.Bytecode[f.CurrentOffset : f.CurrentOffset+n], nil
}

func push(stack *Stack, bytes []byte) {
	*stack = append(*stack, bytes...)
}

func pop(stack *Stack, bytes int) []byte {
	poped := (*stack)[len(*stack)-bytes:]
	*stack = append((*stack)[:len(*stack)-bytes])
	return poped
}

func store(storage Storage, offset uint64, bytes []byte) {
	var b []byte
	b = append(b, bytes...)
	storage[offset] = b
}

func load(storage Storage, offset uint64) []byte {
	fmt.Println(offset, storage[offset], storage[offset])
	return storage[offset]
}
