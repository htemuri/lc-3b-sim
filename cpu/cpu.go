package cpu

import (
	"lc3b-sim/m/v2/control"
	"lc3b-sim/m/v2/datapath"
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

func (cpu *CPU) Tick() {
	controlStoreOut := control.ControlStore(cpu.state)
	signals := controlStoreOut.DatapathSignals

	cpu.bus = cpu.computeBus(signals)

}

func (cpu *CPU) computeCombinationalLogic(signals control.Signals) {
	ir := cpu.ir.GetValue()
	pc := cpu.pc.GetValue()

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
		enabledCounter += 1
		sr1 := cpu.registerFile.Sr1Out
		sr2 := datapath.SR2Mux(control.SR2Mux(0b100000&ir == 1), 0b11111&ir, cpu.registerFile.Sr2Out)
		bus = datapath.ALU(signals.AluK, sr1, sr2)
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
}
