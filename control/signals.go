package control

type Signals struct {
	ldREG LoadSig
	ldCC  LoadSig
	ldPC  LoadSig
	ldIR  LoadSig
	ldMAR LoadSig
	ldMDR LoadSig
	ldBEN LoadSig

	gatePC     NoYesSig
	gateMDR    NoYesSig
	gateALU    NoYesSig
	gateMARMUX NoYesSig
	gateSHF    NoYesSig

	pcMUX    PCMux2
	drMUX    DRMUX
	sr1MUX   SR1MUX
	addr1MUX ADDR1MUX
	addr2MUX ADDR2MUX
	marMUX   MARMUX
	aluK     ALUOp
	mioEN    NoYesSig
	rw       RW
	dataSize DataSize
	lshf1    NoYesSig
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
	WORD DataSize = iota
	BYTE
)

type ALUOp uint8

const (
	ALU_ADD ALUOp = iota
	ALU_AND
	ALU_XOR
	ALU_PASSA
)

type SR2Mux uint8

const (
	IR SR2Mux = iota
	SR2OUT
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
