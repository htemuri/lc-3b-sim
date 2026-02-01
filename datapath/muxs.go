package datapath

import (
	"lc3b-sim/m/v2/control"
)

func PCMux(PCMux2 control.PCMux2, busInput, adderInput, pcPlus2Input uint16) uint16 {
	switch PCMux2 {
	case control.PCPLUS2:
		return pcPlus2Input
	case control.ADDER:
		return adderInput
	case control.BUS:
		return busInput
	default:
		panic("invalid control signal")
	}
}

func SR1Mux(controlSig control.SR1MUX, irInput11_9, irInput8_6 uint8) gpRegister {
	if controlSig == control.SR1MUX_11_9 {
		return gpRegister(irInput11_9 & 0b111)
	} else {
		return gpRegister(irInput8_6 & 0b111)
	}
}

func SR2Mux(controlSig control.SR2Mux, irInput, sr2outInput uint16) uint16 {
	switch controlSig {
	case control.SR2Mux_IR:
		return irInput
	case control.SR2Mux_SR2OUT:
		return sr2outInput
	default:
		panic("invalid control signal")
	}
}

func MARMux(marMux control.MARMUX, seven0Input uint16, adderInput uint16) uint16 {
	if marMux == control.MARMUX_7_0 {
		return uint16(seven0Input)
	} else {
		return adderInput
	}
}

func Addr1Mux(addr1Mux control.ADDR1MUX, pcInput uint16, sr1OutInput uint16) uint16 {
	if addr1Mux == control.ADDR1MUX_PC {
		return pcInput
	} else {
		return sr1OutInput
	}
}

func Addr2Mux(addr2Mux control.ADDR2MUX, top11Input, top9Input, top6Input uint16) uint16 {
	switch addr2Mux {
	case control.ADDR2MUX_ZERO:
		return 0
	case control.ADDR2MUX_PCoffset11:
		return top11Input
	case control.ADDR2MUX_PCoffset9:
		return top9Input
	case control.ADDR2MUX_offset6:
		return top6Input
	default:
		panic("invalid control signal")
	}
}

func DrMux(drMuxSig control.DRMUX, irInput11_9 uint8) gpRegister {
	if drMuxSig == control.DRMUX_11_9 {
		return gpRegister(irInput11_9)
	} else {
		return gpRegister(0b111)
	}
}

func MioENMux(mioEN control.NoYesSig, busInput, memoryInput uint16) uint16 {
	if mioEN == control.NoYesSig_NO {
		return busInput
	} else {
		return memoryInput
	}
}
