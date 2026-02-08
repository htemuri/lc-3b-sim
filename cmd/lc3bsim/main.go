package main

import (
	"lc3b-sim/m/v2/cpu"
	"log/slog"
	"os"
)

func main() {

	var logger slog.Logger = *slog.New(slog.NewTextHandler(os.Stdout, nil))
	var instructions []uint16

	instructions = append(instructions, 0x54A0)
	instructions = append(instructions, 0x5260)
	instructions = append(instructions, 0x1265)
	instructions = append(instructions, 0x1481)
	instructions = append(instructions, 0x127F)
	instructions = append(instructions, 0x03FD)
	instructions = append(instructions, 0xF025)
	// instructions = append(instructions, 0x0000)
	// instructions = append(instructions, 0x0000)
	var cpu cpu.CPU
	cpu.Init(0x3000, instructions, logger)
	cpu.Run()

	// var u int8 = int8(datapath.SEXT(0b11111) + uint16(5)) // In binary: 11111111
	// // Casting to int8 re-interprets 11111111 as -1 in two's complement
	// // var s int8 = int8(u)
	// log.Printf("%d", int8(datapath.SEXT(0b111111101, 9))) // Output: -1

}
