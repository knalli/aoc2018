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

func findGridFuelCell(serialNo int, gridSize int, squareSize int) (int, int, int) {
	grid := buildFuelGrid(serialNo, gridSize)

	maxPowerLevelX := 0
	maxPowerLevelY := 0
	maxPowerLevelSum := math.MinInt16
	for y := 1; y <= gridSize-squareSize; y++ {
		for x := 1; x <= gridSize-squareSize; x++ {
			// compute for square
			// fmt.Printf("Checking %d,%d -> %d,%d\n", x, y, x+squareSize-1, y+squareSize-1)
			powerLevelSum := 0
			for i := 0; i < squareSize; i++ {
				for j := 0; j < squareSize; j++ {
					powerLevelSum += grid[y-1+i][x-1+j]
				}
			}
			if powerLevelSum > maxPowerLevelSum {
				maxPowerLevelSum = powerLevelSum
				maxPowerLevelX = x
				maxPowerLevelY = y
			}
		}
	}

	return maxPowerLevelX, maxPowerLevelY, maxPowerLevelSum
}

func findMaxFuelGridSize(serialNo int, gridSize int) (int, int, int, int) {

	maxPowerLevelX := 0
	maxPowerLevelY := 0
	maxPowerLevelSum := math.MinInt16
	maxSquareSize := 0

	for squareSize := 1; squareSize <= gridSize; squareSize++ {
		x, y, powerLevelSum := findGridFuelCell(serialNo, gridSize, squareSize)

		if powerLevelSum > maxPowerLevelSum {
			maxPowerLevelSum = powerLevelSum
			maxPowerLevelX = x
			maxPowerLevelY = y
			maxSquareSize = squareSize

			fmt.Printf("⏳ Found %d,%d power=%d, size=%d …\n", x, y, powerLevelSum, squareSize)
		}
	}

	return maxPowerLevelX, maxPowerLevelY, maxPowerLevelSum, maxSquareSize
}

func buildFuelGrid(serialNo int, size int) [][]int {
	grid := make([][]int, size)
	for y := 1; y <= size; y++ {
		grid[y-1] = make([]int, size)
		for x := 1; x <= size; x++ {
			rackId := x + 10
			powerLevel := rackId * y
			powerLevel += serialNo
			powerLevel *= rackId
			powerLevel = (powerLevel % 1000) / 100 // "The hundreds digit"
			powerLevel -= 5
			grid[y-1][x-1] = powerLevel
		}
	}
	return grid
}
