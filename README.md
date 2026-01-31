# lc-3b-sim
Simulating a microarchitecture that satisfies the LC-3b ISA specification

The goal for this project is to learn how CPU microarchitecture while abstracting away transistor/latch/flip flop/logic gate design. 

## Todo:

> Adding to this list as I go

- [x] Build data path
  - [x] Adder
  - [x] ALU
  - [x] Register Interface
  - [x] Register File
  - [x] Memory
  - [x] Special registers
  - [x] Bit shifter
  - [x] Multiplexers
  - [x] Sign extender
- [x] Build control logic
  - [x] Microcode in the form of the control store
  - [x] Microsequencer


## LC-3b Microarchitecture 
![Diagram of an lc3b microarchitecture datapath](public/lc3b-microarch.png)


## LC-3b Finite State Machine
![Diagram of the lc3b finite state machine](public/lc3b-fsm.png)