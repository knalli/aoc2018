package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 5
const AocDayName = "day05"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	polymerUnits := lines[0]
	fmt.Printf("Input puzzle (polymer) contains %d chars\n", len(polymerUnits))
	reactedPolymerUnits := reducePolymerUnits(polymerUnits)
	fmt.Printf("After reacting, %d units are left\n", len(reactedPolymerUnits))
	fmt.Println()

	dayless.PrintStepHeader(2)
	minLength, minLengthChar := findShortestReactedPolymer(polymerUnits)
	fmt.Printf("Length of the shortest polymer: %d (%c)\n", minLength, minLengthChar)
	fmt.Println()
}

func reducePolymerUnits(s string) string {
	return string(reducePolymerUnits2([]rune(s)))
}

func reducePolymerUnits2(runes []rune) []rune {
	for {
		changed := false
		for i := 0; i < len(runes); i++ {
			if i+1 >= len(runes) {
				break // EOL
			}
			r := runes[i]
			s := runes[i+1]

			if r != s && r%32 == s%32 { // mod 32 for matching a=A, b=B, ... (polarity non-matches)
				runes = append(runes[0:i], runes[i+2:]...) // remove reacting part
				changed = true
				break
			}
		}
		if !changed {
			break
		}
	}
	return runes
}

func findShortestReactedPolymer(s string) (int, int32) {
	defer dayless.TimeTrack(time.Now(), "findShortestReactedPolymer")
	var minLength = len(s)
	var minLengthChar = 'a' - 1
	for i := 'a'; i <= 'z'; i++ {
		m := i % 32
		s2 := strings.Map(func(r rune) rune {
			if r%32 == m { // mod 32 for matching a&A, b&B, ...
				return -1
			}
			return r
		}, s)
		length := len(reducePolymerUnits2([]rune(s2)))
		if minLength > length {
			minLength = length
			minLengthChar = i
		}
	}
	return minLength, minLengthChar
}
