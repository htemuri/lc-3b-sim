package cpu

import (
	"fmt"
	"lc3b-sim/m/v2/control"
	"lc3b-sim/m/v2/datapath"
	"log"
	"log/slog"
)

type CPU struct {

	// state
	state uint8

	bus uint16

	// registers
	pc               datapath.Register16bit
	ir               datapath.Register16bit
	registerFile     datapath.RegisterFile
	mar              datapath.Register16bit
	mdr              datapath.Register16bit
	conditionalCodes datapath.ConditionalCodes

	// memory
	memory datapath.Memory

	// temp vars so dont have to repeat ops
	adderOutput uint16

	halted bool // if we've hit the trap instruction

	logger slog.Logger
}

func (cpu *CPU) Init(pcStart uint16, instructions []uint16, logger slog.Logger) {
	cpu.logger = logger
	cpu.state = 18
	cpu.pc.UpdateValue(control.LoadSig_LOAD, pcStart)
	cpu.pc.Commit()
	var programMemory [65536]uint8
	for i := 0; i < len(instructions)*2; i += 2 {
		programMemory[pcStart+uint16(i)] = uint8(instructions[i/2])
		programMemory[pcStart+uint16(i)+1] = uint8(instructions[i/2] >> 8)
	}
	cpu.memory.Init(programMemory)
}

func (cpu *CPU) Run() {
	// i := 0
	for !cpu.halted && cpu.pc.GetValue() < 0x301f {
		cpu.Tick()
		// cpu.PrintRegisterFile()
		// log.Println("---------------------")
		// fmt.Scanln()
	}
}

func (cpu *CPU) Tick() {
	cpu.memory.Commit()

	if cpu.state == 18 || cpu.state == 19 {
		// cpu.logger.Info(fmt.Sprintf("PC: %x", cpu.pc.GetValue()-2))
		// // cpu.logger.Info(fmt.Sprintf("IR: %16b", cpu.ir.GetValue()))
		// cpu.logger.Info(fmt.Sprintf("BEN: %v", cpu.calculateBen()))

		cpu.logger.Info(fmt.Sprintf("R1: %d; R2: %d", int8(cpu.registerFile.GetRegisters()[1]), int8(cpu.registerFile.GetRegisters()[2])))
		// cpu.logger.Info(fmt.Sprintf("Condition codes: N: %v, Z: %v, P: %v", cpu.conditionalCodes.N, cpu.conditionalCodes.Z, cpu.conditionalCodes.P))
	}

	// if cpu.state == 32 {
	// 	cpu.logger.Info(fmt.Sprintf("IR: %16b", cpu.ir.GetValue()))
	// }

	// if cpu.state == 33 {
	// 	cpu.logger.Info(fmt.Sprintf("MDR: %v", cpu.mdr))
	// 	cpu.logger.Info(fmt.Sprintf("MEM: %v", cpu.memory.DataOut))

	// }

	if cpu.state == 30 {
		cpu.logger.Info("Halting CPU due to TRAP instruction")
		cpu.halted = true
		return
	}
	controlStoreOut := control.ControlStore(cpu.state)
	signals := controlStoreOut.DatapathSignals

	cpu.computeCombinationalLogic(signals)
	cpu.bus = cpu.computeBus(signals)
	cpu.updateMemory(signals)
	cpu.updateRegisters(signals)

	// clock edge commits
	cpu.pc.Commit()
	cpu.ir.Commit()
	cpu.registerFile.Commit()
	cpu.mar.Commit()
	cpu.mdr.Commit()
	cpu.conditionalCodes.Commit()

	// calculate new state

	ir15to11 := uint8((cpu.ir.GetValue() >> 11)) & 0b11111
	cpu.state = control.Microsequencer(controlStoreOut.COND, cpu.calculateBen(), cpu.memory.Ready, ir15to11, controlStoreOut.J, controlStoreOut.IRD)
	// cpu.logger.Info(fmt.Sprintf("New State: %d", cpu.state))

}

// func (cpu *CPU) PrintRegisterFile() {
// 	cpu.registerFile.PrintValues()
// }

func (cpu *CPU) computeCombinationalLogic(signals control.Signals) {
	ir := cpu.ir.GetValue()
	pc := cpu.pc.GetValue()

	// register file reads
	irInput11_9 := uint8((ir >> 9)) & 0b111
	irInput8_6 := uint8((ir >> 6)) & 0b111
	cpu.registerFile.Read(datapath.SR1Mux(signals.Sr1MUX, irInput11_9, irInput8_6), datapath.GetSR2Input(ir))

	// adder
	addr2muxOutput := datapath.Addr2Mux(signals.Addr2MUX, datapath.SEXT(ir&0b1111111111, 11), datapath.SEXT(ir&0b111111111, 9), datapath.SEXT(ir&0b111111, 6))
	addr1muxOutput := datapath.Addr1Mux(signals.Addr1MUX, pc, cpu.registerFile.Sr1Out)
	cpu.adderOutput = datapath.Adder(datapath.LSHF1(signals.Lshf1, addr2muxOutput), addr1muxOutput)

}

func (cpu *CPU) computeBus(signals control.Signals) uint16 {
	enabledCounter := 0
	var bus uint16

	ir := cpu.ir.GetValue()

	if signals.GatePC == control.NoYesSig_YES {
		enabledCounter += 1
		bus = cpu.pc.GetValue()
	}

	if signals.GateMDR == control.NoYesSig_YES {
		enabledCounter += 1
		mdr := cpu.mdr.GetValue()
		mar0 := cpu.mar.GetValue() & 0b1
		if signals.DataSize == control.BYTE {
			if mar0 == 0 {
				bus = mdr & 0b11111111 // bus = mdr[7:0]
			} else {
				bus = (mdr & 0b1111111100000000) >> 8 // bus = mdr[15:8]
			}
			bus = datapath.SEXT(bus, 8) // sign extend if byte
		} else {
			bus = mdr
		}
	}

	if signals.GateALU == control.NoYesSig_YES {
		// log.Println("ALU Gate opened")
		enabledCounter += 1
		sr1 := cpu.registerFile.Sr1Out
		sr2 := datapath.SR2Mux(control.SR2Mux(0b100000&ir>>5 == 1), datapath.SEXT(0b11111&ir, 5), cpu.registerFile.Sr2Out)
		// log.Println("SR2: ", sr2)
		// log.Println("ALUK: ", signals.AluK)
		bus = datapath.ALU(signals.AluK, sr1, sr2)
		log.Println("Calculated BUS: ", bus)
	}

	if signals.GateMARMUX == control.NoYesSig_YES {
		enabledCounter += 1
		bus = datapath.MARMux(signals.MarMUX, datapath.ZEXTandLSHF1(uint8(ir)), cpu.adderOutput)
	}

	if signals.GateSHF == control.NoYesSig_YES {
		enabledCounter += 1
		ir := cpu.ir.GetValue()
		bus = datapath.SHF(cpu.registerFile.Sr1Out, uint8(ir&0b111111))
	}

	if enabledCounter > 1 {
		panic("more than 1 gate open at same time")
	}
	return bus
}

func (cpu *CPU) updateRegisters(signals control.Signals) {
	cpu.pc.UpdateValue(signals.LdPC, datapath.PCMux(signals.PcMUX, cpu.bus, cpu.adderOutput, datapath.PlusTwo(cpu.pc.GetValue())))
	cpu.ir.UpdateValue(signals.LdIR, cpu.bus)
	cpu.mar.UpdateValue(signals.LdMAR, cpu.bus)
	cpu.conditionalCodes.SetCC(signals.LdCC, cpu.bus)

	// assuming no keyboard and sdram for now so the mioen mux will always either have bus input or memory input
	mar0 := cpu.mar.GetValue() & 0b1
	var modifiedBus uint16
	if signals.DataSize == control.BYTE {
		if mar0 == 0 {
			modifiedBus = uint16(uint8(cpu.bus)) // modifiedBus = bus[7:0]
		} else {
			modifiedBus = uint16(uint8(cpu.bus)) << 8 // modifiedBus[15:8] = bus[7:0]
		}
	} else {
		modifiedBus = cpu.bus
	}
	var memData uint16
	if cpu.memory.Ready {
		memData = cpu.memory.DataOut
		cpu.memory.PendingRead = false
	}
	cpu.mdr.UpdateValue(signals.LdMDR, datapath.MioENMux(signals.MioEN, modifiedBus, memData))
	irInput11_9 := uint8((cpu.ir.GetValue() >> 9)) & 0b111
	cpu.registerFile.Write(signals.LdREG, datapath.DrMux(signals.DrMUX, irInput11_9), cpu.bus)
}

func (cpu *CPU) updateMemory(signals control.Signals) {
	if signals.MioEN == control.NoYesSig_YES {
		mar := cpu.mar.GetValue()
		mar0 := mar & 0b1
		if signals.Rw == control.RW_WR {
			var we0, we1 bool
			if signals.DataSize == control.BYTE {
				if mar0 == 0 {
					we0 = true
					we1 = false
				} else {
					we0 = false
					we1 = true
				}
			} else {
				we0 = true
				we1 = true
			}
			cpu.memory.Write(mar, cpu.mdr.GetValue(), we0, we1)
		} else {
			cpu.memory.Read(mar)
		}
	}
}

func (cpu *CPU) calculateBen() bool {
	ir := cpu.ir.GetValue()
	ir11 := (ir & (1 << 11)) == 1<<11
	ir10 := (ir & (1 << 10)) == 1<<10
	ir9 := (ir & (1 << 9)) == 1<<9
	ben := ir11 && cpu.conditionalCodes.N || ir10 && cpu.conditionalCodes.Z || ir9 && cpu.conditionalCodes.P
	return ben
}
