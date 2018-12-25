package main

import (
	"container/list"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strings"
	"time"
)

const AocDay = 25
const AocDayName = "day25"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")

	dayless.PrintStepHeader(1)
	fmt.Printf("🎉 Found %d constellations\n", solve(parsePoints(lines)))
	fmt.Println()
}

type Point struct {
	X int
	Y int
	Z int
	T int
}

func (p *Point) ToString() string {
	return fmt.Sprintf("%d,%d,%d,%d", p.X, p.Y, p.Z, p.T)
}

func parsePoints(lines []string) (result []Point) {
	for _, line := range lines {
		parts := strings.Split(line, ",")
		result = append(result, Point{
			X: dayless.ParseInt(parts[0]),
			Y: dayless.ParseInt(parts[1]),
			Z: dayless.ParseInt(parts[2]),
			T: dayless.ParseInt(parts[3]),
		})
	}
	return result
}

func solve(points []Point) int {
	constellations := make(map[int][]Point)

	unprocessed := list.New()
	for _, p := range points {
		unprocessed.PushBack(p)
	}

	for ; unprocessed.Len() > 0; {
		c := len(constellations)
		constellations[c] = make([]Point, 0)
		nextElement := unprocessed.Front()
		startPoint := nextElement.Value.(Point)
		unprocessed.Remove(nextElement)
		constellations[c] = append(constellations[c], startPoint)
		for {
			found := false
			morePoints := list.New()
			for _, cPoint := range constellations[c] {
				size := unprocessed.Len()
				for n := unprocessed.Front(); size > 0 && n != nil; size-- {
					point := n.Value.(Point)
					t := n // remove must be done after $next()
					n = n.Next()
					// fmt.Printf("%s -> %s [%d]\n", cPoint.ToString(), point.ToString(), getManhattenDistance(cPoint, point))
					if getManhattenDistance(cPoint, point) <= 3 {
						morePoints.PushBack(point)
						unprocessed.Remove(t)
					}
				}
			}
			for i := morePoints.Front(); i != nil; i = i.Next() {
				found = true
				constellations[c] = append(constellations[c], i.Value.(Point))
			}
			if !found || unprocessed.Len() == 0 {
				break
			}
		}
	}

	return len(constellations)
}

func getManhattenDistance(from Point, to Point) int {
	return intAbs(to.X-from.X) + intAbs(to.Y-from.Y) + intAbs(to.Z-from.Z) + intAbs(to.T-from.T)
}

func intAbs(n int) int {
	if n > 0 {
		return n
	} else {
		return -n
	}
}
