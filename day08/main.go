package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const AocDay = 8
const AocDayName = "day08"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	line, _ := dayless.ReadFileToString(AocDayName + "/puzzle.txt")
	sum, value := resolveTreeData(*line)
	fmt.Printf("Sum of all metadata: %d\n", sum)
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Printf("The value of the root node: %d\n", value)
	fmt.Println()
}

func resolveTreeData(line string) (int, int) {
	sum, value, _ := walkTree(splitAsNumbers(line), 0)
	return sum, value
}

func splitAsNumbers(line string) []int {
	parts := strings.Split(line, " ")
	numbers := make([]int, len(parts))
	for i, part := range parts {
		n, _ := strconv.Atoi(part)
		numbers[i] = n
	}
	return numbers
}

func walkTree(numbers []int, start int) (sum int, value int, distance int) {

	p := start
	numChildren := numbers[p]
	p++
	numMetadata := numbers[p]
	p++

	childValues := make([]int, numChildren)
	for i := 0; i < numChildren; i++ {
		childSum, childValue, childDistance := walkTree(numbers, p)
		childValues[i] = childValue // for value (part2)
		sum += childSum             // for global sum (part1)
		p += childDistance
	}

	// collect meta
	for i := 0; i < numMetadata; i++ {
		entry := numbers[p]
		p++
		sum += entry
		if len(childValues) > 0 {
			if entry <= len(childValues) {
				value += childValues[entry-1]
			} // or 0
		} else {
			value += entry
		}
	}

	distance = p - start

	return sum, value, distance
}
