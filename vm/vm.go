package vm

type Storage map[uint64][]byte
type Consts map[uint64][]byte
type Stack []byte

type VirtualMachine struct {
	Stack     *Stack
	Functions *[]Function
}

type Function struct {
	Consts   Consts
	Storage  *Storage
	Bytecode []byte
}

func push(stack *Stack, bytes []byte) {
	*stack = append(*stack, bytes...)
}

func pop(stack *Stack, bytes int) []byte {
	poped := (*stack)[len(*stack)-bytes:]
	*stack = append((*stack)[:len(*stack)-bytes])
	return poped
}

func store(storage *Storage, offset uint64, bytes []byte) {
	(*storage)[offset] = bytes
}

func load(storage *Storage, offset uint64) []byte {
	return (*storage)[offset]
}
