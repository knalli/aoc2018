package main

import (
	"de.knallisworld/aoc/aoc2018/day22/lib"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"github.com/RyanCarrier/dijkstra"
	"strings"
	"time"
)

const AocDay = 22
const AocDayName = "day22"

const (
	rocky   int = 0
	wet     int = 1
	narrow  int = 2
	gear    int = 0
	torch   int = 1
	neither int = 2
)

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	depth, target := parseInput(lines)
	grid := buildCave(depth, target, 1000)
	riskLevel := getRiskLevel(grid, target)
	fmt.Printf("The total risk level: %d\n", riskLevel)
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Printf("Fewest number of minutes to target (aka lowest cost): %d\n", findFastestWay(grid, *grid.GetPoint(0, 0), target))
	fmt.Println()
}

func parseInput(lines []string) (depth int, target lib.Point) {
	depth = dayless.ParseInt(strings.Split(lines[0], ": ")[1])
	coordinates := strings.Split(lines[1], ": ")[1]
	target = lib.Point{
		X:       dayless.ParseInt(strings.Split(coordinates, ",")[0]),
		Y:       dayless.ParseInt(strings.Split(coordinates, ",")[1]),
		Erosion: resolveErosionLevel(0, depth),
	}
	return depth, target
}

func buildCave(depth int, target lib.Point, additions int) lib.Grid {
	grid := lib.NewGrid(lib.NewPoint(0, 0, resolveErosionLevel(0, depth)), &target)
	for y := 0; y < grid.Height(); y++ {
		for x := 0; x < grid.Width(); x++ {
			if grid.GetPoint(x, y) != nil {
				continue
			}
			if y == 0 {
				grid.AddPoint(lib.NewPoint(x, y, resolveErosionLevel(x*16807, depth)))
			} else if x == 0 {
				grid.AddPoint(lib.NewPoint(x, y, resolveErosionLevel(y*48271, depth)))
			} else {
				grid.AddPoint(lib.NewPoint(x, y, resolveErosionLevel(grid.GetPoint(x-1, y).Erosion*grid.GetPoint(x, y-1).Erosion, depth)))
			}
		}
	}
	return grid.Expand(additions, func(g lib.Grid, x int, y int) *lib.Point {
		if y == 0 {
			return lib.NewPoint(x, y, resolveErosionLevel(x*16807, depth))
		} else if x == 0 {
			return lib.NewPoint(x, y, resolveErosionLevel(y*48271, depth))
		} else {
			return lib.NewPoint(x, y, resolveErosionLevel(g.GetPoint(x-1, y).Erosion*g.GetPoint(x, y-1).Erosion, depth))
		}
	})
}

func resolveErosionLevel(geologicIndex int, depth int) int {
	return (geologicIndex + depth) % 20183
}

func getRiskLevel(grid lib.Grid, target lib.Point) (total int) {
	total = 0
	for y := 0; y <= target.Y; y++ {
		for x := 0; x <= target.X; x++ {
			total += grid.GetPoint(x, y).Erosion % 3
		}
	}
	return total
}

func findFastestWay(grid lib.Grid, start lib.Point, target lib.Point) int64 {

	// as the graph provides only "string" as key, this identifies both coordinates and tool: $y_$x_$tool
	payloadTpl := "%d_%d_%d"

	graph := dijkstra.NewGraph()
	regionTools := make(map[int][]int)
	for _, region := range []int{rocky, wet, narrow} {
		regionTools[region] = validToolsByRegion(region)
	}

	for y := 0; y < grid.Height(); y++ {
		for x := 0; x < grid.Width(); x++ {
			currTools, _ := regionTools[grid.GetPoint(x, y).Erosion%3]
			nextTools := []int{gear, torch, neither}
			for _, currTool := range currTools {
				for _, nextTool := range nextTools {
					if nextTool == currTool {
						continue
					}
					currId := fmt.Sprintf(payloadTpl, y, x, currTool)
					nextId := fmt.Sprintf(payloadTpl, y, x, nextTool)
					graph.AddMappedVertex(currId)
					graph.AddMappedVertex(nextId)
					_ = graph.AddMappedArc(currId, nextId, 7)
				}
			}
			grid.ForEachAdjacent(x, y, func(p lib.Point) {
				currTools, _ := regionTools[grid.GetPoint(x, y).Erosion%3]
				nextTools, _ := regionTools[grid.GetPoint(p.X, p.Y).Erosion%3]
				for _, currTool := range currTools {
					for _, nextTool := range nextTools {
						if nextTool == currTool {
							currId := fmt.Sprintf(payloadTpl, y, x, currTool)
							nextId := fmt.Sprintf(payloadTpl, p.Y, p.X, nextTool)
							graph.AddMappedVertex(currId)
							graph.AddMappedVertex(nextId)
							_ = graph.AddMappedArc(currId, nextId, 1)
						}
					}
				}
			})
		}
	}

	startIdx, err := graph.GetMapping(fmt.Sprintf(payloadTpl, 0, 0, torch))
	if err != nil {
		panic(err)
	}
	targetIdx, err := graph.GetMapping(fmt.Sprintf(payloadTpl, target.Y, target.X, torch))
	if err != nil {
		panic(err)
	}
	path, err := graph.Shortest(startIdx, targetIdx)
	if err != nil {
		panic(err)
	}

	return path.Distance
}

func validToolsByRegion(region int) []int {
	switch region {
	case rocky:
		return []int{gear, torch}
	case wet:
		return []int{gear, neither}
	case narrow:
		return []int{torch, neither}
	}
	panic("invalid combo")
}
