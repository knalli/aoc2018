package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const AocDay = 10
const AocDayName = "day10"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	lookIntoTheStars(lines)
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Println("See above")
	fmt.Println()
}

type point struct {
	positionX int
	positionY int
	velocityX int
	velocityY int
}

func lookIntoTheStars(lines []string) {
	points := parseLine(lines)

	// look for smallest room
	minI := 0
	minDx := math.MaxInt16
	minDy := math.MaxInt16
	maxWait := 40000
	for i := 1; i < maxWait; i++ {
		movePoints(points, 1)
		dx, dy := measurePoints(points)
		if minDx > dx && minDy > dy {
			minI = i
			minDx = dx
			minDy = dy
		}
	}

	fmt.Printf("Minimum @ #%ds (%d,%d)\n", minI, minDx, minDy)
	fmt.Println("The sky looks like...")
	movePoints(points, -(maxWait - 1 - minI)) // move back
	printPoints(points, true)
}

func parseLine(lines []string) []point {
	// position=<-43063,  10936> velocity=< 4, -1>
	pattern := regexp.MustCompile("position=<\\s*(-?\\d+),\\s*(-?\\d+)> velocity=<\\s*(-?\\d+),\\s*(-?\\d+)>")

	result := make([]point, len(lines))
	for i, line := range lines {
		match := pattern.FindStringSubmatch(line)
		positionX, _ := strconv.Atoi(match[1])
		positionY, _ := strconv.Atoi(match[2])
		velocityX, _ := strconv.Atoi(match[3])
		velocityY, _ := strconv.Atoi(match[4])
		result[i] = point{
			positionX: positionX,
			positionY: positionY,
			velocityX: velocityX,
			velocityY: velocityY,
		}
	}

	return result
}

func movePoints(points []point, seconds int) {
	for i := range points {
		point := &points[i]
		point.positionX += point.velocityX * seconds
		point.positionY += point.velocityY * seconds
	}
}

func measurePoints(points []point) (int, int) {
	minX := ^int(0)
	maxX := -minX - 1
	minY := ^int(0)
	maxY := -minY - 1
	for _, point := range points {
		if minX > point.positionX {
			minX = point.positionX
		} else if maxX < point.positionX {
			maxX = point.positionX
		}
		if minY > point.positionY {
			minY = point.positionY
		} else if maxY < point.positionY {
			maxY = point.positionY
		}
	}
	dx := int(math.Abs(float64(minX)) + math.Abs(float64(maxX)))
	dy := int(math.Abs(float64(minY)) + math.Abs(float64(maxY)))
	return dx, dy
}

func printPoints(points []point, skipEmptyLines bool) {
	minX := ^int(0)
	maxX := -minX - 1
	minY := ^int(0)
	maxY := -minY - 1
	for _, point := range points {
		if minX > point.positionX {
			minX = point.positionX
		} else if maxX < point.positionX {
			maxX = point.positionX
		}
		if minY > point.positionY {
			minY = point.positionY
		} else if maxY < point.positionY {
			maxY = point.positionY
		}
	}
	dx := int(math.Abs(float64(minX)) + math.Abs(float64(maxX)))
	dy := int(math.Abs(float64(minY)) + math.Abs(float64(maxY)))
	grid := make([][]int, dy+1)
	for i := 0; i <= dy; i++ {
		grid[i] = make([]int, dx+1)
	}
	for _, point := range points {
		ax := point.positionX - minX
		ay := point.positionY - minY
		grid[ay][ax] = 1
	}

	s := ""
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] == 1 {
				s += "#"
			} else {
				s += "."
			}
		}
		if !skipEmptyLines || strings.Contains(s, "#") {
			fmt.Println(s)
		}
		s = ""
	}
}
