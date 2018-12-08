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
	fmt.Printf("Sum of all metadata: %d", getSumOfAllMetadata(*line))
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Println("N/A")
	fmt.Println()
}

func getSumOfAllMetadata(line string) int {
	sum, _ := infixWalk2(splitAsNumbers(line), 0)
	return sum
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

func infixWalk2(numbers []int, start int) (metadataSum int, distance int) {

	p := start
	quantityChildren := numbers[p]
	p++
	quantityMetadata := numbers[p]
	p++

	for i := 0; i < quantityChildren; i++ {
		childMetadataSum, childDistance := infixWalk2(numbers, p)
		metadataSum += childMetadataSum
		p += childDistance
	}

	// collect meta
	for i := 0; i < quantityMetadata; i++ {
		metadataSum += numbers[p]
		p++
	}

	distance = p - start

	return metadataSum, distance
}
