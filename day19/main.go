package main

import (
	d16 "de.knallisworld/aoc/aoc2018/day16/lib"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 19
const AocDayName = "day19"
const trace = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	pcRegister, instructions := readInstructions(lines)

	dayless.PrintStepHeader(1)
	result := compute(instructions, pcRegister, false)
	fmt.Printf("🎉The final result state of the registers: 0=%d, 1=%d, 2=%d, 3=%d, 4=%d, 5=%d\n",
		result.Registers[0],
		result.Registers[1],
		result.Registers[2],
		result.Registers[3],
		result.Registers[4],
		result.Registers[5],
	)
	fmt.Println()

	dayless.PrintStepHeader(2)
	result = compute(instructions, pcRegister, true)
	fmt.Printf("🎉The final result state of the registers: 0=%d, 1=%d, 2=%d, 3=%d, 4=%d, 5=%d\n",
		result.Registers[0],
		result.Registers[1],
		result.Registers[2],
		result.Registers[3],
		result.Registers[4],
		result.Registers[5],
	)
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

func compute(instructions []d16.Instruction, ip int, part2 bool) d16.State {
	defer dayless.TimeTrack(time.Now(), "compute")

	state := &d16.State{
		Registers: []int{0, 0, 0, 0, 0, 0},
	}
	if part2 {
		state.Registers[0] = 1
	}

	for ; state.Registers[ip] < len(instructions); {
		instruction := instructions[state.Registers[ip]]

		// optimize
		if part2 && state.Registers[ip] == 2 {
			for i := 1; i <= state.Registers[3]; i++ {
				if state.Registers[3]%i == 0 {
					state.Registers[0] += i
				}
			}

			state.Registers[ip] = 15
		} else {
			next := instruction.ExecuteInto(state)

			if trace {
				fmt.Printf("👉 ip=%2d [%4d, %4d, %4d, %4d, %4d, %4d] %s [%4d, %4d, %4d, %4d, %4d, %4d]\n",
					state.Registers[ip],
					state.Registers[0],
					state.Registers[1],
					state.Registers[2],
					state.Registers[3],
					state.Registers[4],
					state.Registers[5],
					instruction.ToString(),
					next.Registers[0],
					next.Registers[1],
					next.Registers[2],
					next.Registers[3],
					next.Registers[4],
					next.Registers[5],
				)
			}
			state = next
		}

		state.Registers[ip]++
	}

	return *state
}
