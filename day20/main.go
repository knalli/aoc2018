package main

import (
	d15Grid "de.knallisworld/aoc/aoc2018/day15/grid"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"time"
)

const AocDay = 20
const AocDayName = "day20"

const (
	unknown rune = '❓'
	room    rune = '🕳'
	wall    rune = '🧱'
	doorH   rune = '🚪'
	doorV   rune = '🚪'
)

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	line, _ := dayless.ReadFileToString(AocDayName + "/puzzle.txt")
	dayless.PrintStepHeader(1)
	compute(*line, 1000, 1000)
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Println("See Part1")
	fmt.Println()
}

func compute(regex string, dimension int, countMinDoors int) d15Grid.Grid {
	fmt.Printf("👉 Regex: %s\n", regex)
	grid := buildGrid(dimension)
	originStart := d15Grid.NewPoint(dimension/2, dimension/2)
	visit(grid, originStart, regex, 0)
	grid, points := compactGrid(grid, []d15Grid.Point{originStart})
	start := points[0]
	grid.WalkAndModify(func(p d15Grid.Point, value rune, cb func(v rune)) {
		if value == unknown {
			cb(wall)
		}
	})
	fmt.Println()
	fmt.Print(grid.ToString())
	fmt.Println()
	fmt.Printf("👉 Start point = %s, grid dimension (height=%d, width=%d)\n", start.ToString(), grid.Height(), grid.Width())
	fmt.Println()

	furthestPoint, distance, totalRoomsMinDistance := findFurthestPointWithShortestPath(grid, start, countMinDoors*2)
	fmt.Printf("🎉 %s (∂=%d, doors=%d)\n", furthestPoint.ToString(), distance, distance/2)
	fmt.Printf("🎉 At least %d rooms pass %d doors\n", totalRoomsMinDistance, countMinDoors)

	return grid
}

func findFurthestPointWithShortestPath(grid d15Grid.Grid, start d15Grid.Point, countRoomsMinDistance int) (point d15Grid.Point, distance int, totalRoomsMinDistance int) {
	distance = 0
	totalRoomsMinDistance = 0
	grid.Walk(func(p d15Grid.Point, value rune) {
		if value != room && p != start {
			return
		}
		path := grid.GetShortestPath2(
			start,
			p,
			func(v rune) bool {
				return v != wall
			},
			nil,
		)
		d := len(path)
		if d > 0 {
			if d > distance {
				distance = d
				point = p
			}
			if d >= countRoomsMinDistance {
				totalRoomsMinDistance++
			}
		}
	})
	return point, distance, totalRoomsMinDistance
}

func buildGrid(dimension int) d15Grid.Grid {
	return d15Grid.NewByDimension(dimension, dimension, unknown)
}

func visit(grid d15Grid.Grid, start d15Grid.Point, regex string, position int) (d15Grid.Point, int) {

	point := start

	for {
		grid.PutValueAtXY(point, room)
		switch regex[position] {
		case '^':
			// start
			position++
		case '$':
			// end-of-line
			position++
			return point, position
		case 'N':
			grid.PutValueAtXY(point.Top(), doorH)
			point = point.Top().Top()
			position++
		case 'E':
			grid.PutValueAtXY(point.Right(), doorV)
			point = point.Right().Right()
			position++
		case 'S':
			grid.PutValueAtXY(point.Bottom(), doorH)
			point = point.Bottom().Bottom()
			position++
		case 'W':
			grid.PutValueAtXY(point.Left(), doorV)
			point = point.Left().Left()
			position++
		case '(':
			position++
			point, position = visit(grid, point, regex, position)
		case ')':
			position++
			// end-of-sub
			return point, position
		case '|':
			point = start
			position++
		default:
			panic(fmt.Sprintf("unknown char in regex %c", regex[position]))
		}
	}
}

func compactGrid(grid d15Grid.Grid, points []d15Grid.Point) (d15Grid.Grid, []d15Grid.Point) {
	min, max := findGridMinMax(grid, true)
	dx := max.X - min.X + 1
	dy := max.Y - min.Y + 1
	resultGrid := d15Grid.NewByDimension(dy, dx, unknown)
	for y := min.Y; y <= max.Y; y++ {
		for x := min.X; x <= max.X; x++ {
			resultGrid.PutValueAtXY(d15Grid.Point{X: x - min.X, Y: y - min.Y}, grid.GetValueAtXY(d15Grid.Point{X: x, Y: y}))
		}
	}
	// transform all given points also
	resultPoints := make([]d15Grid.Point, len(points))
	for i, p := range points {
		resultPoints[i] = d15Grid.Point{X: p.X - min.X, Y: p.Y - min.Y}
	}
	return resultGrid, resultPoints
}

func findGridMinMax(grid d15Grid.Grid, includeOuter bool) (d15Grid.Point, d15Grid.Point) {

	min := d15Grid.Point{X: 0, Y: 0}
	max := d15Grid.Point{X: grid.Width() - 1, Y: grid.Height() - 1}
	center := d15Grid.Point{X: grid.Width() / 2, Y: grid.Height() / 2}

	outerLineOffset := 0
	if includeOuter {
		outerLineOffset = 1
	}

	// to-top: find first line of all-unknowns
	for y := center.Y; y >= 0; y-- {
		allUnknown := true
		for x := 0; x < grid.Width(); x++ {
			if grid.GetValueAtXY(d15Grid.Point{X: x, Y: y}) != unknown {
				allUnknown = false
				break
			}
		}
		if allUnknown {
			min.Y = y + 1 - outerLineOffset
			break
		}
	}
	// to-bottom: find first line of all-unknowns
	for y := center.Y; y < grid.Height(); y++ {
		allUnknown := true
		for x := 0; x < grid.Width(); x++ {
			if grid.GetValueAtXY(d15Grid.Point{X: x, Y: y}) != unknown {
				allUnknown = false
				break
			}
		}
		if allUnknown {
			max.Y = y - 1 + outerLineOffset
			break
		}
	}
	// to-left: find first line of all-unknowns
	for x := center.X; x >= 0; x-- {
		allUnknown := true
		for y := 0; y < grid.Height(); y++ {
			if grid.GetValueAtXY(d15Grid.Point{X: x, Y: y}) != unknown {
				allUnknown = false
				break
			}
		}
		if allUnknown {
			min.X = x + 1 - outerLineOffset
			break
		}
	}
	// to-right: find first line of all-unknowns
	for x := center.X; x < grid.Width(); x++ {
		allUnknown := true
		for y := 0; y < grid.Height(); y++ {
			if grid.GetValueAtXY(d15Grid.Point{X: x, Y: y}) != unknown {
				allUnknown = false
				break
			}
		}
		if allUnknown {
			max.X = x - 1 + outerLineOffset
			break
		}
	}

	return min, max
}
