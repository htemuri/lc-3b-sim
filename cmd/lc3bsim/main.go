package main

import (
	"lc3b-sim/m/v2/cpu"
	"log/slog"
	"os"
)

func main() {
	var cpu cpu.CPU

	pcStart := uint16(0x3000)
	instructions := []uint16{
		0xE006, // LEA R0, #6 (Load Effective Address)
		0x5260, // AND R1, R1, #0 (Clear R1)
		0x1267, // ADD R1, R1, #7 (Set R1 = 7)
		0x7200, // STW R1, R0, #0 (Store Word)
		0x6400, // LDW R2, R0, #0 (Load Word)
		0xF025, // HALT
		0x0000, // Padding (0x300C)
		0xFFFF, // STORAGE SLOT (0x300E)
	}

	cpu.Init(pcStart, instructions, *slog.New(slog.NewTextHandler(os.Stdout, nil)))
	cpu.Run()
}
