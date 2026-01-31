package main

import (
	"fmt"
	"lc3b-sim/m/v2/datapath"
)

// using 16 bit unsigned ints because each word in lc-3b is 16 bits and we're assuming 2s complement form

func main() {
	x := uint16(0b1010100101000111)
	fmt.Printf("%016b", datapath.SHF(x, 0b00110001))
}
