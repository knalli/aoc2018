package main

import (
	"de.knallisworld/aoc/aoc2018/day15/grid"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"math"
	"regexp"
	"time"
)

const AocDay = 17
const AocDayName = "day17"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	clayPoints, topLeftPoint, bottomRightPoint := parsePoints(lines)
	tiles := buildTiles(clayPoints)
	drainWater2(tiles, 500, 0, topLeftPoint, bottomRightPoint)
	fmt.Println(tiles.ToStringSkip(topLeftPoint, bottomRightPoint))
	fmt.Printf("Tiles with water: %d\n", countGroundType(tiles, topLeftPoint, bottomRightPoint, settledWater)+countGroundType(tiles, topLeftPoint, bottomRightPoint, dryWater))
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Printf("Tiles with water (retained): %d\n", countGroundType(tiles, topLeftPoint, bottomRightPoint, settledWater))
	fmt.Println()
}

type cellType rune

const (
	sand         cellType = '.'
	clay         cellType = '#'
	spring       cellType = '+'
	settledWater cellType = '~'
	dryWater     cellType = '|'
)

func (c cellType) Rune() rune {
	switch c {
	case clay:
		return '🧱'
		// return '#'
	case sand:
		return '🕳'
		// return '.'
	case dryWater:
		return '💧'
		// return '|'
	case settledWater:
		return '🌊'
		// return '~'
	case spring:
		return '💦'
		// return '+'
	default:
		panic("invalid cell type")
	}
}

type ground [][]cellType

func (g ground) ToString() string {
	result := ""
	maxY := len(g)
	for y := 0; y < maxY; y++ {
		maxX := len(g[y])
		runes := make([]rune, maxX)
		for x := 0; x < maxX; x++ {
			runes[x] = g[y][x].Rune()
		}
		result += string(runes)
		result += "\n"
	}
	return result
}

func (g ground) ToStringSkip(topLeftPoint grid.Point, bottomRightPoint grid.Point) string {
	result := ""
	for y, row := range g {
		if topLeftPoint.Y <= y && y <= bottomRightPoint.Y {
			selected := row
			runes := make([]rune, len(selected))
			for i, r := range selected {
				runes[i] = r.Rune()
			}
			result += string(runes)
			result += "\n"
		}
	}
	return result
}

type pressureType string

const (
	flood   pressureType = "flood"
	flow    pressureType = "flow"
	blocked pressureType = "blocked"
)

func parsePoints(lines []string) (points []grid.Point, topLeft grid.Point, bottomRight grid.Point) {
	pattern := regexp.MustCompile("(\\w)=(\\d+), (\\w)=(\\d+)..(\\d+)")
	for _, line := range lines {
		if m := pattern.FindStringSubmatch(line); len(m) < 6 {
			panic("input line failed to match")
		} else if m[1] == "x" {
			x := dayless.ParseInt(m[2])
			fromY := dayless.ParseInt(m[4])
			toY := dayless.ParseInt(m[5])
			for y := fromY; y <= toY; y++ {
				points = append(points, grid.Point{
					X: x,
					Y: y,
				})
			}
		} else {
			y := dayless.ParseInt(m[2])
			fromX := dayless.ParseInt(m[4])
			toX := dayless.ParseInt(m[5])
			for x := fromX; x <= toX; x++ {
				points = append(points, grid.Point{
					X: x,
					Y: y,
				})
			}
		}
	}

	minX := math.MaxInt16
	minY := math.MaxInt16
	maxX := math.MinInt16
	maxY := math.MinInt16
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return points, grid.Point{X: minX, Y: minY}, grid.Point{X: maxX, Y: maxY}
}

func buildTiles(clays []grid.Point) (grid ground) {
	minX := math.MaxInt16
	minY := math.MaxInt16
	maxX := math.MinInt16
	maxY := math.MinInt16
	claysPerRow := make(map[int][]int)
	for _, p := range clays {
		x := p.X
		y := p.Y
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		if v, exist := claysPerRow[y]; exist {
			claysPerRow[y] = append(v, x)
		} else {
			claysPerRow[y] = []int{x}
		}
	}

	maxX++
	maxY++

	maxX += 10
	maxY += 10

	grid = make([][]cellType, maxY)
	for y := 0; y < maxY; y++ {
		row := make([]cellType, maxX)
		for x := 0; x < maxX; x++ {
			if claysInRow, exist := claysPerRow[y]; exist && contains(claysInRow, x) {
				row[x] = clay
			} else {
				row[x] = sand
			}
		}
		grid[y] = row
	}

	return grid
}

func drainWater3(ground ground, springX int, springY int) {

	ground[springY][springX] = spring

	x := springX
	y := springY

	path := stack.New()
	path.Push(grid.Point{X: x, Y: y + 1})
	drainWater3Iterative(ground, path)
}

func drainWater3Iterative(ground ground, stack *stack.Stack) {

	visited := make(map[grid.Point]bool)

	for ; stack.Len() > 0; {
		current := stack.Peek().(grid.Point)
		visited[current] = true

		ground[current.Y][current.X] = dryWater

		downPoint := grid.Point{X: current.X, Y: current.Y + 1}
		if downPoint.Y < len(ground) {
			if _, alreadyVisited := visited[downPoint]; !alreadyVisited {
				visited[downPoint] = true
				switch ground[downPoint.Y][downPoint.X] {
				case sand, dryWater:
					stack.Push(downPoint)
					continue
				}
			} else {
				if ground[downPoint.Y][downPoint.X] == dryWater {
					// if down is still dry-water, than there is a endless flow
					stack.Pop()
					continue
				}
			}
		} else {
			stack.Pop()
			continue // end-of-grid, flowing forever, go back
		}

		leftPoint := grid.Point{X: current.X - 1, Y: current.Y}
		if leftPoint.X >= 0 {
			if _, alreadyVisited := visited[leftPoint]; !alreadyVisited {
				visited[leftPoint] = true
				switch ground[leftPoint.Y][leftPoint.X] {
				case sand, dryWater:
					stack.Push(leftPoint)
					continue
				}
			}
		} else {
			stack.Pop()
			continue // end-of-grid, flowing forever, go back
		}

		rightPoint := grid.Point{X: current.X + 1, Y: current.Y}
		if rightPoint.X < len(ground[rightPoint.Y]) {
			if _, alreadyVisited := visited[rightPoint]; !alreadyVisited {
				visited[rightPoint] = true
				switch ground[rightPoint.Y][rightPoint.X] {
				case sand, dryWater:
					stack.Push(rightPoint)
					continue
				}
			}
		} else {
			stack.Pop()
			continue // end-of-grid, flowing forever, go back
		}

		leftBlocked := -1
		for x := current.X; 0 <= x; x-- {
			s := ground[current.Y][x]
			if s == sand {
				break
			} else if s == clay {
				leftBlocked = x
				break
			}
		}
		rightBlocked := -1
		for x := current.X + 1; x < len(ground[current.Y]); x++ {
			s := ground[current.Y][x]
			if s == sand {
				break
			} else if s == clay {
				rightBlocked = x
				break
			}
		}
		if leftBlocked > -1 && rightBlocked > -1 {
			for x := leftBlocked + 1; x < rightBlocked; x++ {
				ground[current.Y][x] = settledWater
			}
		}
		stack.Pop()
	}
}

func drainWater1(ground ground, springX int, springY int) {

	ground[springY][springX] = spring

	x := springX
	y := springY

	path := grid.Points{grid.Point{X: x, Y: y + 1}}
	drainWater1Recursive(ground, path)
}

func drainWater2(ground ground, springX int, springY int, topLeftPoint grid.Point, bottomRightPoint grid.Point) {

	ground[springY][springX] = spring
	drainWater2Recursive(ground, grid.Point{X: springX, Y: springY}, topLeftPoint, bottomRightPoint)
}

func drainWater1Recursive(ground ground, path grid.Points) (pressureType pressureType) {

	current := path[len(path)-1]

	switch ground[current.Y][current.X] {
	case clay:
		return blocked
	}

	ground[current.Y][current.X] = dryWater

	downPoint := grid.Point{X: current.X, Y: current.Y + 1}
	if !(downPoint.Y < len(ground)) {
		return flow
	}
	downPath := make(grid.Points, len(path)+1)
	copy(downPath, path)
	downPath[len(path)] = downPoint
	downPressure := drainWater1Recursive(ground, downPath)
	switch downPressure {
	case flow:
		return flow
	}

	leftPoint := grid.Point{X: current.X - 1, Y: current.Y}
	leftPressure := flow
	leftAlreadyVisited := containsPath(path, leftPoint)
	if !leftAlreadyVisited {
		leftPath := make(grid.Points, len(path)+1)
		copy(leftPath, path)
		leftPath[len(path)] = leftPoint
		leftPressure = drainWater1Recursive(ground, leftPath)
	}

	rightPoint := grid.Point{X: current.X + 1, Y: current.Y}
	rightPressure := flow
	rightAlreadyVisited := containsPath(path, rightPoint)
	if !rightAlreadyVisited {
		rightPath := make(grid.Points, len(path)+1)
		copy(rightPath, path)
		rightPath[len(path)] = rightPoint
		rightPressure = drainWater1Recursive(ground, rightPath)
	}

	if !leftAlreadyVisited && rightAlreadyVisited {
		return leftPressure
	} else if leftAlreadyVisited && !rightAlreadyVisited {
		return rightPressure
	} else if !leftAlreadyVisited && !rightAlreadyVisited {
		if leftPressure == blocked && rightPressure == blocked {
			for xi := 0; ground[current.Y][current.X-xi] == dryWater; xi++ {
				ground[current.Y][current.X-xi] = settledWater
			}
			for xi := 1; ground[current.Y][current.X+xi] == dryWater; xi++ {
				ground[current.Y][current.X+xi] = settledWater
			}
			return flood
		}
	}
	return flow
}

func drainWater2Recursive(ground ground, p grid.Point, topLeftPoint grid.Point, bottomRightPoint grid.Point) {
	if p.Y > bottomRightPoint.Y {
		// out-of-grid
		return
	}
	if ground[p.Y][p.X] != spring && !(isOpen(ground, p)) {
		// no space for water
		return
	}
	if !(isOpen(ground, grid.Point{X: p.X, Y: p.Y + 1})) {
		leftX := p.X
		for ; isOpen(ground, grid.Point{X: leftX, Y: p.Y}) && !isOpen(ground, grid.Point{X: leftX, Y: p.Y + 1}); {
			ground[p.Y][leftX] = dryWater
			if leftX < topLeftPoint.X {
				break
			}
			leftX--
		}
		rightX := p.X + 1
		for ; isOpen(ground, grid.Point{X: rightX, Y: p.Y}) && !isOpen(ground, grid.Point{X: rightX, Y: p.Y + 1}); {
			ground[p.Y][rightX] = dryWater
			if rightX > bottomRightPoint.X {
				break
			}
			rightX++
		}
		if isOpen(ground, grid.Point{X: leftX, Y: p.Y + 1}) || isOpen(ground, grid.Point{X: rightX, Y: p.Y + 1}) {
			// any open -> drain again (recursive)
			drainWater2Recursive(ground, grid.Point{X: leftX, Y: p.Y}, topLeftPoint, bottomRightPoint)
			drainWater2Recursive(ground, grid.Point{X: rightX, Y: p.Y}, topLeftPoint, bottomRightPoint)
		} else if ground[p.Y][leftX] == clay && ground[p.Y][rightX] == clay {
			for i := leftX + 1; i < rightX; i++ {
				ground[p.Y][i] = settledWater
			}
		}
	} else {
		// open -> fall through
		ground[p.Y][p.X] = dryWater
		drainWater2Recursive(ground, grid.Point{X: p.X, Y: p.Y + 1}, topLeftPoint, bottomRightPoint)
		if ground[p.Y+1][p.X] == settledWater {
			drainWater2Recursive(ground, p, topLeftPoint, bottomRightPoint) // restart because of new data under
		}
	}
}

func isOpen(g ground, p grid.Point) bool {
	return g[p.Y][p.X] == sand || g[p.Y][p.X] == dryWater
}

func countGroundType(ground ground, topLeftPoint grid.Point, bottomRightPoint grid.Point, search cellType) (total int) {
	total = 0

	for y := topLeftPoint.Y; y <= bottomRightPoint.Y; y++ {
		for x := 0; x < len(ground[y]); x++ {
			if ground[y][x] == search {
				total++
			}
		}
	}

	return total
}

func containsPath(array []grid.Point, search grid.Point) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}

func contains(array []int, search int) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}
