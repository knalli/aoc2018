package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

const AocDay = 2
const AocDayName = "day02"
const debug = true
const trace = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, err := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("The checksum is: %d\n", checksum(lines))
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Printf("This letters are common: %s\n", findBestCommonLinePart(lines))
	fmt.Println()
}

func checksum(lines []string) int {
	var counts = countCharactersDuplications(lines)
	return counts[2] * counts[3]
}

func countCharactersDuplications(lines []string) map[int]int {
	var counts = make(map[int]int)
	for _, line := range lines {
		var chars = countCharacters(line)
		for i := 2; i < len(line); i++ {
			if containsValue(chars, i) {
				counts[i]++
				if trace {
					log.Printf("Line '%s' contains IDs exactly %d times", line, i)
				}
			}
		}
	}
	return counts
}

func countCharacters(line string) map[int32]int {
	var chars = make(map[int32]int)
	for _, b := range line {
		// b is decimal (not the char)
		chars[b]++
	}
	return chars
}

func containsValue(data map[int32]int, search int) bool {
	for _, n := range data {
		if n == search {
			return true
		}
	}
	return false
}

func findBestCommonLinePart(lines []string) string {

	var bestLine1 string
	var bestLine2 string
	var bestDiff = math.MaxInt16

	for _, line1 := range lines {
		for _, line2 := range lines {
			if line1 == line2 {
				continue
			}
			diff := diffString(line1, line2)
			if diff < bestDiff {
				bestDiff = diff
				bestLine1 = line1
				bestLine2 = line2
			}
			if trace {
				fmt.Printf("Line '%s' and '%s' differ by: %d\n", line1, line2, diff)
			}
		}
	}

	if debug {
		fmt.Printf("BEST MATCH: Line '%s' and '%s' differ by: %d\n", bestLine1, bestLine2, bestDiff)
	}

	return extractCommonParts(bestLine1, bestLine2)
}

/**
Returns the number of different characters
*/
func diffString(line1, line2 string) int {

	// ensure first line is the smaller one
	if len(line1) > len(line2) {
		return diffString(line2, line1)
	}

	var result int
	for i := range line1 {
		if line1[i] != line2[i] {
			result++
		}
	}
	result += len(line2) - len(line1) // add the rest if line2 is actually longer

	return result
}

func extractCommonParts(line1 string, line2 string) string {

	if len(line1) != len(line2) {
		panic(errors.New("the length of both strings are not equal (not supported)"))
	}

	result := ""
	for i := range line1 {
		if line1[i] == line2[i] {
			result += string(line1[i])
		}
	}

	return result
}
