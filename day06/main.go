package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const AocDay = 6
const AocDayName = "day06"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	grid, coordinates, _ := createGrid(lines)
	fillGridManhattenDistance(grid, coordinates)
	// fmt.Println(renderGrid(grid))
	maxAreaSize, _ := getLargestFiniteArea(grid)
	fmt.Printf("Size of the largest area: %d\n", maxAreaSize)
	fmt.Println()

	dayless.PrintStepHeader(2)
	grid, _, _ = createGrid(lines)
	fillGridManhattenDistanceSum(grid, coordinates)
	sum := 0
	for _, row := range grid.rows {
		for _, col := range row {
			if col == 1 {
				sum++
			}
		}
	}
	fmt.Printf("Size of region: %d\n", sum)
	// fmt.Println(renderGrid(grid))
	fmt.Println()
}

type coordinate struct {
	x  int
	y  int
	id int
}

type grid struct {
	rows map[int]map[int]int
}

func createGrid(lines []string) (grid, []coordinate, error) {
	var coordinates []coordinate
	maxX := 0
	maxY := 0
	for i, line := range lines {
		parts := strings.Split(line, ", ")
		if len(parts) != 2 {
			return grid{}, nil, errors.New("invalid size of numbers in line")
		}
		x, err := strconv.Atoi(parts[0])
		if err != nil {
			return grid{}, nil, errors.New("invalid number")
		}
		y, err := strconv.Atoi(parts[1])
		if err != nil {
			return grid{}, nil, errors.New("invalid number")
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		coordinates = append(coordinates, coordinate{x: x, y: y, id: i})
	}
	maxX++
	maxY++
	var rows = make(map[int]map[int]int, maxX)
	for x := 0; x < maxX; x++ {
		var row = make(map[int]int, maxY)
		for y := 0; y < maxY; y++ {
			row[y] = -1

		}
		rows[x] = row
	}
	return grid{rows: rows}, coordinates, nil
}

func fillGridManhattenDistance(grid grid, coordinates []coordinate) {
	for x := 0; x < len(grid.rows); x++ {
		row := grid.rows[x]
		for y := 0; y < len(row); y++ {
			closestCoordinate, _, err := calcClosestManhattenDistance(coordinates, coordinate{x: x, y: y, id: -1})
			if err != nil {
				row[y] = -1
			} else {
				row[y] = closestCoordinate.id
			}
		}
	}
}

func fillGridManhattenDistanceSum(grid grid, coordinates []coordinate) {
	for x := 0; x < len(grid.rows); x++ {
		row := grid.rows[x]
		for y := 0; y < len(row); y++ {
			sum := 0
			for _, c := range coordinates {
				sum += calcManhattenDistance(c, coordinate{x: x, y: y, id: -1})
			}
			if sum < 10000 {
				row[y] = 1
			} else {
				row[y] = -1
			}
		}
	}
}

func calcManhattenDistance(from coordinate, to coordinate) int {
	return int(math.Abs(float64(from.x-to.x))) + int(math.Abs(float64(from.y-to.y)))
}

func calcClosestManhattenDistance(froms []coordinate, to coordinate) (coordinate, int, error) {
	min := math.MaxInt16
	minC := coordinate{}
	for _, from := range froms {
		distance := calcManhattenDistance(from, to)
		if distance < min {
			min = distance
			minC = from
		} else if distance == min {
			return coordinate{}, 0, errors.New("same distance")
		}
	}
	return minC, min, nil
}

func getLargestFiniteArea(grid grid) (int, int) {
	// skip edges (0, max)
	infiniteAreaIds := make(map[int]struct{})
	for _, id := range grid.rows[0] {
		if _, ok := infiniteAreaIds[id]; !ok {
			infiniteAreaIds[id] = struct{}{}
		}
	}
	for _, id := range grid.rows[len(grid.rows)-1] {
		if _, ok := infiniteAreaIds[id]; !ok {
			infiniteAreaIds[id] = struct{}{}
		}
	}
	for _, row := range grid.rows {
		id := row[0]
		if _, ok := infiniteAreaIds[id]; !ok {
			infiniteAreaIds[id] = struct{}{}
		}
		id = row[len(row)-1]
		if _, ok := infiniteAreaIds[id]; !ok {
			infiniteAreaIds[id] = struct{}{}
		}
	}

	var areaSizes = make(map[int]int)
	for x := 0; x < len(grid.rows)-1; x++ {
		row := grid.rows[x]
		for y := 0; y < len(row)-1; y++ {
			id := row[y]
			if _, ok := infiniteAreaIds[id]; ok {
				continue // infinite area
			}
			areaSizes[id] += 1
		}
	}

	maxSize := 0
	maxId := 0
	for id, size := range areaSizes {
		if size > maxSize {
			maxSize = size
			maxId = id
		}
	}

	return maxSize, maxId
}

func renderGrid(grid grid) string {
	var result = ""
	for x := 0; x < len(grid.rows); x++ {
		row := grid.rows[x]
		result += fmt.Sprintf("%03d | ", x)
		for y := 0; y < len(row); y++ {
			col := fmt.Sprintf("%3d", row[y])
			result += col
		}
		result += "\n"
	}
	return result
}
