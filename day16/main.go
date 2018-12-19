package main

import (
	"de.knallisworld/aoc/aoc2018/day16/lib"
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

type Sample struct {
	Before      lib.State
	Instruction lib.UnknownInstruction
	After       lib.State
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

func extractState(matches []string) lib.State {
	registers := make([]int, 4)
	if len(matches) != 5 {
		panic("unexpected size of matches (should be 5)")
	}
	for i := 0; i < 4; i++ {
		registers[i] = parseInt(matches[i+1])
	}
	return lib.State{
		Registers: registers,
	}
}

func extractUnknownInstruction(line string) lib.UnknownInstruction {
	parts := strings.Split(line, " ")

	return lib.UnknownInstruction{
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

func resolveOpCodeBehaviours(before lib.State, after lib.State, instruction lib.UnknownInstruction) []lib.OpCode {
	matchingOpCodes := make([]lib.OpCode, 0)

	for _, oc := range lib.ALL_OPCODES {
		testInstruction := lib.Instruction{
			OpCode: oc,
			Input1: instruction.Input1,
			Input2: instruction.Input2,
			Output: instruction.Output,
		}
		if testInstruction.Execute(before).Equal(after) {
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

func resolveOpCodes(samples Samples) (result map[lib.OpCode]int) {
	result = make(map[lib.OpCode]int)
	ocCandidates := make(map[lib.OpCode][]int)
	for _, oc := range lib.ALL_OPCODES {
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

func reduceCandidatesIfMatchingOne(operations map[lib.OpCode][]int) {

	uniqueOperations := make(map[lib.OpCode]int)

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

func readInstructions(lines []string, linesOffset int, opcodeMap map[lib.OpCode]int) (result []lib.Instruction) {
	// flip map
	codeMap := make(map[int]lib.OpCode)
	for op, code := range opcodeMap {
		codeMap[code] = op
	}

	for i := linesOffset; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		result = append(result, lib.Instruction{
			OpCode: codeMap[parseInt(parts[0])],
			Input1: parseInt(parts[1]),
			Input2: parseInt(parts[2]),
			Output: parseInt(parts[3]),
		})
	}

	return result
}

func compute(instructions []lib.Instruction) lib.State {
	state := lib.State{
		Registers: []int{0, 0, 0, 0},
	}

	for _, instruction := range instructions {
		state = instruction.Execute(state)
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
