package main

import "fmt"

func main() {
	fmt.Println("hello, world!")
}

// using 16 bit unsigned ints because each word in lc-3b is 16 bits and we're assuming 2s complement form

// Adder
func Adder(input1, input2 uint16) uint16 {
	return input1 + input2
}

// ALU
type ALUOp uint8

const (
	ALU_ADD ALUOp = iota
	ALU_AND
	ALU_XOR
	ALU_PASS
)

func ALU(
	op ALUOp, // two bit control signal for function
	a, b uint16,
) uint16 {
	switch op {
	case ALU_ADD: // 00
		return a + b
	case ALU_AND: // 01
		return a & b
	case ALU_XOR: // 10
		return a ^ b
	case ALU_PASS: // 11
		return a
	default:
		panic("invalid ALU operation")
	}
}

// General purpose registers - 8x16 bit

type GPRegister uint8

const (
	R0 GPRegister = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
)

type RegisterFile struct {
	regs [8]uint16

	sr1Out uint16
	sr2Out uint16

	// adding these because writes are sequential and depend on the clock
	pendingWrite bool
	writeData    uint16
	writeReg     GPRegister
}

func (rf *RegisterFile) Read(
	sr1, sr2 GPRegister,
) {
	rf.sr1Out = rf.regs[sr1]
	rf.sr2Out = rf.regs[sr2]
}

func (rf *RegisterFile) Write(
	ldREG bool,
	dr GPRegister,
	data uint16,
) {
	if ldREG {
		rf.pendingWrite = true
		rf.writeData = data
		rf.writeReg = dr
	}
}

func (rf *RegisterFile) Commit() {
	if rf.pendingWrite {
		rf.regs[rf.writeReg] = rf.writeData
		rf.pendingWrite = false
	}
}

type DataSize uint8

const (
	WORD DataSize = iota
	BYTE
)

type Memory struct {
	mem [65536]uint8 // 2**16 addresses with 8 bit addressability. words are 16bits and are aligned if their addresses differ only in bit [0]

	pendingWrite bool
	dataSize     DataSize
	writeEnable1 bool
	writeEnable2 bool
	writeDate    uint16
}

func (m *Memory) Read(
	address uint16,
) uint8 {
	return 0
}
