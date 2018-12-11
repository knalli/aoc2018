package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"math"
	"strconv"
	"time"
)

const AocDay = 11
const AocDayName = "day11"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	line, _ := dayless.ReadFileToString(AocDayName + "/puzzle.txt")
	serialNo, _ := strconv.Atoi(*line)
	topLeftX, topLeftRight, powerSum := findGridFuelCell(serialNo, 300, 3)
	fmt.Printf("👉 The largtest power square is @ %d,%d with the power sum of %d\n", topLeftX, topLeftRight, powerSum)
	fmt.Printf("🎉 Final answer: %d,%d\n", topLeftX, topLeftRight)
	fmt.Println()

	dayless.PrintStepHeader(2)
	topLeftX, topLeftRight, powerSum, maxSquareSize := findMaxFuelGridSize(serialNo, 300)
	fmt.Printf("👉 The largtest total power square is @ %d,%d with the power sum of %d and a max square of %d\n", topLeftX, topLeftRight, powerSum, maxSquareSize)
	fmt.Printf("🎉 Final answer: %d,%d,%d\n", topLeftX, topLeftRight, powerSum)
	fmt.Println()
}

type cell struct {
	rackId     int
	powerLevel int
}

func findGridFuelCell(serialNo int, size int, squareSize int) (int, int, int) {
	grid := buildFuelGrid(serialNo, size)

	maxTLCoordinateX := 0
	maxTLCoordinateY := 0
	maxPowerLevelSum := math.MinInt16
	for y := 1; y <= size-squareSize; y++ {
		for x := 1; x <= size-squareSize; x++ {
			// compute for square
			// fmt.Printf("Checking %d,%d -> %d,%d\n", x, y, x+squareSize-1, y+squareSize-1)
			sum := 0
			for i := 0; i < squareSize; i++ {
				for j := 0; j < squareSize; j++ {
					sum += grid[y-1+i][x-1+j].powerLevel
				}
			}
			if sum > maxPowerLevelSum {
				maxPowerLevelSum = sum
				maxTLCoordinateX = x
				maxTLCoordinateY = y
			}
		}
	}

	return maxTLCoordinateX, maxTLCoordinateY, maxPowerLevelSum
}

func findMaxFuelGridSize(serialNo int, size int) (int, int, int, int) {

	maxTLCoordinateX := 0
	maxTLCoordinateY := 0
	maxPowerLevelSum := math.MinInt16
	maxSquareSize := 0

	for n := 1; n <= size; n++ {
		x, y, powerLevelSum := findGridFuelCell(serialNo, size, n)

		if powerLevelSum > maxPowerLevelSum {
			maxPowerLevelSum = powerLevelSum
			maxTLCoordinateX = x
			maxTLCoordinateY = y
			maxSquareSize = n

			fmt.Printf("⏳ Found %d,%d power=%d, size=%d …\n", x, y, powerLevelSum, n)
		}
	}

	return maxTLCoordinateX, maxTLCoordinateY, maxPowerLevelSum, maxSquareSize
}

func buildFuelGrid(serialNo int, size int) [][]cell {
	grid := make([][]cell, size)
	for y := 1; y <= size; y++ {
		grid[y-1] = make([]cell, size)
		for x := 1; x <= size; x++ {
			rackId := x + 10
			temp := rackId * y
			temp += serialNo
			temp *= rackId
			temp = (temp % 1000) / 100 // "The hundreds digit"
			temp -= 5
			grid[y-1][x-1] = cell{rackId: rackId, powerLevel: temp}
		}
	}
	return grid
}
