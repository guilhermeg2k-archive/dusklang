package vm

const (
	ILOAD_CONST byte = 0
	ILOAD       byte = 1
	ISTORE      byte = 2
	IADD        byte = 3
	ISUB        byte = 4
	IMULT       byte = 5
	IDIV        byte = 6
	IMOD        byte = 7

	FLOAD_CONST byte = 8
	FLOAD       byte = 9
	FSTORE      byte = 10
	FADD        byte = 11
	FSUB        byte = 12
	FMULT       byte = 13
	FDIV        byte = 14

	BOLOAD_CONST byte = 15
	BOLOAD       byte = 16
	BOSTORE      byte = 18

	ICMP_EQUALS         byte = 19
	ICMP_LESS_EQUALS    byte = 20
	ICMP_GREATER_EQUALS byte = 21
	ICMP_LESS_THEN      byte = 22
	ICMP_GREATER_THEN   byte = 23

	FCMP_EQUALS         byte = 24
	FCMP_LESS_EQUALS    byte = 25
	FCMP_GREATER_EQUALS byte = 26
	FCMP_LESS_THEN      byte = 27
	FCMP_GREATER_THEN   byte = 28

	FUNCCALL byte = 29
	PRINT         = 99
)
