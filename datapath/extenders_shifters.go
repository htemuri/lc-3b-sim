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

func SEXT(input uint16) uint16 { // extend most sig bit of input till 16th bit of output (for ex: 110 -> 1111111111111110)
	for i := 16; i > 0; i-- {
		var bitmask uint16 = 1 << (i - 1)
		if bitmask&input == bitmask {
			// cant do arithmetic shift bc working with unsigned ints so will have to bitwise OR the bits
			var extenders uint16 = 0b1111111111111111
			extenders = (extenders >> i) << i
			return extenders | input
		}
	}
	return 0
}
