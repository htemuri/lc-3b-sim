package control

type Signals struct {
	LdREG LoadSig
	LdCC  LoadSig
	LdPC  LoadSig
	LdIR  LoadSig
	LdMAR LoadSig
	LdMDR LoadSig
	LdBEN LoadSig

	GatePC     NoYesSig
	GateMDR    NoYesSig
	GateALU    NoYesSig
	GateMARMUX NoYesSig
	GateSHF    NoYesSig

	PcMUX    PCMux2
	DrMUX    DRMUX
	Sr1MUX   SR1MUX
	Addr1MUX ADDR1MUX
	Addr2MUX ADDR2MUX
	MarMUX   MARMUX
	AluK     ALUOp
	MioEN    NoYesSig
	Rw       RW
	DataSize DataSize
	Lshf1    NoYesSig
}

type LoadSig bool

const (
	LoadSig_NO   LoadSig = false
	LoadSig_LOAD LoadSig = true
)

type NoYesSig bool

const (
	NoYesSig_NO  NoYesSig = false
	NoYesSig_YES NoYesSig = true
)

type PCMux2 uint8

const (
	PCPLUS2 PCMux2 = iota
	BUS
	ADDER
)

type DataSize uint8

const (
	BYTE DataSize = iota
	WORD
)

type ALUOp uint8

const (
	ALU_ADD ALUOp = iota
	ALU_AND
	ALU_XOR
	ALU_PASSA
)

type SR2Mux bool

const (
	SR2Mux_SR2OUT SR2Mux = false
	SR2Mux_IR     SR2Mux = true
)

type DRMUX uint8

const (
	DRMUX_11_9 DRMUX = iota
	DRMUX_R7
)

type SR1MUX uint8

const (
	SR1MUX_11_9 SR1MUX = iota
	SR1MUX_8_6
)

type ADDR1MUX uint8

const (
	ADDR1MUX_PC ADDR1MUX = iota
	ADDR1MUX_BaseR
)

type ADDR2MUX uint8

const (
	ADDR2MUX_ZERO ADDR2MUX = iota
	ADDR2MUX_offset6
	ADDR2MUX_PCoffset9
	ADDR2MUX_PCoffset11
)

type MARMUX uint8

const (
	MARMUX_7_0 MARMUX = iota
	MARMUX_ADDER
)

type RW bool

const (
	RW_RD = false
	RW_WR = true
)
