package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 12
const AocDayName = "day12"
const debug = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	state, rules := readPuzzle(lines)
	runNextGenerations(20, state, rules)
	fmt.Printf("🎉 The sum of the numbers of all pots which contain a plant: %d\n", sumOfFilledItemIndexes(state))
	fmt.Println()

	dayless.PrintStepHeader(2)
	// TODO after generation ~166, this are adding 73 for each constantly, running for 50000000000 is too slow
	// runNextGenerations(50000000000-20, state, rules)
	runNextGenerations(180, state, rules)
	fmt.Printf("🎉 The sum of the numbers of all pots which contain a plant: %d\n", sumOfFilledItemIndexes(state)+((50000000000-200)*73))
	fmt.Println()
}

type state struct {
	value  string
	offset int
}

type rule struct {
	pattern string
	filled  string
}

func readPuzzle(lines []string) (*state, []rule) {
	// initial state:
	initialState := strings.Split(lines[0], " ")[2]

	rules := make([]rule, len(lines)-2)
	for i := 2; i < len(lines); i++ {
		// #.### => .
		parts := strings.Split(lines[i], " ")
		rules[i-2] = rule{
			pattern: parts[0],
			filled:  parts[2],
		}
	}

	return &state{value: initialState, offset: 0}, rules
}

type memory struct {
	value      string
	generation int64
}

func runNextGenerations(generations int64, state *state, rules []rule) {
	defer dayless.TimeTrack(time.Now(), fmt.Sprintf("runNextGenerations=%d", generations))

	generationSums := make([]int64, 0)
	generationSums = append(generationSums, sumOfFilledItemIndexes(state))

	for g := int64(0); g < generations; g++ {
		runNextGeneration(state, rules)
		generationSums = append(generationSums, sumOfFilledItemIndexes(state))
		if debug {
			printState(g+1, state)
		}
	}

	for i, x := range generationSums {
		p := int64(0)
		if i > 0 {
			p = int64(x) - int64(generationSums[i-1])
		}
		fmt.Printf("👉 %6d = %6d [%6d]\n", i, x, p)
	}
}

func runNextGeneration(state *state, rules []rule) {

	// fmt.Println("Next generation running...")

	// ensure at 4 empty slots (".") are prefixed
	for strings.Contains(state.value[:4], "#") {
		state.offset++
		state.value = "." + state.value
	}
	// ensure at 4 empty slots (".") are suffixed
	for strings.Contains(state.value[len(state.value)-4:], "#") {
		state.value += "."
	}

	nextValue := state.value
	// ensure at least 2 left and right (above checked this are empty slots always)
	for i := 2; i < len(state.value)-2; i++ {
		sub := state.value[i-2 : i+3]
		matched := false
		for r, rule := range rules {
			if sub == rule.pattern {
				if debug {
					fmt.Printf("Rule %d has matched\n", r)
				}
				nextValue = nextValue[:i] + rule.filled + nextValue[i+1:]
				matched = true
				break
			}
		}
		if !matched {
			nextValue = nextValue[:i] + "." + nextValue[i+1:]
		}
	}
	state.value = nextValue

	// fmt.Println("Next generation completed")

}

func printState(generation int64, state *state) {
	fmt.Printf("%2d [offset=%2d] ", generation, state.offset)
	fmt.Print(state.value)
	fmt.Println()
}

func sumOfFilledItemIndexes(state *state) int64 {
	sum := int64(0)
	for i := 0; i < len(state.value); i++ {
		p := i - state.offset
		if state.value[i:i+1] == "#" {
			sum += int64(p)
		}
	}
	return sum
}
