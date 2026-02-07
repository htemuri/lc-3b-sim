package cpu

import (
	"lc3b-sim/m/v2/control"
	"lc3b-sim/m/v2/datapath"
	"log"
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
}

func (cpu *CPU) Init(pcStart uint16, programData [65536]uint8) {
	cpu.state = 18
	cpu.pc.UpdateValue(control.LoadSig_LOAD, pcStart)
	cpu.pc.Commit()
	cpu.memory.Init(programData)
}

func (cpu *CPU) Tick() {
	log.Println("STATE: ", cpu.state)
	log.Printf("PC: 0x%x", cpu.pc.GetValue())
	controlStoreOut := control.ControlStore(cpu.state)
	signals := controlStoreOut.DatapathSignals

	cpu.computeCombinationalLogic(signals)
	cpu.bus = cpu.computeBus(signals)
	log.Printf("BUS value: 0x%x", cpu.bus)
	log.Printf("IR value: 0x%x", cpu.ir.GetValue())
	cpu.updateRegisters(signals)
	cpu.updateMemory(signals)

	// clock edge commits
	cpu.pc.Commit()
	cpu.ir.Commit()
	cpu.registerFile.Commit()
	cpu.mar.Commit()
	cpu.mdr.Commit()
	cpu.conditionalCodes.Commit()
	cpu.memory.Commit()

	// calculate new state

	ir15to11 := uint8((cpu.ir.GetValue() >> 11)) & 0b11111
	log.Println("Memory state: ", cpu.memory.Ready)
	cpu.memory.PrintMem(0x2)
	cpu.state = control.Microsequencer(controlStoreOut.COND, cpu.calculateBen(), cpu.memory.Ready, ir15to11, controlStoreOut.J, controlStoreOut.IRD)
}

func (cpu *CPU) PrintRegisterFile() {
	cpu.registerFile.PrintValues()
}

func (cpu *CPU) computeCombinationalLogic(signals control.Signals) {
	ir := cpu.ir.GetValue()
	pc := cpu.pc.GetValue()

	// register file reads
	irInput11_9 := uint8((ir >> 9)) & 0b111
	irInput8_6 := uint8((ir >> 6)) & 0b111
	cpu.registerFile.Read(datapath.SR1Mux(signals.Sr1MUX, irInput11_9, irInput8_6), datapath.GetSR2Input(ir))

	// adder
	addr2muxOutput := datapath.Addr2Mux(signals.Addr2MUX, ir&0b1111111111, ir&0b111111111, ir&0b111111)
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
			bus = datapath.SEXT(bus) // sign extend if byte
		} else {
			bus = mdr
		}
	}

	if signals.GateALU == control.NoYesSig_YES {
		// log.Println("ALU Gate opened")
		enabledCounter += 1
		sr1 := cpu.registerFile.Sr1Out
		sr2 := datapath.SR2Mux(control.SR2Mux(0b100000&ir>>5 == 1), 0b11111&ir, cpu.registerFile.Sr2Out)
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
	cpu.mdr.UpdateValue(signals.LdMDR, datapath.MioENMux(signals.MioEN, modifiedBus, cpu.memory.DataOut))
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
