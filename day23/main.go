package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"math"
	"regexp"
	"time"
)

const AocDay = 23
const AocDayName = "day23"
const trace = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	bots := parseLines(lines)

	dayless.PrintStepHeader(1)
	fmt.Printf("Nanobots in range of the strongest one: %d\n", len(findBotsInRange(bots, findStrongestNanobot(bots))))
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Printf("The shortest manhatten distance to any of the largest bot clusters: %d\n", findClosestCluster(bots))
	fmt.Println()
}

// pos=<4,0,0>, r=3
type NanoBot struct {
	X     int
	Y     int
	Z     int
	Range int
}

func (n NanoBot) ToString() string {
	return fmt.Sprintf("pos=<%d,%d,%d>, r=%d", n.X, n.Y, n.Z, n.Range)
}

func parseLines(lines []string) (bots []NanoBot) {
	pattern := regexp.MustCompile("pos=<(-?\\d+),(-?\\d+),(-?\\d+)>, r=(\\d+)")
	for _, line := range lines {
		matches := pattern.FindStringSubmatch(line)
		if len(matches) != 5 {
			panic("failed parsing line")
		}
		bots = append(bots, NanoBot{
			X:     dayless.ParseInt(matches[1]),
			Y:     dayless.ParseInt(matches[2]),
			Z:     dayless.ParseInt(matches[3]),
			Range: dayless.ParseInt(matches[4]),
		})
	}
	return bots
}

func findStrongestNanobot(bots []NanoBot) *NanoBot {
	c := -1
	maxRange := -1
	for i := range bots {
		bot := &bots[i]
		if bot.Range > maxRange {
			c = i
			maxRange = bot.Range
		}
	}
	return &bots[c]
}

func findBotsInRange(bots []NanoBot, bot *NanoBot) (result []*NanoBot) {
	for i := range bots {
		b := &bots[i]
		if getManhattenDistanceByBots(b, bot) <= bot.Range {
			result = append(result, b)
		}
	}
	return result
}

type Point struct {
	X int
	Y int
	Z int
}

func (p Point) ToString() string {
	return fmt.Sprintf("(%d,%d,%d)", p.X, p.Y, p.Z)
}

func findClosestCluster(bots []NanoBot) int {

	defer dayless.TimeTrack(time.Now(), "findClosestCluster")

	minX := math.MaxInt64
	maxX := math.MinInt64
	minY := math.MaxInt64
	maxY := math.MinInt64
	minZ := math.MaxInt64
	maxZ := math.MinInt64

	for i := range bots {
		bot := &bots[i]
		if bot.X < minX {
			minX = bot.X
		}
		if bot.X > maxX {
			maxX = bot.X
		}
		if bot.Y < minY {
			minY = bot.Y
		}
		if bot.Y > maxY {
			maxY = bot.Y
		}
		if bot.Z < minZ {
			minZ = bot.Z
		}
		if bot.Z > maxZ {
			maxZ = bot.Z
		}
	}
	if trace {
		fmt.Printf("👉 Problem area coordicates (X,Y,Z): min(%d,%d,%d), max=(%d,%d,%d)\n", minX, minY, minZ, maxX, maxY, maxZ)
	}

	// use binary search (splitting problem area in large sub-problems)
	dist := 1
	for ; dist < maxX-minX; {
		dist *= 2
	}

	for {
		pov := Point{0, 0, 0}
		povMinDistance := math.MaxInt64
		best := Point{0, 0, 0}
		bestTotal := 0

		if trace {
			fmt.Printf("⏳ (%d,%d,%d) - (%d,%d,%d) dist=%d last-best=%s, best-total=%d\n", minX, minY, minZ, maxX, maxY, maxZ, dist, best.ToString(), bestTotal)
		}

		for x := minX; x <= maxX; x += dist {
			for y := minY; y <= maxY; y += dist {
				for z := minZ; z <= maxZ; z += dist {
					total := 0
					for _, bot := range bots {
						d := getManhattenDistance(x, y, z, bot.X, bot.Y, bot.Z)
						if (d-bot.Range)/dist <= 0 { // "d-range" <=0 is in range, "/dist" ensures we are within our binary search area
							total++
						}
					}
					if total < bestTotal {
						continue
					} else if total > bestTotal {
						bestTotal = total
						povMinDistance = getManhattenDistance(x, y, z, pov.X, pov.Y, pov.Z)
						best = Point{X: x, Y: y, Z: z}
					} else {
						d := getManhattenDistance(x, y, z, pov.X, pov.Y, pov.Z)
						if d < povMinDistance {
							povMinDistance = d
							best = Point{X: x, Y: y, Z: z}
						}
					}
				}
			}
		}

		if trace {
			fmt.Printf("👉 (%d,%d,%d) - (%d,%d,%d) dist=%d last-best=%s, best-total=%d\n", minX, minY, minZ, maxX, maxY, maxZ, dist, best.ToString(), bestTotal)
		}

		if dist == 1 {
			return povMinDistance
		} else {
			// reduce by 2
			minX = best.X - dist
			maxX = best.X + dist
			minY = best.Y - dist
			maxY = best.Y + dist
			minZ = best.Z - dist
			maxZ = best.Z + dist
			dist /= 2
		}
	}
}

func getManhattenDistanceByBots(from *NanoBot, to *NanoBot) int {
	return getManhattenDistance(from.X, from.Y, from.Z, to.X, to.Y, to.Z)
}

func getManhattenDistance(fromX, fromY, fromZ, toX, toY, toZ int) int {
	return abs(fromX, toX) + abs(fromY, toY) + abs(fromZ, toZ)
}

func abs(from int, to int) int {
	return int(math.Abs(float64(to) - float64(from)))
}
