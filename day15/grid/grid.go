package grid

import (
	"fmt"
	queue2 "github.com/golang-collections/collections/queue"
	"github.com/jupp0r/go-priority-queue"
)

type Grid struct {
	fields [][]rune
}

func NewByChars(fields [][]rune) Grid {
	return Grid{fields: fields}
}

func NewByStrings(lines []string) Grid {
	fields := make([][]rune, len(lines))
	for i, line := range lines {
		fields[i] = []rune(line)
	}
	return Grid{fields: fields}
}

func (g Grid) Width() int {
	return len(g.fields[0])
}

func (g Grid) Height() int {
	return len(g.fields)
}

func (g Grid) GetValueAtXY(p Point) rune {
	return g.fields[p.Y][p.X]
}

func (g Grid) GetValueAtXYAsString(p Point) string {
	return string(g.GetValueAtXY(p))
}

func (g Grid) PutValueAtXY(p Point, value rune) {
	g.fields[p.Y][p.X] = value
}

func (g Grid) ToString() string {
	result := ""
	height := g.Height()
	for y := 0; y < height; y++ {
		result += string(g.fields[y])
		result += "\n"
	}
	return result
}

func (g Grid) ToStringWithBorder() string {
	result := "  "
	height := g.Height()
	width := g.Width()
	for i := 0; i < width; i++ {
		result += fmt.Sprintf("%d", i%10)
	}
	result += "\n"
	for y := 0; y < height; y++ {
		result += fmt.Sprintf("%d ", y%10)
		result += string(g.fields[y])
		result += "\n"
	}
	return result
}

func (g Grid) Walk(f func(p Point, value rune)) {
	g.WalkBreakable(func(p Point, value rune) bool {
		f(p, value)
		return true
	})
}

func (g Grid) WalkBreakable(f func(p Point, value rune) bool) bool {
	height := g.Height()
	width := g.Width()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !f(NewPoint(x, y), g.fields[y][x]) {
				return false
			}
		}
	}
	return true
}

func (g Grid) ForEachAdjacent(x int, y int, f func(p Point, value rune)) {
	g.ForEachAdjacentBreakable(x, y, func(p Point, value rune) bool {
		f(p, value)
		return true
	})
}

func (g Grid) ForEachAdjacentBreakable(x int, y int, f func(p Point, value rune) bool) *Point {
	height := g.Height()
	width := g.Width()
	if y > 0 {
		p := Point{X: x, Y: y - 1}
		v := g.GetValueAtXY(p)
		if !f(p, v) {
			return &p
		}
	}
	if x > 0 {
		p := Point{X: x - 1, Y: y}
		v := g.GetValueAtXY(p)
		if !f(p, v) {
			return &p
		}
	}
	if x+1 < width {
		p := Point{X: x + 1, Y: y}
		v := g.GetValueAtXY(p)
		if !f(p, v) {
			return &p
		}
	}
	if y+1 < height {
		p := Point{X: x, Y: y + 1}
		v := g.GetValueAtXY(p)
		if !f(p, v) {
			return &p
		}
	}
	return nil
}

func (g Grid) IsAdjacent(p1 Point, p2 Point) bool {
	if r := g.ForEachAdjacentBreakable(p1.X, p1.Y, func(p Point, _ rune) bool {
		return p != p2 // break if found
	}); r != nil {
		return *r == p2
	} else {
		return false
	}
}

type Point struct {
	X int
	Y int
}

func NewPoint(x int, y int) Point {
	return Point{X: x, Y: y}
}

func (p Point) ToString() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func (p Point) GetGridPosition(grid *Grid) int {
	return grid.Width()*p.Y + p.X
}

type Points []Point

func (p Points) ToString() string {
	result := ""
	for _, point := range p {
		if len(result) > 0 {
			result += ", "
		}
		result += point.ToString()
	}
	return result
}

func (g Grid) GetShortestPath_BFS(start Point, goal Point, isOpenFn func(value rune) bool, heuristic func(from Point, to Point) int) []Point {
	frontier := queue2.New()
	frontier.Enqueue(start)
	cameFrom := make(map[Point]*Point)
	for frontier.Len() > 0 {
		current := frontier.Dequeue().(Point)
		g.ForEachAdjacent(current.X, current.Y, func(next Point, value rune) {
			if isOpenFn(value) {
				if _, exist := cameFrom[next]; !exist {
					frontier.Enqueue(next)
					cameFrom[next] = &current
				}
			}
		})
	}

	// fmt.Printf("FINISHED\n")

	path := make(Points, 0)
	// path = append(path, goal)

	x := &goal
	for {
		next, exist := cameFrom[*x]
		if !exist || next == nil || *next == start || *next == goal {
			break
		}
		path = append(path, *next)
		x = next
	}

	// path = append(path, start)

	path2 := make(Points, len(path))
	for i := range path2 {
		path2[i] = path[len(path)-i-1]
	}

	// fmt.Printf("PATH: %s -> %s = %s\n", start.ToString(), goal.ToString(), path2.ToString())

	return path2
}

func (g Grid) GetShortestPath(start Point, goal Point, isOpenFn func(value rune) bool, heuristic func(from Point, to Point) int) []Point {
	frontier := queue2.New()
	frontier.Enqueue([]Point{start})
	cameFrom := make(map[Point]*Point)
	goals := []Point{goal}
	for frontier.Len() > 0 {
		path := frontier.Dequeue().([]Point)
		current := path[len(path)-1]
		if contains(goals, current) {
			return path
		}
		g.ForEachAdjacent(current.X, current.Y, func(next Point, value rune) {
			if isOpenFn(value) {
				if _, exist := cameFrom[next]; !exist {
					frontier.Enqueue(append(path, next))
					cameFrom[next] = &current
				}
			}
		})
	}

	return []Point{}
}

func (g Grid) GetShortestPathMulti(start Point, goals []Point, isOpenFn func(value rune) bool, heuristic func(from Point, to Point) int) []Point {
	frontier := queue2.New()
	frontier.Enqueue([]Point{start})
	cameFrom := make(map[Point]*Point)
	for frontier.Len() > 0 {
		path := frontier.Dequeue().([]Point)
		current := path[len(path)-1]
		if contains(goals, current) {
			return path
		}
		g.ForEachAdjacent(current.X, current.Y, func(next Point, value rune) {
			if isOpenFn(value) {
				if _, exist := cameFrom[next]; !exist {
					frontier.Enqueue(append(path, next))
					cameFrom[next] = &current
				}
			}
		})
	}

	return []Point{}
}

func contains(array []Point, search Point) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}

func (g Grid) GetShortestPath2(start Point, goal Point, isOpenFn func(value rune) bool, heuristic func(from Point, to Point) int) []Point {
	frontier := pq.New()
	frontier.Insert(start, 0)
	cameFrom := make(map[Point]*Point)
	costSoFar := make(map[Point]int)
	cameFrom[start] = nil
	costSoFar[start] = 0

	for frontier.Len() > 0 {
		x, _ := frontier.Pop()
		current := x.(Point)
		if current == goal {
			break
		}
		// fmt.Printf("  POINT: %s -> %s\n", current.ToString(), goal.ToString())
		g.ForEachAdjacent(current.X, current.Y, func(next Point, value rune) {
			// fmt.Printf("    ADJACENT: %s = %c\n", next.ToString(), value)
			if isOpenFn(value) {
				newCost := costSoFar[current] + 1
				if c, exist := costSoFar[next]; !exist || newCost < c {
					costSoFar[next] = newCost
					priority := newCost // + heuristic(goal, next)
					frontier.Insert(next, float64(priority))
					cameFrom[next] = &current
				}
			}
		})
	}

	// fmt.Printf("FINISHED\n")

	path := make(Points, 0)
	path = append(path, goal)

	x := &goal
	for {
		next, exist := cameFrom[*x]
		if !exist || next == nil || *next == start || *next == goal {
			break
		}
		path = append(path, *next)
		x = next
	}

	path = append(path, start)

	if len(path) == 2 {
		return []Point{}
	}

	path2 := make(Points, len(path))
	for i := range path2 {
		path2[i] = path[len(path)-i-1]
	}

	// fmt.Printf("PATH: %s -> %s = %s\n", start.ToString(), goal.ToString(), path2.ToString())

	return path2
}
