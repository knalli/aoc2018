package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"math"
	"strconv"
	"time"
)

const AocDay = 14
const AocDayName = "day14"
const debug = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	thresholdMinStr, _ := dayless.ReadFileToString(AocDayName + "/puzzle.txt")
	thresholdMin, _ := strconv.Atoi(*thresholdMinStr)
	fmt.Printf("🎉 The scores of the ren recipes after the puzzle input: %d\n", cook1(thresholdMin, 10, 2))
	fmt.Println()

	dayless.PrintStepHeader(2)
	for _, str := range []string{"51589", "01245", "92510", "59414", *thresholdMinStr} {
		fmt.Printf("🎉 %s first appears after %d recipes\n", str, cook2(str, 2))
	}
	fmt.Println()
}

func cook1(needle int, topX int, workerSize int) int {

	recipes := make([]int, 1000000000)
	offset := 0

	recipes[offset] = 3
	offset++
	recipes[offset] = 7
	offset++
	length := offset

	workers := make([]int, workerSize)
	for i := 0; i < workerSize; i++ {
		workers[i] = i
	}

	fmt.Println("👉 Let's start cooking...")
	for {
		// create score
		score := 0
		for i := 0; i < workerSize; i++ {
			workerScore := recipes[workers[i]]
			score += workerScore
		}
		// split digits & create new recipes
		for _, c := range fmt.Sprintf("%d", score) {
			d := int(c) - 48 // ASCII
			recipes[offset] = d
			offset++
		}
		length = offset

		for i := 0; i < workerSize; i++ {
			workerScore := recipes[workers[i]]
			workers[i] = (workers[i] + 1 + workerScore) % length
		}

		if debug {
			fmt.Printf("%3d ", length)
			fmt.Print(workers)
			fmt.Print(" ")
			fmt.Println(recipes)
		}

		if length >= needle+topX {
			break
		}
	}

	return convertArrayIntoInt(recipes[needle : needle+topX])
}

func cook2(needle string, workerSize int) int {

	recipeDigits := make([]int, 1000000000)
	recipeRunes := make([]rune, len(recipeDigits))
	offset := 0

	recipeDigits[offset] = 3
	offset++
	recipeDigits[offset] = 7
	offset++
	length := offset

	workers := make([]int, workerSize)
	for i := 0; i < workerSize; i++ {
		workers[i] = i
	}

	fmt.Println("👉 Let's start cooking...")

	found := false
	for !found {
		// create score
		score := 0
		for i := 0; i < workerSize; i++ {
			workerScore := recipeDigits[workers[i]]
			score += workerScore
		}
		// split digits & create new recipes
		for _, c := range fmt.Sprintf("%d", score) {
			d := int(c) - 48 // ASCII
			recipeDigits[offset] = d
			recipeRunes[offset] = c
			offset++
			length = offset

			// Compare if the needle is at the end of the processed numbers
			if length >= len(needle) && string(recipeRunes[length-len(needle):length]) == needle {
				found = true
				break
			}
		}

		for i := 0; i < workerSize; i++ {
			workerScore := recipeDigits[workers[i]]
			workers[i] = (workers[i] + 1 + workerScore) % length
		}

		if debug {
			fmt.Printf("%3d ", length)
			fmt.Print(workers)
			fmt.Print(" ")
			fmt.Println(recipeDigits)
		}
	}

	return length - len(needle)
}

func getPow(n int) int {
	result := 0
	x := float64(n)
	for x > 1 {
		x /= 10
		result++
	}
	return result
}

func convertArrayIntoInt(numbers []int) int {
	result := 0
	for c, v := range numbers {
		result += int(math.Pow10(len(numbers)-c-1)) * v
	}
	return result
}
