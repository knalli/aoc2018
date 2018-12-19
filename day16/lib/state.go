package lib

type State struct {
	Registers []int
}

func (s State) Clone() State {
	registers := make([]int, len(s.Registers))
	copy(registers, s.Registers)
	return State{
		Registers: registers,
	}
}

func (s State) Equal(o State) bool {
	if len(s.Registers) != len(o.Registers) {
		return false
	}
	for i := range s.Registers {
		if s.Registers[i] != o.Registers[i] {
			return false
		}
	}
	return true
}
