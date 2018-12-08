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
	sum, value := getSumOfAllMetadataAndValues(*line)
	fmt.Printf("Sum of all metadata: %d\n", sum)
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Printf("The value of the root node: %d\n", value)
	fmt.Println()
}

func getSumOfAllMetadataAndValues(line string) (int, int) {
	sum, value, _ := infixWalk2(splitAsNumbers(line), 0)
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

func infixWalk2(numbers []int, start int) (metadataSum int, value int, distance int) {

	p := start
	quantityChildren := numbers[p]
	p++
	quantityMetadata := numbers[p]
	p++

	childValues := make([]int, quantityChildren)
	for i := 0; i < quantityChildren; i++ {
		childMetadataSum, childValue, childDistance := infixWalk2(numbers, p)
		childValues[i] = childValue     // for value (part2)
		metadataSum += childMetadataSum // for global sum (part1)
		p += childDistance
	}

	// collect meta
	for i := 0; i < quantityMetadata; i++ {
		entry := numbers[p]
		p++
		metadataSum += entry
		if len(childValues) > 0 {
			if entry <= len(childValues) {
				value += childValues[entry-1]
			} // or 0
		} else {
			value += entry
		}
	}

	distance = p - start

	return metadataSum, value, distance
}
