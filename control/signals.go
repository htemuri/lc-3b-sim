package control

type DataSize uint8

const (
	WORD DataSize = iota
	BYTE
)

type ALUOp uint8

const (
	ALU_ADD ALUOp = iota
	ALU_AND
	ALU_XOR
	ALU_PASS
)

type LD_REG bool
