package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const AocDay = 1
const AocDayName = "day01"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	fmt.Println("--- Part One ---")
	lines, err := dayless.ReadFileToArray(AocDayName + "/puzzle1.txt")
	if err != nil {
		panic(err)
		return
	}
	numbers, err := transformStringsToNumbers(lines)
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("The resulting frequency is: %d\n", sumNumbers(numbers))
	fmt.Println()

	fmt.Println("--- Part Two ---")
	interval, index, result := findFirstDuplicateResult(numbers, 1)
	fmt.Printf("The first frequency reached is: %d (interval=%d, index=%d)\n", result, interval, index)
	fmt.Println()
}

func transformStringsToNumbers(lines []string) ([]int, error) {
	var numbers []int
	for _, line := range lines {
		i, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		if line[0] == '+' {
			numbers = append(numbers, i)
		} else if line[0] == '-' {
			numbers = append(numbers, -i)
		} else {
			return nil, errors.New("invalid sign found")
		}
	}
	return numbers, nil
}

func sumNumbers(numbers []int) int {
	result := 0
	for _, i := range numbers {
		result += i
	}
	return result
}

func findFirstDuplicateResult(numbers []int, min int) (interval int, index int, result int) {

	var counts = make(map[int]int)

	result = 0
	counts[result] = 1

	for {
		interval++
		for index, i := range numbers {
			result += i
			counts[result] = counts[result] + 1
			if counts[result] > min {
				return interval, index, result
			}
		}
	}
}
