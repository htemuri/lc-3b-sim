package datapath

import "lc3b-sim/m/v2/control"

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
