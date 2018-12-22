package lib

import "fmt"

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

func (s State) ToString() string {
	result := ""
	for i, v := range s.Registers {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%d", v)
	}
	return result
}
