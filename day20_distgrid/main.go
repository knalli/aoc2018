package main

import (
	d15Grid "de.knallisworld/aoc/aoc2018/day15/grid"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"gopkg.in/karalabe/cookiejar.v1/exts/mathext"
	"math"
	"time"
)

const AocDay = 20
const AocDayName = "day20"

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

func compute(regex string, dimension int, countMinDoors int) {
	fmt.Printf("👉 Regex: %s\n", regex)
	start := d15Grid.NewPoint(dimension/2, dimension/2)
	distances := make(map[d15Grid.Point]int)
	visit(distances, start, regex, 0)

	// look for the max distance
	distance := math.MinInt16
	totalRoomsMinDistance := 0
	for _, d := range distances {
		if d > distance {
			distance = d
		}
		if d >= countMinDoors {
			totalRoomsMinDistance++
		}
	}
	fmt.Printf("🎉 Largest number of doors: ∂=%d\n", distance)
	fmt.Printf("🎉 At least %d rooms pass %d doors\n", totalRoomsMinDistance, countMinDoors)
}

func visit(distances map[d15Grid.Point]int, start d15Grid.Point, regex string, position int) (d15Grid.Point, int) {

	point := start

	for {
		switch regex[position] {
		case '^':
			// start
			position++
		case '$':
			// end-of-line
			position++
			return point, position
		case 'N':
			next := point.Top()
			if nextDist, exist := distances[next]; exist {
				distances[next] = mathext.MinInt(nextDist, distances[point]+1)
			} else {
				distances[next] = distances[point] + 1
			}
			point = next
			position++
		case 'E':
			next := point.Right()
			if nextDist, exist := distances[next]; exist {
				distances[next] = mathext.MinInt(nextDist, distances[point]+1)
			} else {
				distances[next] = distances[point] + 1
			}
			point = next
			position++
		case 'S':
			next := point.Bottom()
			if nextDist, exist := distances[next]; exist {
				distances[next] = mathext.MinInt(nextDist, distances[point]+1)
			} else {
				distances[next] = distances[point] + 1
			}
			point = next
			position++
		case 'W':
			next := point.Left()
			if nextDist, exist := distances[next]; exist {
				distances[next] = mathext.MinInt(nextDist, distances[point]+1)
			} else {
				distances[next] = distances[point] + 1
			}
			point = next
			position++
		case '(':
			position++
			point, position = visit(distances, point, regex, position)
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
