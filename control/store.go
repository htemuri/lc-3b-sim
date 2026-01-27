package control

// control store should be 35 bits x 2**6 (possible state combinations)

type ControlStore struct {
	store map[int][]int
}

type ControlStoreOutput struct {
	IRD             bool
	COND            uint8
	J               uint8
	DatapathSignals Signals
}

func (cs *ControlStore) GetSignals(state uint8) ControlStoreOutput {
	return ControlStoreOutput{}
}
