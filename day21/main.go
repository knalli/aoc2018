package main

import (
	d16 "de.knallisworld/aoc/aoc2018/day16/lib"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 21
const AocDayName = "day21"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	pcRegister, instructions := readInstructions(lines)

	dayless.PrintStepHeader(1)
	compute(instructions, pcRegister, true)
	fmt.Println()

	dayless.PrintStepHeader(2)
	compute(instructions, pcRegister, false)
	fmt.Println()
}

func readInstructions(lines []string) (ip int, result []d16.Instruction) {
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if line[0:4] == "#ip " {
			ip = dayless.ParseInt(line[4:])
		} else {
			parts := strings.Split(line, " ")
			result = append(result, d16.Instruction{
				OpCode: parseOpCodeByString(parts[0]),
				Input1: dayless.ParseInt(parts[1]),
				Input2: dayless.ParseInt(parts[2]),
				Output: dayless.ParseInt(parts[3]),
			})
		}
	}

	return ip, result
}

func parseOpCodeByString(str string) d16.OpCode {
	s := strings.ToLower(str)
	for _, oc := range d16.ALL_OPCODES {
		a := oc.ToString()
		if a == s {
			return oc
		}
	}
	panic(fmt.Sprintf("unknown opcode '%s'", str))
}

func compute(instructions []d16.Instruction, ip int, part1 bool) (*d16.State, int) {
	defer dayless.TimeTrack(time.Now(), "compute")

	state := &d16.State{
		Registers: []int{0, 0, 0, 0, 0, 0},
	}

	set := make(map[int]bool)
	last := 0

	i := 0
	for ; state.Registers[ip] < len(instructions); {

		instruction := instructions[state.Registers[ip]]
		state = instruction.ExecuteInto(state)

		if part1 {
			// If the first time here, it is the first number we are looking for...
			if state.Registers[ip] == 28 {
				fmt.Printf("👉 Program would be halted with R0=%d\n", state.Registers[2])
				break
			}
		} else {
			if state.Registers[ip] == 28 {
				if _, exist := set[state.Registers[2]]; exist {
					fmt.Printf("👉 Program would be halted with R0=%d, and this has been found twice now. Last:  %d\n", state.Registers[2], last)
					break
				} else {
					set[state.Registers[2]] = true
					last = state.Registers[2]
				}
			}
		}

		state.Registers[ip]++
		i++
	}

	return state, i
}
