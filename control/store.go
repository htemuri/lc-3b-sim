package control

// The control store is basically taking the finite state machine which you can see in the README
// and converting it to the signals needed in the datapath to accomplish the state and move on.
// So for ex, if i'm on state 1...
//			- I need to do a microoperation which is adding sr1 + op2 and storing it in the dest register.
// 			  To do that, i need to set my alu control signal to add, select which bits from the instruction are
// 			  my sr1, and set the load signals for the registers im writing to. I also enable the tristate buffer
// 			  that allows my alu output to flow through the bus to my dest register and to my conditional code registers (n,z,p).
// 			- I set the `J` value to 18 to declare that once I finish my microoperation, I should move to state 18

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

// control store should be 35 bits x 2**6 (possible state combinations)
var store map[int]ControlStoreOutput = map[int]ControlStoreOutput{
	0:  {COND: COND_BRANCH, J: 18},
	1:  {J: 18, DatapathSignals: Signals{aluK: ALU_ADD, sr1MUX: SR1MUX_8_6, ldCC: LoadSig_LOAD, ldREG: LoadSig_LOAD, gateALU: NoYesSig_YES}},
	2:  {J: 29, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD, gateMARMUX: NoYesSig_YES}},
	3:  {J: 24, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD, gateMARMUX: NoYesSig_YES}},
	4:  {COND: COND_ADDRMODE, J: 20},
	5:  {J: 18, DatapathSignals: Signals{aluK: ALU_AND, sr1MUX: SR1MUX_8_6, ldCC: LoadSig_LOAD, ldREG: LoadSig_LOAD, gateALU: NoYesSig_YES}},
	6:  {J: 25, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD, lshf1: NoYesSig_YES, gateMARMUX: NoYesSig_YES}},
	7:  {J: 23, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, addr2MUX: ADDR2MUX_offset6, marMUX: MARMUX_ADDER, addr1MUX: ADDR1MUX_BaseR, ldMAR: LoadSig_LOAD, lshf1: NoYesSig_YES, gateMARMUX: NoYesSig_YES}},
	8:  {},
	9:  {J: 18, DatapathSignals: Signals{aluK: ALU_XOR, sr1MUX: SR1MUX_8_6, ldCC: LoadSig_LOAD, ldREG: LoadSig_LOAD, gateALU: NoYesSig_YES}},
	10: {},
	11: {},
	12: {J: 18, DatapathSignals: Signals{pcMUX: ADDER, addr1MUX: ADDR1MUX_BaseR, ldPC: LoadSig_LOAD, sr1MUX: SR1MUX_8_6}},
	13: {J: 18, DatapathSignals: Signals{sr1MUX: SR1MUX_8_6, ldREG: LoadSig_LOAD, gateSHF: NoYesSig_YES, ldCC: LoadSig_LOAD}},
	14: {J: 18, DatapathSignals: Signals{addr1MUX: ADDR1MUX_PC, addr2MUX: ADDR2MUX_PCoffset9, lshf1: NoYesSig_YES, marMUX: MARMUX_ADDER, ldREG: LoadSig_LOAD, ldCC: LoadSig_LOAD, gateMARMUX: NoYesSig_YES}},
	15: {J: 28, DatapathSignals: Signals{marMUX: MARMUX_7_0, ldMAR: LoadSig_LOAD, gateMARMUX: NoYesSig_YES}},
	16: {COND: COND_MEMREADY, J: 18, DatapathSignals: Signals{rw: RW_WR, dataSize: WORD, mioEN: NoYesSig_YES}},
	17: {COND: COND_MEMREADY, J: 19, DatapathSignals: Signals{rw: RW_WR, mioEN: NoYesSig_YES, dataSize: BYTE}},
	18: {J: 33, DatapathSignals: Signals{gatePC: NoYesSig_YES, ldMAR: LoadSig_LOAD, pcMUX: PCPLUS2, ldPC: LoadSig_LOAD}},
	19: {J: 33, DatapathSignals: Signals{gatePC: NoYesSig_YES, ldMAR: LoadSig_LOAD, pcMUX: PCPLUS2, ldPC: LoadSig_LOAD}},
	20: {J: 18, DatapathSignals: Signals{drMUX: DRMUX_R7, gatePC: NoYesSig_YES, ldREG: LoadSig_LOAD, sr1MUX: SR1MUX_8_6, addr1MUX: ADDR1MUX_BaseR, pcMUX: ADDER, ldPC: LoadSig_LOAD}},
	21: {J: 18, DatapathSignals: Signals{drMUX: DRMUX_R7, gatePC: NoYesSig_YES, ldREG: LoadSig_LOAD, addr2MUX: ADDR2MUX_PCoffset11, lshf1: NoYesSig_YES, addr1MUX: ADDR1MUX_PC, pcMUX: ADDER, ldPC: LoadSig_LOAD}},
	22: {J: 18, DatapathSignals: Signals{addr2MUX: ADDR2MUX_PCoffset9, lshf1: NoYesSig_YES, addr1MUX: ADDR1MUX_PC, pcMUX: ADDER, ldPC: LoadSig_LOAD}},
	23: {J: 16, DatapathSignals: Signals{aluK: ALU_PASSA, dataSize: WORD, mioEN: NoYesSig_NO, ldMDR: LoadSig_LOAD, gateALU: NoYesSig_YES}},
	24: {J: 17, DatapathSignals: Signals{aluK: ALU_PASSA, dataSize: BYTE, mioEN: NoYesSig_NO, ldMDR: LoadSig_LOAD, gateALU: NoYesSig_YES}},
	25: {COND: COND_MEMREADY, J: 27, DatapathSignals: Signals{mioEN: NoYesSig_YES, ldMDR: LoadSig_LOAD, dataSize: WORD}},
	26: {},
	27: {J: 18, DatapathSignals: Signals{dataSize: WORD, gateMDR: NoYesSig_YES, ldREG: LoadSig_LOAD, ldCC: LoadSig_LOAD}},
	28: {COND: COND_MEMREADY, J: 30, DatapathSignals: Signals{mioEN: NoYesSig_YES, ldMDR: LoadSig_LOAD, drMUX: DRMUX_R7, gatePC: NoYesSig_YES, ldREG: LoadSig_LOAD}},
	29: {COND: COND_MEMREADY, J: 31, DatapathSignals: Signals{rw: RW_RD, dataSize: BYTE, mioEN: NoYesSig_YES, ldMDR: LoadSig_LOAD}},
	30: {J: 18, DatapathSignals: Signals{dataSize: WORD, gateMDR: NoYesSig_YES, pcMUX: BUS, ldPC: LoadSig_LOAD}},
	31: {J: 18, DatapathSignals: Signals{dataSize: BYTE, gateMDR: NoYesSig_YES, ldREG: LoadSig_LOAD, ldCC: LoadSig_LOAD}},
	32: {IRD: true, DatapathSignals: Signals{ldBEN: LoadSig_LOAD}},
	33: {COND: COND_MEMREADY, J: 35, DatapathSignals: Signals{rw: RW_RD, dataSize: WORD, mioEN: NoYesSig_YES, ldMDR: LoadSig_LOAD}},
	34: {},
	35: {J: 32, DatapathSignals: Signals{dataSize: WORD, gateMDR: NoYesSig_YES, ldIR: LoadSig_LOAD}},
}

func ControlStore(state uint8) ControlStoreOutput {

	return ControlStoreOutput{}
}
