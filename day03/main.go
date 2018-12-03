package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const AocDay = 3
const AocDayName = "day03"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, err := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	if err != nil {
		panic(err)
	}
	claims, err := transformClaims(lines)
	if err != nil {
		panic(err)
	}
	fabric := getMaxDimensions(claims)
	fmt.Printf("The input puzzle has a maximum dimension of %dx%d.\n", fabric.width, fabric.height)
	fmt.Printf("There are %d inches wich are within two or more claims.\n", numberOfOverlappingBoxes(fabric, claims))
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Printf("These ids are not overlapping: %v", resolveNonOverlappingIds(fabric, claims))
	fmt.Println()
}

type claim struct {
	id     int
	x      int
	y      int
	width  int
	height int
}

type grid struct {
	width  int
	height int
}

func transformClaims(lines []string) ([]claim, error) {
	defer dayless.TimeTrack(time.Now(), "transformClaims")
	var result []claim
	for _, line := range lines {
		parts := strings.Split(line, " ")
		id, err := strconv.Atoi(parts[0][1:])
		if err != nil {
			return nil, err
		}
		coordinates := strings.Split(parts[2][0:len(parts[2])-1], ",")
		dimensions := strings.Split(parts[3], "x")
		x, err := strconv.Atoi(coordinates[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(coordinates[1])
		if err != nil {
			return nil, err
		}
		w, err := strconv.Atoi(dimensions[0])
		if err != nil {
			return nil, err
		}
		h, err := strconv.Atoi(dimensions[1])
		if err != nil {
			return nil, err
		}
		result = append(result, claim{
			id:     id,
			x:      x,
			y:      y,
			width:  w,
			height: h,
		})
	}
	return result, nil
}

func getMaxDimensions(claims []claim) grid {
	var width = 0
	var height = 0

	for _, claim := range claims {
		if width < claim.x+claim.width {
			width = claim.x + claim.width
		}
		if height < claim.y+claim.height {
			height = claim.y + claim.height
		}
	}

	return grid{
		width:  width + 1,
		height: height + 1,
	}
}

func numberOfOverlappingBoxes(fabric grid, claims []claim) int {

	result := 0
	m := buildMatrix(fabric, claims)

	for x := 0; x < len(m); x++ {
		for y := 0; y < len(m[x]); y++ {
			if len(m[x][y]) >= 2 {
				result++
			}
		}
	}

	return result
}

func resolveNonOverlappingIds(fabric grid, claims []claim) []int {

	result := make([]int, 0)
	m := buildMatrix(fabric, claims)

	for _, claim := range claims {
		found := true
		for i := 0; i < claim.width; i++ {
			for j := 0; j < claim.height; j++ {
				x := claim.x + i
				y := claim.y + j
				if len(m[x][y]) != 1 {
					found = false
					break
				}
			}
			if !found {
				break
			}
		}

		if found {
			result = appendIfMissing(result, claim.id)
		}
	}

	return result
}

func appendIfMissing(slice []int, i int) []int {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func buildMatrix(grid grid, claims []claim) [][][]int {
	defer dayless.TimeTrack(time.Now(), "buildMatrix")
	m := make([][][]int, grid.width)
	for i := range m {
		m[i] = make([][]int, grid.height)
		for j := range m[i] {
			m[i][j] = make([]int, 0)
		}
	}
	for _, claim := range claims {
		for x := claim.x; x < claim.x+claim.width; x++ {
			for y := claim.y; y < claim.y+claim.height; y++ {
				m[x][y] = append(m[x][y], claim.id)
			}
		}
	}
	return m
}
