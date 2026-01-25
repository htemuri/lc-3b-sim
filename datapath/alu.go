package datapath

import "lc3b-sim/m/v2/control"

func ALU(
	op control.ALUOp, // two bit control signal for function
	a, b uint16,
) uint16 {
	switch op {
	case control.ALU_ADD:
		return a + b
	case control.ALU_AND:
		return a & b
	case control.ALU_XOR:
		return a ^ b
	case control.ALU_PASSA:
		return a
	default:
		panic("invalid ALU operation")
	}
}
