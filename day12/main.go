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
	result := runNextGenerations(20, state, rules)
	fmt.Printf("🎉 The sum of the numbers of all pots which contain a plant: %d\n", result)
	fmt.Println()

	dayless.PrintStepHeader(2)
	// TODO after generation ~166, this are adding 73 for each constantly, running for 50000000000 is too slow
	result = runNextGenerations(50000000000-20, state, rules)
	fmt.Printf("🎉 The sum of the numbers of all pots which contain a plant: %d\n", result)
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

func runNextGenerations(generations int64, state *state, rules []rule) int64 {
	defer dayless.TimeTrack(time.Now(), fmt.Sprintf("runNextGenerations=%d", generations))

	lastSum := int64(0)
	lastDiff := int64(0)
	diffStrikes := 0

	g := int64(0)
	strikeLengthRequired := int(float64(generations) * 0.00000001)
	if strikeLengthRequired < 100 {
		strikeLengthRequired = int(generations)
	}
	fmt.Printf("👉 Require a strike-length of at least %d\n", strikeLengthRequired)
	for ; g < generations && diffStrikes < strikeLengthRequired; g++ {
		runNextGeneration(state, rules)
		sum := sumOfFilledItemIndexes(state)
		diff := sum - lastSum
		if debug {
			fmt.Printf("👉 %6d = %6d [%6d]\n", g, sum, diff)
		}
		lastSum = sum
		if lastDiff == diff {
			diffStrikes++
		} else {
			diffStrikes = 0 // reset
			lastDiff = diff
		}
	}

	if g < generations {
		fmt.Printf("👉 After %d generations, a %d-strike of diff %d has been detected.\n", g, diffStrikes, lastDiff)
		fmt.Printf("👉 Assuming this value is constant for the %d generations left.\n", generations-g)
	}

	return lastSum + ((generations - g) * lastDiff)
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
