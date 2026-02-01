package main

import (
	"fmt"
	"lc3b-sim/m/v2/cpu"
)

// using 16 bit unsigned ints because each word in lc-3b is 16 bits and we're assuming 2s complement form

func main() {
	// x := uint16(0b1010100101000111)
	// fmt.Printf("%016b", uint8(x))

	var instructions []uint16
	var programMemory [65536]uint8

	instructions = append(instructions, 0b0001000000100001) //  ADD  DR=0  SR1=0 imm5=1
	programMemory[0x3000] = uint8(instructions[0])
	programMemory[0x3001] = uint8(instructions[0] >> 8)

	var cpu cpu.CPU
	cpu.Init(0x3000, programMemory)
	for {
		cpu.Tick()
		cpu.PrintRegisterFile()
		fmt.Scanln()
	}
}
