package lib

import "fmt"

type UnknownInstruction struct {
	OpCode int
	Input1 int
	Input2 int
	Output int
}

func (i UnknownInstruction) ToString() string {
	return fmt.Sprintf("%d %3d %3d %d", i.OpCode, i.Input1, i.Input2, i.Output)
}

type Instruction struct {
	OpCode OpCode
	Input1 int
	Input2 int
	Output int
}

func (i Instruction) ToString() string {
	return fmt.Sprintf("%s %3d %3d %3d", i.OpCode, i.Input1, i.Input2, i.Output)
}

func (i Instruction) Execute(s State) State {
	result := s.Clone()
	r := result.Registers
	switch i.OpCode {
	case ADDR:
		r[i.Output] = r[i.Input1] + r[i.Input2]
	case ADDI:
		r[i.Output] = r[i.Input1] + i.Input2
	case MULR:
		r[i.Output] = r[i.Input1] * r[i.Input2]
	case MULI:
		r[i.Output] = r[i.Input1] * i.Input2
	case BANR:
		r[i.Output] = r[i.Input1] & r[i.Input2]
	case BANI:
		r[i.Output] = r[i.Input1] & i.Input2
	case BORR:
		r[i.Output] = r[i.Input1] | r[i.Input2]
	case BORI:
		r[i.Output] = r[i.Input1] | i.Input2
	case SETR:
		r[i.Output] = r[i.Input1]
	case SETI:
		r[i.Output] = i.Input1
	case GTIR:
		if i.Input1 > r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case GTRI:
		if r[i.Input1] > i.Input2 {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case GTRR:
		if r[i.Input1] > r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case EQIR:
		if i.Input1 == r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case EQRI:
		if r[i.Input1] == i.Input2 {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case EQRR:
		if r[i.Input1] == r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	default:
		panic(fmt.Sprintf("found unknown opcode '%s'", i.OpCode))
	}
	return result
}


func (i Instruction) ExecuteInto(s *State) *State {
	result := s.Clone()
	r := result.Registers
	switch i.OpCode {
	case ADDR:
		r[i.Output] = r[i.Input1] + r[i.Input2]
	case ADDI:
		r[i.Output] = r[i.Input1] + i.Input2
	case MULR:
		r[i.Output] = r[i.Input1] * r[i.Input2]
	case MULI:
		r[i.Output] = r[i.Input1] * i.Input2
	case BANR:
		r[i.Output] = r[i.Input1] & r[i.Input2]
	case BANI:
		r[i.Output] = r[i.Input1] & i.Input2
	case BORR:
		r[i.Output] = r[i.Input1] | r[i.Input2]
	case BORI:
		r[i.Output] = r[i.Input1] | i.Input2
	case SETR:
		r[i.Output] = r[i.Input1]
	case SETI:
		r[i.Output] = i.Input1
	case GTIR:
		if i.Input1 > r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case GTRI:
		if r[i.Input1] > i.Input2 {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case GTRR:
		if r[i.Input1] > r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case EQIR:
		if i.Input1 == r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case EQRI:
		if r[i.Input1] == i.Input2 {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case EQRR:
		if r[i.Input1] == r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	default:
		panic(fmt.Sprintf("found unknown opcode '%s'", i.OpCode))
	}
	return &result
}
