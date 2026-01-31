package control

// The control store is basically taking the finite state machine which you can see in the README
// and converting it to the signals needed in the datapath to accomplish the state and move on.
// So for ex, if i'm on state 1...
//			- I need to do a microoperation which is adding sr1 + op2 and storing it in the dest register.
// 			  To do that, i need to set my alu control signal to add, select which bits from the instruction are
// 			  my sr1, and set the load signals for the registers im writing to. I also enable the tristate buffer
// 			  that allows my alu output to flow through the bus to my dest register and to my conditional code registers (n,z,p).
// 			- I set the `J` value to 18 to declare that once I finish my microoperation, I shouLd move to state 18

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

// control store shouLd be 35 bits x 2**6 (possible state combinations)
var store map[uint8]ControlStoreOutput = map[uint8]ControlStoreOutput{
	0:  {COND: COND_BRANCH, J: 18},
	1:  {J: 18, DatapathSignals: Signals{AluK: ALU_ADD, Sr1MUX: SR1MUX_8_6, LdCC: LoadSig_LOAD, LdREG: LoadSig_LOAD, GateALU: NoYesSig_YES}},
	2:  {J: 29, DatapathSignals: Signals{Sr1MUX: SR1MUX_8_6, Addr2MUX: ADDR2MUX_offset6, MarMUX: MARMUX_ADDER, Addr1MUX: ADDR1MUX_BaseR, LdMAR: LoadSig_LOAD, GateMARMUX: NoYesSig_YES}},
	3:  {J: 24, DatapathSignals: Signals{Sr1MUX: SR1MUX_8_6, Addr2MUX: ADDR2MUX_offset6, MarMUX: MARMUX_ADDER, Addr1MUX: ADDR1MUX_BaseR, LdMAR: LoadSig_LOAD, GateMARMUX: NoYesSig_YES}},
	4:  {COND: COND_ADDRMODE, J: 20},
	5:  {J: 18, DatapathSignals: Signals{AluK: ALU_AND, Sr1MUX: SR1MUX_8_6, LdCC: LoadSig_LOAD, LdREG: LoadSig_LOAD, GateALU: NoYesSig_YES}},
	6:  {J: 25, DatapathSignals: Signals{Sr1MUX: SR1MUX_8_6, Addr2MUX: ADDR2MUX_offset6, MarMUX: MARMUX_ADDER, Addr1MUX: ADDR1MUX_BaseR, LdMAR: LoadSig_LOAD, Lshf1: NoYesSig_YES, GateMARMUX: NoYesSig_YES}},
	7:  {J: 23, DatapathSignals: Signals{Sr1MUX: SR1MUX_8_6, Addr2MUX: ADDR2MUX_offset6, MarMUX: MARMUX_ADDER, Addr1MUX: ADDR1MUX_BaseR, LdMAR: LoadSig_LOAD, Lshf1: NoYesSig_YES, GateMARMUX: NoYesSig_YES}},
	8:  {},
	9:  {J: 18, DatapathSignals: Signals{AluK: ALU_XOR, Sr1MUX: SR1MUX_8_6, LdCC: LoadSig_LOAD, LdREG: LoadSig_LOAD, GateALU: NoYesSig_YES}},
	10: {},
	11: {},
	12: {J: 18, DatapathSignals: Signals{PcMUX: ADDER, Addr1MUX: ADDR1MUX_BaseR, LdPC: LoadSig_LOAD, Sr1MUX: SR1MUX_8_6}},
	13: {J: 18, DatapathSignals: Signals{Sr1MUX: SR1MUX_8_6, LdREG: LoadSig_LOAD, GateSHF: NoYesSig_YES, LdCC: LoadSig_LOAD}},
	14: {J: 18, DatapathSignals: Signals{Addr1MUX: ADDR1MUX_PC, Addr2MUX: ADDR2MUX_PCoffset9, Lshf1: NoYesSig_YES, MarMUX: MARMUX_ADDER, LdREG: LoadSig_LOAD, LdCC: LoadSig_LOAD, GateMARMUX: NoYesSig_YES}},
	15: {J: 28, DatapathSignals: Signals{MarMUX: MARMUX_7_0, LdMAR: LoadSig_LOAD, GateMARMUX: NoYesSig_YES}},
	16: {COND: COND_MEMREADY, J: 18, DatapathSignals: Signals{Rw: RW_WR, DataSize: WORD, MioEN: NoYesSig_YES}},
	17: {COND: COND_MEMREADY, J: 19, DatapathSignals: Signals{Rw: RW_WR, MioEN: NoYesSig_YES, DataSize: BYTE}},
	18: {J: 33, DatapathSignals: Signals{GatePC: NoYesSig_YES, LdMAR: LoadSig_LOAD, PcMUX: PCPLUS2, LdPC: LoadSig_LOAD}},
	19: {J: 33, DatapathSignals: Signals{GatePC: NoYesSig_YES, LdMAR: LoadSig_LOAD, PcMUX: PCPLUS2, LdPC: LoadSig_LOAD}},
	20: {J: 18, DatapathSignals: Signals{DrMUX: DRMUX_R7, GatePC: NoYesSig_YES, LdREG: LoadSig_LOAD, Sr1MUX: SR1MUX_8_6, Addr1MUX: ADDR1MUX_BaseR, PcMUX: ADDER, LdPC: LoadSig_LOAD}},
	21: {J: 18, DatapathSignals: Signals{DrMUX: DRMUX_R7, GatePC: NoYesSig_YES, LdREG: LoadSig_LOAD, Addr2MUX: ADDR2MUX_PCoffset11, Lshf1: NoYesSig_YES, Addr1MUX: ADDR1MUX_PC, PcMUX: ADDER, LdPC: LoadSig_LOAD}},
	22: {J: 18, DatapathSignals: Signals{Addr2MUX: ADDR2MUX_PCoffset9, Lshf1: NoYesSig_YES, Addr1MUX: ADDR1MUX_PC, PcMUX: ADDER, LdPC: LoadSig_LOAD}},
	23: {J: 16, DatapathSignals: Signals{AluK: ALU_PASSA, DataSize: WORD, MioEN: NoYesSig_NO, LdMDR: LoadSig_LOAD, GateALU: NoYesSig_YES}},
	24: {J: 17, DatapathSignals: Signals{AluK: ALU_PASSA, DataSize: BYTE, MioEN: NoYesSig_NO, LdMDR: LoadSig_LOAD, GateALU: NoYesSig_YES}},
	25: {COND: COND_MEMREADY, J: 27, DatapathSignals: Signals{MioEN: NoYesSig_YES, LdMDR: LoadSig_LOAD, DataSize: WORD}},
	26: {},
	27: {J: 18, DatapathSignals: Signals{DataSize: WORD, GateMDR: NoYesSig_YES, LdREG: LoadSig_LOAD, LdCC: LoadSig_LOAD}},
	28: {COND: COND_MEMREADY, J: 30, DatapathSignals: Signals{MioEN: NoYesSig_YES, LdMDR: LoadSig_LOAD, DrMUX: DRMUX_R7, GatePC: NoYesSig_YES, LdREG: LoadSig_LOAD}},
	29: {COND: COND_MEMREADY, J: 31, DatapathSignals: Signals{Rw: RW_RD, DataSize: BYTE, MioEN: NoYesSig_YES, LdMDR: LoadSig_LOAD}},
	30: {J: 18, DatapathSignals: Signals{DataSize: WORD, GateMDR: NoYesSig_YES, PcMUX: BUS, LdPC: LoadSig_LOAD}},
	31: {J: 18, DatapathSignals: Signals{DataSize: BYTE, GateMDR: NoYesSig_YES, LdREG: LoadSig_LOAD, LdCC: LoadSig_LOAD}},
	32: {IRD: true, DatapathSignals: Signals{LdBEN: LoadSig_LOAD}},
	33: {COND: COND_MEMREADY, J: 35, DatapathSignals: Signals{Rw: RW_RD, DataSize: WORD, MioEN: NoYesSig_YES, LdMDR: LoadSig_LOAD}},
	34: {},
	35: {J: 32, DatapathSignals: Signals{DataSize: WORD, GateMDR: NoYesSig_YES, LdIR: LoadSig_LOAD}},
	// initing 36..64 to prevent out of bounds state index
	36: {}, 37: {}, 38: {}, 39: {}, 40: {}, 41: {}, 42: {}, 43: {}, 44: {}, 45: {}, 46: {}, 47: {}, 48: {}, 49: {}, 50: {}, 51: {}, 52: {},
	53: {}, 54: {}, 55: {}, 56: {}, 57: {}, 58: {}, 59: {}, 60: {}, 61: {}, 62: {}, 63: {},
}

func ControlStore(state uint8) ControlStoreOutput {
	// state shouLd max be 6 bits so going to remove top 2 bits from var
	state = (state << 2) >> 2
	return store[state] // max state can be is 2**6 - 1 = 63
}
