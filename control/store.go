package control

// control store should be 35 bits x 2**6 (possible state combinations)

type Condition uint8

const (
	COND_UNCONDITIONAL Condition = iota
	COND_MEMREADY
	COND_BRANCH
	COND_ADDRMODE
)

type ControlStoreOutput struct {
	IRD             bool
	COND            Condition
	J               uint8
	DatapathSignals Signals
}

var store map[int]ControlStoreOutput = map[int]ControlStoreOutput{
	// state: [ird, cond1, cond0, j5...j0, ldmar, ldmdr, ldir, ldben, ldreg, ldcc, ldpc, gatepc, gatemdr, gatealu,
	// gatemarmux, gateshf, pcmux(2), drmux, sr1mux, addr1mux, addr2mux(2), marmux, aluk(2), mioen, rw, datasize, lshf1]
	0:  {COND: COND_BRANCH, J: 18},
	1:  {J: 18, DatapathSignals: Signals{aluK: ALU_ADD, sr1MUX: SR1MUX_8_6, ldCC: LoadSig_LOAD, ldREG: LoadSig_LOAD}},
	2:  {J: 29, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD}},
	3:  {J: 24, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD}},
	4:  {COND: COND_ADDRMODE, J: 20},
	5:  {J: 18, DatapathSignals: Signals{aluK: ALU_AND, sr1MUX: SR1MUX_8_6, ldCC: LoadSig_LOAD, ldREG: LoadSig_LOAD}},
	6:  {J: 25, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD, lshf1: NoYesSig_YES}},
	7:  {J: 23, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD, lshf1: NoYesSig_YES}},
	8:  {},
	9:  {J: 18, DatapathSignals: Signals{aluK: ALU_XOR, sr1MUX: SR1MUX_8_6, ldCC: LoadSig_LOAD, ldREG: LoadSig_LOAD}},
	10: {},
	11: {},
	12: {J: 18, DatapathSignals: Signals{pcMUX: ADDER, addr1MUX: ADDR1MUX_BaseR, ldPC: LoadSig_LOAD, sr1MUX: SR1MUX_8_6}},
	13: {J: 18, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, ldREG: LoadSig_LOAD}},
	14: {J: 18, DatapathSignals: Signals{addr1MUX: ADDR1MUX_PC, addr2MUX: ADDR2MUX_PCoffset9, lshf1: NoYesSig_YES, marMUX: MARMUX_ADDER, ldREG: LoadSig_LOAD, ldCC: LoadSig_LOAD}},
	15: {J: 28, DatapathSignals: Signals{marMUX: MARMUX_7_0, ldMAR: LoadSig_LOAD}},
}

func ControlStore(state uint8) ControlStoreOutput {

	return ControlStoreOutput{}
}
