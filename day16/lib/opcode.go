package lib

import "fmt"

type OpCode string

const (
	ADDR OpCode = "addr"
	ADDI OpCode = "addi"
	MULR OpCode = "mulr"
	MULI OpCode = "muli"
	BANR OpCode = "banr"
	BANI OpCode = "bani"
	BORR OpCode = "borr"
	BORI OpCode = "bori"
	SETR OpCode = "setr"
	SETI OpCode = "seti"
	GTIR OpCode = "gtir"
	GTRI OpCode = "gtri"
	GTRR OpCode = "gtrr"
	EQIR OpCode = "eqir"
	EQRI OpCode = "eqri"
	EQRR OpCode = "eqrr"
)

func (o OpCode) ToString() string {
	return fmt.Sprintf("%s", o)
}

var ALL_OPCODES = []OpCode{
	ADDR, ADDI,
	MULR, MULI,
	BANR, BANI,
	BORR, BORI,
	SETR, SETI,
	GTIR, GTRI, GTRR,
	EQIR, EQRI, EQRR,
}
