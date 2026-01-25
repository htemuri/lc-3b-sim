package control

type PCMux2 uint8

const (
	PCPLUS2 PCMux2 = iota
	BUS
	ADDER
)

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

type SR2Mux uint8

const (
	IR SR2Mux = iota
	SR2OUT
)

type LD_REG bool
type LD_PC bool
type LD_MAR bool
