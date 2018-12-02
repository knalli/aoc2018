package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"log"
	"time"
)

const AocDay = 2
const AocDayName = "day02"
const debug = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, err := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("The checksum is: %d", checksum(lines))
	fmt.Println()

	dayless.PrintStepHeader(2)
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
				if debug {
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
