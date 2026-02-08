package datapath

import (
	"lc3b-sim/m/v2/control"
)

func ZEXTandLSHF1(input uint8) uint16 {
	return uint16(input) << 1
}

func LSHF1(lshf1 control.NoYesSig, input uint16) uint16 {
	if lshf1 == control.NoYesSig_NO {
		return input
	} else {
		return input << 1
	}
}

func SEXT(input uint16, numBits uint8) uint16 { // extend most sig bit of input till 16th bit of output (for ex: 110 -> 1111111111111110)

	msb := 0b1 & (input >> (numBits - 1))
	if msb == 1 {
		var extenders uint16 = 0b1111111111111111
		extenders = (extenders >> numBits) << numBits
		return extenders | input
	}
	return uint16(input)

	// for i := 16; i > 0; i-- {
	// 	var bitmask uint16 = 1 << (i - 1)
	// 	if bitmask&input == bitmask {
	// 		// cant do arithmetic shift bc working with unsigned ints so will have to bitwise OR the bits
	// 		var extenders uint16 = 0b1111111111111111
	// 		extenders = (extenders >> i) << i
	// 		return extenders | input
	// 	}
	// }
	// return 0
}

func SHF(sr1 uint16, op uint8) uint16 {
	amount := (op << 4) >> 4
	if 0b00010000&op == 0b00000000 { // if bit 4 == 0 then shift left
		// lshift sr1 by first 4 bits amount
		return sr1 << amount
	} else {
		if 0b00100000&op == 0 { // if bit 5 == 1 then arithmetic shift right
			// logical right shift sr1 by first 4 bits amount
			return sr1 >> amount
		} else {
			// arithmetic right shift sr1 by first 4 bits amount
			bit15 := sr1 >> 15
			var appendMask uint16 = ((bit15 * 0b1111111111111111) >> (16 - amount)) << (16 - amount)
			return (sr1 >> amount) | appendMask

		}
	}
}
