package datapath

import "lc3b-sim/m/v2/control"

// general purpose registers - 8x16 bit

type GPRegister uint8

const (
	R0 GPRegister = iota
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

	sr1Out uint16
	sr2Out uint16

	// adding these because writes are sequential and depend on the clock
	pendingWrite bool
	writeData    uint16
	writeReg     GPRegister
}

func (rf *RegisterFile) Read(
	sr1, sr2 GPRegister,
) {
	rf.sr1Out = rf.regs[sr1]
	rf.sr2Out = rf.regs[sr2]
}

func (rf *RegisterFile) Write(
	ldREG control.LD_REG,
	dr GPRegister,
	data uint16,
) {
	if ldREG {
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

type PC struct {
	value    uint16
	pending  bool
	newValue uint16
}

func (pc *PC) GetPC() uint16 {
	return pc.value
}

func (pc *PC) UpdatePC(ldPC control.LD_PC, address uint16) {
	if ldPC {
		pc.newValue = address
		pc.pending = true
	}
}

func (pc *PC) Commit() {
	if pc.pending {
		pc.value = pc.newValue
		pc.pending = false
	}
}

type MAR struct {
	value        uint16
	nextValue    uint16
	pendingWrite bool
}

func (m *MAR) GetMAR() uint16 {
	return m.value
}

func (m *MAR) UpdateMAR(ldMAR control.LD_MAR, address uint16) {
	if ldMAR {
		m.nextValue = address
		m.pendingWrite = true
	}
}

func (m *MAR) Commit() {
	if m.pendingWrite {
		m.value = m.nextValue
		m.pendingWrite = false
	}
}
