package lib

import (
	"fmt"
)

type Point struct {
	X       int
	Y       int
	Erosion int
}

func (p Point) ToString() string {
	return fmt.Sprintf("(%d,%d) erosion=%d", p.X, p.Y, p.Erosion)
}

func NewPoint(x int, y int, erosion int) *Point {
	return &Point{X: x, Y: y, Erosion: erosion}
}

type Grid struct {
	Points [][]*Point
}

func NewGrid(from *Point, to *Point) Grid {
	points := make([][]*Point, to.Y-from.Y+1)
	for i := range points {
		points[i] = make([]*Point, to.X-from.X+1)
	}
	g := Grid{Points: points}
	g.AddPoint(from)
	g.AddPoint(to)
	return g
}

func (g Grid) AddPoint(p *Point) {
	g.Points[p.Y][p.X] = p
}

func (g Grid) GetPoint(x int, y int) *Point {
	return g.Points[y][x]
}

func (g Grid) Height() int {
	return len(g.Points)
}

func (g Grid) Width() int {
	return len(g.Points[0])
}

func (g Grid) Expand(size int, buildPoint func(g Grid, x int, y int) *Point) Grid {
	newHeight := g.Height() + size
	newWidth := g.Width() + size
	points := make([][]*Point, newHeight)
	for i := range g.Points {
		points[i] = make([]*Point, newWidth)
		copy(points[i], g.Points[i])
	}
	for i := 0; i < size; i++ {
		points[g.Height()+i] = make([]*Point, newWidth)
	}
	r := Grid{Points: points}
	for y := 0; y < r.Height(); y++ {
		for x := 0; x < r.Width(); x++ {
			if r.GetPoint(x, y) == nil {
				r.AddPoint(buildPoint(r, x, y))
			}
		}
	}
	return r
}

func (g Grid) ForEachAdjacent(x int, y int, f func(p Point)) {
	g.ForEachAdjacentBreakable(x, y, func(p Point) bool {
		f(p)
		return true
	})
}

func (g Grid) ForEachAdjacentBreakable(x int, y int, f func(p Point) bool) *Point {
	height := g.Height()
	width := g.Width()
	if y > 0 {
		p := g.GetPoint(x, y-1)
		if !f(*p) {
			return p
		}
	}
	if x > 0 {
		p := g.GetPoint(x-1, y)
		if !f(*p) {
			return p
		}
	}
	if x+1 < width {
		p := g.GetPoint(x+1, y)
		if !f(*p) {
			return p
		}
	}
	if y+1 < height {
		p := g.GetPoint(x, y+1)
		if !f(*p) {
			return p
		}
	}
	return nil
}
