package main

import (
	"fmt"
	"lc3b-sim/m/v2/cpu"
	"log"
)

// using 16 bit unsigned ints because each word in lc-3b is 16 bits and we're assuming 2s complement form

func main() {
	// x := uint16(0b1010100101000111)
	// fmt.Printf("%016b", uint8(x))

	var instructions []uint16
	var programMemory [65536]uint8

	// instructions = append(instructions, 0b0001000000100001) //  ADD  DR=0  SR1=0 imm5=1
	log.Printf("0x3202 = %016b", 0x3202)
	log.Println(0 & 5)

	instructions = append(instructions, 0x5260)
	instructions = append(instructions, 0x14A5)
	instructions = append(instructions, 0x1242)
	instructions = append(instructions, 0x3202)
	instructions = append(instructions, 0xF025)
	// instructions = append(instructions, 0x0000)
	// instructions = append(instructions, 0x3002)
	// instructions = append(instructions, 0xF025)
	// instructions = append(instructions, 0x0000)
	// instructions = append(instructions, 0x0000)

	for i := 0; i < len(instructions)*2; i += 2 {
		programMemory[0x3000+i] = uint8(instructions[i/2])
		programMemory[0x3000+i+1] = uint8(instructions[i/2] >> 8)
	}
	log.Println(programMemory[0x3000:0x3010])

	// for index, instruction := range instructions {
	// 	log.Println(index)
	// 	log.Printf("%08b", uint8(instruction))
	// 	log.Printf("%08b", uint8(instruction>>8))
	// 	programMemory[0x3000+index] = uint8(instruction)
	// 	programMemory[0x3000+index+1] = uint8(instruction >> 8)
	// }
	// log.Printf("0x%x%x", programMemory[0x3009], programMemory[0x3008])
	//  = uint8(instructions[0])

	var cpu cpu.CPU
	cpu.Init(0x3000, programMemory)
	for {
		cpu.Tick()
		cpu.PrintRegisterFile()
		log.Println("---------------------")
		fmt.Scanln()
	}
}
