package datapath

import (
	"lc3b-sim/m/v2/control"
	"log"
)

// general purpose registers - 8x16 bit

type gpRegister uint8

const (
	R0 gpRegister = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
)

type RegisterFile struct {
	regs [8]uint16

	Sr1Out uint16
	Sr2Out uint16

	// adding these because writes are sequential and depend on the clock
	pendingWrite bool
	writeData    uint16
	writeReg     gpRegister
}

func (rf *RegisterFile) Read(
	sr1, sr2 gpRegister,
) {
	rf.Sr1Out = rf.regs[sr1]
	rf.Sr2Out = rf.regs[sr2]
}

func (rf *RegisterFile) Write(
	ldREG control.LoadSig,
	dr gpRegister,
	data uint16,
) {
	if ldREG == control.LoadSig_LOAD {
		rf.pendingWrite = true
		rf.writeData = data
		rf.writeReg = dr
	}
}

func (rf *RegisterFile) Commit() {
	if rf.pendingWrite {
		rf.regs[rf.writeReg] = rf.writeData
		rf.pendingWrite = false
	}
}

func (rf *RegisterFile) PrintValues() {
	log.Println("REGISTER FILE: ", rf.regs)
}

func GetSR2Input(fullIR uint16) gpRegister {
	return gpRegister(fullIR & 0b111)
}

// generic 16 bit register
type Register16bit struct {
	value        uint16
	pendingWrite bool
	newValue     uint16
}

func (r *Register16bit) GetValue() uint16 {
	return r.value
}

func (r *Register16bit) UpdateValue(
	loadSignal control.LoadSig,
	value uint16,
) {
	if loadSignal == control.LoadSig_LOAD {
		r.pendingWrite = true
		r.newValue = value
	}
}

func (r *Register16bit) Commit() {
	if r.pendingWrite {
		r.value = r.newValue
		r.pendingWrite = false
	}
}

type ConditionalCodes struct {
	N bool
	Z bool
	P bool

	pendingUpdate bool
	newN          bool
	newZ          bool
	newP          bool
}

func (cc *ConditionalCodes) SetCC(loadCC control.LoadSig, drValue uint16) {
	if loadCC == control.LoadSig_LOAD {
		cc.newN, cc.newZ, cc.newP = false, false, false
		if (drValue >> 15) == 1 { // assuming 2s complement form for values, if most sig bit is 1, then its neg
			cc.newN = true
		} else if drValue == 0 {
			cc.newZ = true
		} else {
			cc.newP = true
		}
		cc.pendingUpdate = true
	}
}

func (cc *ConditionalCodes) Commit() {
	if cc.pendingUpdate {
		cc.N = cc.newN
		cc.Z = cc.newZ
		cc.P = cc.newP
		cc.pendingUpdate = false
	}
}
