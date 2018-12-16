package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const AocDay = 16
const AocDayName = "day16"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	samples, lastLine := readSamples(lines)
	fmt.Println(samples)

	dayless.PrintStepHeader(1)
	checkInstructionsForOpCodeBehaviours(samples, 3)
	fmt.Println()

	dayless.PrintStepHeader(2)
	result := compute(readInstructions(lines, lastLine, resolveOpCodes(samples)))
	fmt.Printf("🎉The final result state of the registers: 0=%d, 1=%d, 2=%d, 3=%d\n",
		result.Registers[0], result.Registers[1], result.Registers[2], result.Registers[3])
	fmt.Println()
}

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

type UnknownInstruction struct {
	OpCode int
	Input1 int
	Input2 int
	Output int
}

func (i UnknownInstruction) ToString() string {
	return fmt.Sprintf("%d %d %d %d", i.OpCode, i.Input1, i.Input2, i.Output)
}

type Instruction struct {
	OpCode OpCode
	Input1 int
	Input2 int
	Output int
}

func (i Instruction) ToString() string {
	return fmt.Sprintf("%s %d %d %d", i.OpCode, i.Input1, i.Input2, i.Output)
}

type OpCode string

const (
	addr OpCode = "addr"
	addi OpCode = "addi"
	mulr OpCode = "mulr"
	muli OpCode = "muli"
	banr OpCode = "banr"
	bani OpCode = "bani"
	borr OpCode = "borr"
	bori OpCode = "bori"
	setr OpCode = "setr"
	seti OpCode = "seti"
	gtir OpCode = "gtir"
	gtri OpCode = "gtri"
	gtrr OpCode = "gtrr"
	eqir OpCode = "eqir"
	eqri OpCode = "eqri"
	eqrr OpCode = "eqrr"
)

type Sample struct {
	Before      State
	Instruction UnknownInstruction
	After       State
}

type Samples []Sample

func readSamples(lines []string) (Samples, int) {

	bPattern := regexp.MustCompile("Before:\\s+\\[(\\d+), (\\d+), (\\d+), (\\d+)]")
	aPattern := regexp.MustCompile("After:\\s+\\[(\\d+), (\\d+), (\\d+), (\\d+)]")

	result := make(Samples, 0)
	i := 0
	for ; i < len(lines); i += 4 {
		if lines[i] == "" {
			break // empty lines: content of part2 begins
		}
		bMatch := bPattern.FindStringSubmatch(lines[i])
		aMatch := aPattern.FindStringSubmatch(lines[i+2])
		sample := Sample{
			Before:      extractState(bMatch),
			Instruction: extractUnknownInstruction(lines[i+1]),
			After:       extractState(aMatch),
		}
		result = append(result, sample)
	}
	return result, i
}

func extractState(matches []string) State {
	registers := make([]int, 4)
	if len(matches) != 5 {
		panic("unexpected size of matches (should be 5)")
	}
	for i := 0; i < 4; i++ {
		registers[i] = parseInt(matches[i+1])
	}
	return State{
		Registers: registers,
	}
}

func extractUnknownInstruction(line string) UnknownInstruction {
	parts := strings.Split(line, " ")

	return UnknownInstruction{
		OpCode: parseInt(parts[0]),
		Input1: parseInt(parts[1]),
		Input2: parseInt(parts[2]),
		Output: parseInt(parts[3]),
	}
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func executeInstruction(s State, i Instruction) State {
	result := s.Clone()
	r := result.Registers
	switch i.OpCode {
	case addr:
		r[i.Output] = r[i.Input1] + r[i.Input2]
	case addi:
		r[i.Output] = r[i.Input1] + i.Input2
	case mulr:
		r[i.Output] = r[i.Input1] * r[i.Input2]
	case muli:
		r[i.Output] = r[i.Input1] * i.Input2
	case banr:
		r[i.Output] = r[i.Input1] & r[i.Input2]
	case bani:
		r[i.Output] = r[i.Input1] & i.Input2
	case borr:
		r[i.Output] = r[i.Input1] | r[i.Input2]
	case bori:
		r[i.Output] = r[i.Input1] | i.Input2
	case setr:
		r[i.Output] = r[i.Input1]
	case seti:
		r[i.Output] = i.Input1
	case gtir:
		if i.Input1 > r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case gtri:
		if r[i.Input1] > i.Input2 {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case gtrr:
		if r[i.Input1] > r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case eqir:
		if i.Input1 == r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case eqri:
		if r[i.Input1] == i.Input2 {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	case eqrr:
		if r[i.Input1] == r[i.Input2] {
			r[i.Output] = 1
		} else {
			r[i.Output] = 0
		}
	default:
		panic("this opcode is unknown: " + i.OpCode)
	}
	return result
}

func resolveOpCodeBehaviours(before State, after State, instruction UnknownInstruction) []OpCode {
	matchingOpCodes := make([]OpCode, 0)

	for _, oc := range []OpCode{
		addr, addi,
		mulr, muli,
		banr, bani,
		borr, bori,
		setr, seti,
		gtir, gtri, gtrr,
		eqir, eqri, eqrr,
	} {
		testInstruction := Instruction{
			OpCode: oc,
			Input1: instruction.Input1,
			Input2: instruction.Input2,
			Output: instruction.Output,
		}
		if executeInstruction(before, testInstruction).Equal(after) {
			matchingOpCodes = append(matchingOpCodes, oc)
		}
	}

	return matchingOpCodes
}

func checkInstructionsForOpCodeBehaviours(samples Samples, threshold int) {
	total := 0
	for _, sample := range samples {
		behaviours := resolveOpCodeBehaviours(sample.Before, sample.After, sample.Instruction)
		opCodeTotal := len(behaviours)
		if opCodeTotal >= threshold {
			fmt.Printf("👉 Instruction '%s' behave like %d opcodes: %s\n", sample.Instruction.ToString(), opCodeTotal, behaviours)
			total++
		}
	}

	fmt.Printf("🎉 Totally, %d samples behave like %d or more opcodes\n", total, threshold)
}

func resolveOpCodes(samples Samples) (result map[OpCode]int) {
	result = make(map[OpCode]int)
	ocCandidates := make(map[OpCode][]int)
	for _, oc := range []OpCode{
		addr, addi,
		mulr, muli,
		banr, bani,
		borr, bori,
		setr, seti,
		gtir, gtri, gtrr,
		eqir, eqri, eqrr,
	} {
		ocCandidates[oc] = []int{}
	}
	for _, sample := range samples {
		for _, oc := range resolveOpCodeBehaviours(sample.Before, sample.After, sample.Instruction) {
			opCode := sample.Instruction.OpCode
			if !contains(ocCandidates[oc], opCode) {
				ocCandidates[oc] = append(ocCandidates[oc], opCode)
			}
		}
	}
	reduceCandidatesIfMatchingOne(ocCandidates)
	fmt.Printf("👉 Resolved all op codes:\n")
	for operation, operationCodes := range ocCandidates {
		fmt.Printf("  %s: %2d\n", operation, operationCodes[0])
		result[operation] = operationCodes[0]
	}

	return result
}

func reduceCandidatesIfMatchingOne(operations map[OpCode][]int) {

	uniqueOperations := make(map[OpCode]int)

	for {
		for operation, operationCodes := range operations {
			if len(operationCodes) == 1 {
				uniqueOperations[operation] = operationCodes[0]
			}
		}
		// filtering this
		for uniqueOperation, uniqueOperationCode := range uniqueOperations {
			for operation := range operations {
				operationCodes := operations[operation]
				if uniqueOperation != operation {
					operations[operation] = remove(operationCodes, uniqueOperationCode)
				}
			}
		}
		if len(uniqueOperations) == len(operations) {
			break
		}
	}
}

func readInstructions(lines []string, linesOffset int, opcodeMap map[OpCode]int) (result []Instruction) {
	// flip map
	codeMap := make(map[int]OpCode)
	for op, code := range opcodeMap {
		codeMap[code] = op
	}

	for i := linesOffset; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		result = append(result, Instruction{
			OpCode: codeMap[parseInt(parts[0])],
			Input1: parseInt(parts[1]),
			Input2: parseInt(parts[2]),
			Output: parseInt(parts[3]),
		})
	}

	return result
}

func compute(instructions []Instruction) State {
	state := State{
		Registers: []int{0, 0, 0, 0},
	}

	for _, instruction := range instructions {
		state = executeInstruction(state, instruction)
	}

	return state
}

func contains(array []int, search int) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}

func remove(array []int, needle int) (result []int) {
	for _, v := range array {
		if v != needle {
			result = append(result, v)
		}
	}
	return result
}
