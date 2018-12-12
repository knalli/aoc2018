package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"math"
	"strconv"
	"time"
)

const AocDay = 11
const AocDayName = "day11_alternative"

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
	fmt.Printf("🎉 Final answer: %d,%d,%d\n", topLeftX, topLeftRight, maxSquareSize)
	fmt.Println()
}

func findGridFuelCell(serialNo int, gridSize int, squareSize int) (int, int, int) {
	grid := buildFuelGrid(serialNo, gridSize)
	gridSummedArea := buildSummedAreaGrid(grid)

	maxPowerLevelX := 0
	maxPowerLevelY := 0
	maxPowerLevelSum := math.MinInt16
	squareSizeOffset := squareSize - 1
	for y := squareSizeOffset; y < gridSize; y++ {
		for x := squareSizeOffset; x < gridSize; x++ {
			lowerRight := gridSummedArea[y][x]
			up := 0
			if y-squareSizeOffset > 0 {
				// up
				up = gridSummedArea[y-squareSize][x]
			}
			left := 0
			if x-squareSizeOffset > 0 {
				// left
				left = gridSummedArea[y][x-squareSize]
			}
			upperLeft := 0
			if y-squareSizeOffset > 0 && x-squareSizeOffset > 0 {
				// upperLeft
				upperLeft = gridSummedArea[y-squareSize][x-squareSize]
			}
			sum := lowerRight - up - left + upperLeft
			if sum > maxPowerLevelSum {
				maxPowerLevelSum = sum
				maxPowerLevelX = x - squareSizeOffset
				maxPowerLevelY = y - squareSizeOffset
			}
		}
	}

	return maxPowerLevelX + 1, maxPowerLevelY + 1, maxPowerLevelSum
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

func buildSummedAreaGrid(grid [][]int) [][]int {
	sizeY := len(grid)
	result := make([][]int, sizeY)
	for y := 0; y < sizeY; y++ {
		sizeX := len(grid[y])
		result[y] = make([]int, sizeX)
		for x := 0; x < sizeX; x++ {
			powerLevel := grid[y][x]
			up := 0
			if y > 0 {
				up = result[y-1][x]
			}
			left := 0
			if x > 0 {
				left = result[y][x-1]
			}
			upperLeft := 0
			if y > 0 && x > 0 {
				upperLeft = result[y-1][x-1]
			}
			sum := powerLevel + up + left - upperLeft
			result[y][x] = sum
		}
	}
	return result
}
