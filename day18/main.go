package main

import (
	"crypto"
	_ "crypto/sha1"
	"de.knallisworld/aoc/aoc2018/dayless"
	b64 "encoding/base64"
	"fmt"
	"hash"
	"time"
)

const AocDay = 18
const AocDayName = "day18"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")

	dayless.PrintStepHeader(1)
	grid := collectTicks(buildGrid(lines), 10)
	totalWood := countSearch(grid, tree)
	totalLumber := countSearch(grid, lumberjack)
	fmt.Printf("Finally, there are %d woord acres and %d lumberyards: %d * %d = %d\n", totalWood, totalLumber, totalWood, totalLumber, totalWood*totalLumber)
	fmt.Println()

	dayless.PrintStepHeader(2)
	grid = collectTicks(buildGrid(lines), 1000000000)
	totalWood = countSearch(grid, tree)
	totalLumber = countSearch(grid, lumberjack)
	fmt.Printf("Finally, there are %d woord acres and %d lumberyards: %d * %d = %d\n", totalWood, totalLumber, totalWood, totalLumber, totalWood*totalLumber)
	fmt.Println()
}

func buildGrid(lines []string) (grid [][]rune) {
	grid = make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid
}

const open = '.'
const tree = '|'
const lumberjack = '#'

func collectTicks(base [][]rune, iterations int) (grid [][]rune) {

	grid = base
	data := make(map[string]int)

	hashFn := crypto.SHA1.New()
	hashFn.Reset()

	for i := 0; i < iterations; i++ {
		grid = tick(grid)
		h := getHash(grid, hashFn)
		if lastIteration, exist := data[h]; exist {
			fmt.Printf("👉 Hey, this transformation result has been already made (last iteration = %d, now %d)\n", lastIteration, i)
			repeatingSteps := i - lastIteration
			i += ((iterations - i) / repeatingSteps) * repeatingSteps // jump ahead
		}
		data[h] = i
	}
	return grid
}

func tick(base [][]rune) (grid [][]rune) {
	grid = make([][]rune, len(base))
	for y, row := range base {
		grid[y] = make([]rune, len(row))
		copy(grid[y], base[y])
		for x, r := range row {
			switch r {
			case open:
				if countAdjacentSearch(y, x, base, tree) >= 3 {
					grid[y][x] = tree
				}
			case tree:
				if countAdjacentSearch(y, x, base, lumberjack) >= 3 {
					grid[y][x] = lumberjack
				}
			case lumberjack:
				if !(countAdjacentSearch(y, x, base, lumberjack) > 0 && countAdjacentSearch(y, x, base, tree) > 0) {
					grid[y][x] = open
				}
			}
		}
	}
	return grid
}

func countSearch(grid [][]rune, search rune) (total int) {
	for _, row := range grid {
		for _, r := range row {
			if r == search {
				total++
			}
		}
	}
	return total
}

func countAdjacentSearch(y int, x int, grid [][]rune, search rune) (totalOpen int) {
	if y > 0 {
		if x > 0 && grid[y-1][x-1] == search {
			totalOpen++
		}
		if grid[y-1][x] == search {
			totalOpen++
		}
		if x+1 < len(grid[y]) && grid[y-1][x+1] == search {
			totalOpen++
		}
	}
	if x > 0 && grid[y][x-1] == search {
		totalOpen++
	}
	if x+1 < len(grid[y]) && grid[y][x+1] == search {
		totalOpen++
	}
	if y+1 < len(grid) {
		if x > 0 && grid[y+1][x-1] == search {
			totalOpen++
		}
		if grid[y+1][x] == search {
			totalOpen++
		}
		if x+1 < len(grid[y]) && grid[y+1][x+1] == search {
			totalOpen++
		}
	}
	return totalOpen
}

func getHash(grid [][]rune, hashFn hash.Hash) string {
	str := ""
	for _, row := range grid {
		str += string(row)
	}
	// base64 for visible only
	return b64.StdEncoding.EncodeToString(hashFn.Sum([]byte(str)))
}
