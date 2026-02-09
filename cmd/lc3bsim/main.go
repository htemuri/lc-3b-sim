package main

import (
	"lc3b-sim/m/v2/cpu"
	"log/slog"
	"os"
)

func main() {

	var logger slog.Logger = *slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	var instructions []uint16

	// 1. LEA R0, #6 (Load Effective Address)
	// Opcode 1110 (LEA) | DR 000 (R0) | PCoffset9 000000110 (6)
	// Target = PC_next (0x3002) + (6 << 1) = 0x300E
	instructions = append(instructions, 0xE006)

	// 2. AND R1, R1, #0 (Clear R1)
	instructions = append(instructions, 0x5260)

	// 3. ADD R1, R1, #7 (Set R1 = 7)
	instructions = append(instructions, 0x1267)

	// 4. STW R1, R0, #0 (Store Word)
	// Opcode 0111 (STW) | SR 001 (R1) | BaseR 000 (R0) | offset6 000000
	// Mem[R0 + (0<<1)] = R1
	instructions = append(instructions, 0x7200)

	// 5. LDW R2, R0, #0 (Load Word)
	// Opcode 0110 (LDW) | DR 010 (R2) | BaseR 000 (R0) | offset6 000000
	// R2 = Mem[R0 + (0<<1)]
	instructions = append(instructions, 0x6400)

	// 6. HALT
	instructions = append(instructions, 0xF025)

	// 7. Padding (0x300C) - this instruction is skipped/unused
	instructions = append(instructions, 0x0000)

	// 8. STORAGE SLOT (0x300E) - This is where R0 points
	// We initialize it to 0xFFFF so we know for sure STW overwrote it.
	instructions = append(instructions, 0xFFFF)

	var cpu cpu.CPU
	cpu.Init(0x3000, instructions, logger)
	cpu.Run()

	// var u int8 = int8(datapath.SEXT(0b11111) + uint16(5)) // in binary: 11111111
	// // casting to int8 re-interprets 11111111 as -1 in two's complement
	// // var s int8 = int8(u)
	// log.Printf("%016b", 0x6400) // Output: -1

}
