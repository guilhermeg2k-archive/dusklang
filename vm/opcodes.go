package vm

const (
	ILOAD_CONST byte = 0
	ILOAD       byte = 1
	ISTORE      byte = 2
	IPOP        byte = 3
	IADD        byte = 4
	ISUB        byte = 5
	IMULT       byte = 6
	IDIV        byte = 7
	IMOD        byte = 8

	FLOAD_CONST byte = 9
	FLOAD       byte = 10
	FSTORE      byte = 11
	FPOP        byte = 12
	FADD        byte = 13
	FSUB        byte = 14
	FMULT       byte = 15
	FDIV        byte = 16
	FMOD        byte = 17

	BLOAD_CONST byte = 18
	BLOAD       byte = 19
	BSTORE      byte = 20
	BPOP        byte = 21
	BADD        byte = 22
	BSUB        byte = 23
	BMULT       byte = 24
	BDIV        byte = 25
	BMOD        byte = 26

	FUNCCALL byte = 27
	PRINT         = 28
)
