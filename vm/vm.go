package vm

type Frame map[uint64][]byte
type Consts map[uint64][]byte
type Stack []byte
type VirtualMachine struct {
	Stack     *Stack
	Functions *[]Function
}

type Function struct {
	Consts   Consts
	Frame    *Frame
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

func store(frame *Frame, offset uint64, bytes []byte) {
	_frame := *frame
	_frame[offset] = bytes
	frame = &_frame
}

func load(frame *Frame, offset uint64) []byte {
	_frame := *frame
	return _frame[offset]
}
