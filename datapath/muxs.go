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

func SR2Mux(controlSig control.SR2Mux, irInput, sr2outInput uint16) uint16 {
	switch controlSig {
	case control.IR:
		return irInput
	case control.SR2OUT:
		return sr2outInput
	default:
		panic("invalid control signal")
	}
}

func MARMux(marMux control.MARMUX, seven0Input uint8, adderInput uint16) uint16 {
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

func MioENMux(mioEN control.NoYesSig, memoryInput, mioInput uint16) uint16 {
	if mioEN == control.NoYesSig_NO {
		return memoryInput
	} else {
		return mioInput
	}
}
