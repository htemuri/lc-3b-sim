package datapath

import "log"

const MEMORY_READ_LATENCY = 2  // in clock cycles
const MEMORY_WRITE_LATENCY = 5 // in clock cycles

type Memory struct {
	mem [65536]uint8 // 2**16 addresses with 8 bit addressability. words are 16bits and are aligned if their addresses differ only in bit [0]

	mar uint16 // memory address register

	Ready bool // memory ready signal

	pendingRead bool
	readCycles  uint8 // simulated latency for reads
	DataOut     uint16

	pendingWrite bool
	writeCycles  uint8 // simulated latency for writes
	writeEnable1 bool
	writeEnable0 bool
	writeData    uint16
}

func (m *Memory) Init(data [65536]uint8) {
	m.mem = data
}

func (m *Memory) Read(mar uint16) {
	if m.pendingRead {
		log.Printf("already pending read")
	}
	m.pendingRead = true
	m.mar = mar
	m.readCycles = MEMORY_READ_LATENCY
	m.Ready = false
}

func (m *Memory) Write(
	mar uint16,
	mdr uint16,
	writeEnable0,
	writeEnable1 bool,
) {
	if m.pendingWrite {
		log.Printf("already pending write")
	}
	m.pendingWrite = true
	m.mar = mar
	m.writeCycles = MEMORY_WRITE_LATENCY
	m.writeData = mdr
	m.Ready = false
	m.writeEnable0 = writeEnable0
	m.writeEnable1 = writeEnable1
}

func (m *Memory) Commit() {
	if m.pendingRead {
		if m.readCycles > 0 {
			m.readCycles--
		} else {
			m.DataOut = uint16(m.mem[m.mar+1])<<8 | uint16(m.mem[m.mar])
			m.pendingRead = false
			m.Ready = true
		}
	}
	if m.pendingWrite {
		if m.writeCycles > 0 {
			m.writeCycles--
		} else {
			// if both enabled then mar should be word alligned and instruction should be store word
			if m.writeEnable1 && m.writeEnable0 {
				m.mem[m.mar+1] = uint8(m.writeData >> 8)
				m.mem[m.mar] = uint8(m.writeData)
			} else if m.writeEnable1 { // if just WE1 enabled, write mdr[15:8] to mem[mar]
				m.mem[m.mar] = uint8(m.writeData >> 8)
			} else if m.writeEnable0 { // if just WE0 enabled, write mdr[7:0] to mem[mar]
				m.mem[m.mar] = uint8(m.writeData)
			}
			m.pendingWrite = false
			m.Ready = true
		}
	}
}
