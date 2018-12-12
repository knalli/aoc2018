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

type cell struct {
	rackId     int
	powerLevel int
}

func findGridFuelCell(serialNo int, size int, squareSize int) (int, int, int) {
	grid := buildFuelGrid(serialNo, size)
	gridSummedArea := buildSummedAreaGrid(grid)

	maxTLCoordinateX := 0
	maxTLCoordinateY := 0
	maxPowerLevelSum := math.MinInt16
	squareSizeOffset := squareSize - 1
	for y := squareSizeOffset; y < size; y++ {
		for x := squareSizeOffset; x < size; x++ {
			lowerRight := gridSummedArea[y][x].powerLevel
			up := 0
			if y-squareSizeOffset > 0 {
				// up
				up = gridSummedArea[y-squareSize][x].powerLevel
			}
			left := 0
			if x-squareSizeOffset > 0 {
				// left
				left = gridSummedArea[y][x-squareSize].powerLevel
			}
			upperLeft := 0
			if y-squareSizeOffset > 0 && x-squareSizeOffset > 0 {
				// upperLeft
				upperLeft = gridSummedArea[y-squareSize][x-squareSize].powerLevel
			}
			sum := lowerRight - up - left + upperLeft
			if sum > maxPowerLevelSum {
				maxPowerLevelSum = sum
				maxTLCoordinateX = x - squareSizeOffset
				maxTLCoordinateY = y - squareSizeOffset
			}
		}
	}

	return maxTLCoordinateX + 1, maxTLCoordinateY + 1, maxPowerLevelSum
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

func buildSummedAreaGrid(grid [][]cell) [][]cell {
	sizeY := len(grid)
	result := make([][]cell, sizeY)
	for y := 0; y < sizeY; y++ {
		sizeX := len(grid[y])
		result[y] = make([]cell, sizeX)
		for x := 0; x < sizeX; x++ {
			powerLevel := grid[y][x].powerLevel
			up := 0
			if y > 0 {
				up = result[y-1][x].powerLevel
			}
			left := 0
			if x > 0 {
				left = result[y][x-1].powerLevel
			}
			upperLeft := 0
			if y > 0 && x > 0 {
				upperLeft = result[y-1][x-1].powerLevel
			}
			sum := powerLevel + up + left - upperLeft
			result[y][x] = cell{rackId: 0, powerLevel: sum}
		}
	}
	return result
}
