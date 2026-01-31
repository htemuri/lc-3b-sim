package control

func Microsequencer(
	cond Condition,
	ben bool,
	r bool,
	ir15to11 uint8, // IR[15:11]
	j uint8, // should actually be uint6
	ird bool,
) uint8 { // only care about the first 6 bits in output

	if ird {
		return ir15to11 >> 1 // IR[15:12]
	}

	cond1 := uint8((0b10&cond)>>1) == 1
	cond0 := uint8(0b1&cond) == 1
	ir11 := 0b1&ir15to11 == 1

	var branch, ready, addrMode uint8
	if cond1 && !cond0 && ben {
		branch = 1 << 2
	}
	if !cond1 && cond0 && r {
		ready = 1 << 1
	}
	if cond1 && cond0 && ir11 {
		addrMode = 1
	}

	newJ2 := (0b100 & j) | branch
	newJ1 := (0b10 & j) | ready
	newJ0 := (0b1 & j) | addrMode
	newJ := ((j >> 3) << 3) | newJ2 | newJ1 | newJ0

	return newJ
}
