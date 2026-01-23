package main

import "fmt"

// using 16 bit unsigned ints because each word in lc-3b is 16 bits and we're assuming 2s complement form

func main() {
	x := uint16(0b0101011011111111)
	fmt.Printf("%016b", uint8(x>>8))
}
