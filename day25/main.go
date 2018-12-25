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

	// list of all points which have to be processed (goal is to clear this list)
	unprocessed := list.New()
	for _, p := range points {
		unprocessed.PushBack(p)
	}

	for unprocessed.Len() > 0 {
		// This is starting a completely new constellation with index 'c'
		c := len(constellations)
		constellations[c] = make([]Point, 0)

		// Starting with the next available point (and adding it into the current constellations)...
		nextElement := unprocessed.Front()
		startPoint := nextElement.Value.(Point)
		unprocessed.Remove(nextElement)
		constellations[c] = append(constellations[c], startPoint)

		for { // do.. while foundAny && unprocessed > 0
			// search for unprocessed points which will belongs to this constellation...
			foundAny := false
			morePoints := list.New()
			for _, cPoint := range constellations[c] {
				// which means any constellation's point can be used for looking up the distance to any unprocessed one...
				size := unprocessed.Len()
				for i := unprocessed.Front(); size > 0 && i != nil; size-- { // ensure we are running only once for all max
					elem := i // remove must be done after $next()
					point := i.Value.(Point)
					i = i.Next()
					if getManhattenDistance(cPoint, point) <= 3 {
						morePoints.PushBack(point)
						unprocessed.Remove(elem)
					}
				}
			}
			// If at least one additional point has been found (constellation has been grown), we will have to search at least one time again...
			for i := morePoints.Front(); i != nil; i = i.Next() {
				foundAny = true
				constellations[c] = append(constellations[c], i.Value.(Point))
			}
			// But: Stop whenever nothing has been found (no other unprocessed point is matching this constellation) or there aren't any left
			if !foundAny || unprocessed.Len() == 0 {
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
