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

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	bots := parseLines(lines)

	dayless.PrintStepHeader(1)
	fmt.Printf("Nanobots in range of the strongest one: %d\n", len(findBotsInRange(bots, findStrongestNanobot(bots))))
	fmt.Println()

	dayless.PrintStepHeader(2)
	fmt.Println("N/A")
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

func findBotsInRange(bots []NanoBot, bot *NanoBot) (result []NanoBot) {
	for i := range bots {
		b := &bots[i]
		if getManhattenDistance(b, bot) <= bot.Range {
			result = append(result, *b)
		}
	}
	return result
}

func getManhattenDistance(from *NanoBot, to *NanoBot) int {
	return abs(from.X, to.X) + abs(from.Y, to.Y) + abs(from.Z, to.Z)
}

func abs(from int, to int) int {
	return int(math.Abs(float64(from) - float64(to)))
}
