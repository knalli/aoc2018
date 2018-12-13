package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"errors"
	"fmt"
	"sort"
	"time"
)

const AocDay = 13
const AocDayName = "day13"
const debug = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")

	dayless.PrintStepHeader(1)
	trackMap, carts := explodeToRunesAndCarts(lines)
	for {
		if debug {
			printTrackMap(trackMap, carts)
			fmt.Println()
		}
		if err := tick(trackMap, carts, false); err != nil {
			fmt.Println(err)
			break
		}
	}
	if debug {
		printTrackMap(trackMap, carts)
		fmt.Println()
	}
	fmt.Println()

	dayless.PrintStepHeader(2)
	trackMap, carts = explodeToRunesAndCarts(lines)
	for {
		if debug {
			printTrackMap(trackMap, carts)
			fmt.Println()
		}
		if err := tick(trackMap, carts, true); err != nil {
			fmt.Println(err)
			break
		}
	}
	if debug {
		printTrackMap(trackMap, carts)
		fmt.Println()
	}
	fmt.Println()
}

func printTrackMap(runes [][]rune, carts Carts) {

	// helper structure
	h := make([][]int, len(runes))
	for i := range h {
		h[i] = make([]int, len(runes[0]))
		for j := range h[i] {
			h[i][j] = -1
		}
	}
	for i := range carts {
		cart := &carts[i]
		if cart.disabled {
			continue
		}
		h[cart.y][cart.x] = i
	}

	for y := 0; y < len(runes); y++ {
		lineRunes := make([]rune, 0)
		for x := 0; x < len(runes[y]); x++ {
			if cartIdx := h[y][x]; cartIdx > -1 {
				lineRunes = append(lineRunes, carts[cartIdx].direction)
			} else {
				lineRunes = append(lineRunes, runes[y][x])
			}
		}
		fmt.Println(string(lineRunes))
	}
}

func explodeToRunesAndCarts(lines []string) ([][]rune, Carts) {

	maxLineLength := 0
	for _, line := range lines {
		m := len(line)
		if maxLineLength < m {
			maxLineLength = m
		}
	}

	runes := make([][]rune, len(lines))
	for i, line := range lines {
		lineRunes := []rune(line)
		for len(lineRunes) < maxLineLength {
			// ensure all lines have same length
			lineRunes = append(lineRunes, ' ')
		}
		runes[i] = lineRunes
	}

	carts := make([]Cart, 0)
	for y := 0; y < len(runes); y++ {
		for x := 0; x < len(runes[y]); x++ {
			direction := runes[y][x]
			switch direction {
			case '^':
				carts = append(carts, Cart{x: x, y: y, direction: direction})
				runes[y][x] = '|'
			case '>':
				carts = append(carts, Cart{x: x, y: y, direction: direction})
				runes[y][x] = '-'
			case 'v':
				carts = append(carts, Cart{x: x, y: y, direction: direction})
				runes[y][x] = '|'
			case '<':
				carts = append(carts, Cart{x: x, y: y, direction: direction})
				runes[y][x] = '-'
			}
		}
	}

	return runes, carts
}

type vector struct {
	x int
	y int
}

type Cart struct {
	x         int
	y         int
	direction rune
	turns     int
	disabled  bool
}

type Carts []Cart

func (c Cart) sortable() int {
	return c.y*1000 + c.x
}

func tick(runes [][]rune, carts Carts, avoidCrash bool) error {

	sort.Slice(carts, func(i, j int) bool {
		return carts[i].sortable() < carts[j].sortable()
	})

	for i := range carts {
		cart := &carts[i]
		if cart.disabled {
			continue
		}

		x := cart.x
		y := cart.y

		currDirection := cart.direction

		currVector := getVectorByDirection(currDirection)
		nextTrack := runes[y+currVector.y][x+currVector.x]

		// Check about possible crash?
		var crashedCart *Cart
		for j := range carts {
			if carts[j].disabled {
				continue
			}
			if carts[j].y == y+currVector.y && carts[j].x == x+currVector.x { // any cart on next coordinates?
				crashedCart = &carts[j]
			}
		}
		if crashedCart != nil {
			if avoidCrash {
				cart.disabled = true
				crashedCart.disabled = true
				continue
			} else {
				return errors.New(fmt.Sprintf("💥 Crash detected while driving (%d,%d) -> (%d,%d)", x, y, x+currVector.x, y+currVector.y))
			}
		}

		var nextCartDirection rune

		if nextTrack == '+' {
			nextRelativeDirection := []rune("<^>")[cart.turns%3]
			cart.turns += 1
			nextCartDirection = getEffectiveDirection(currDirection, nextRelativeDirection)
		} else {
			test := string([]rune{currDirection, nextTrack})
			switch test {
			case "^|", ">/", "<\\":
				nextCartDirection = '^'
			case "^/", ">-", "v\\":
				nextCartDirection = '>'
			case "^\\", "v/", "<-":
				nextCartDirection = '<'
			case ">\\", "v|", "</":
				nextCartDirection = 'v'
			default:
				panic(fmt.Sprintf("not expected: '%s' (%d,%d) [%c]", test, x, y, currDirection))
			}
		}

		// Move cart to next position
		cart.direction = nextCartDirection
		cart.x += currVector.x
		cart.y += currVector.y
	}

	if avoidCrash {
		processedCarts := 0
		for _, cart := range carts {
			if !cart.disabled {
				processedCarts++
			}
		}
		if processedCarts == 1 {
			for _, cart := range carts {
				if !cart.disabled {
					return errors.New(fmt.Sprintf("💥 Finally, only one cart is remainig at (%d,%d)", cart.x, cart.y))
				}
			}
		}
	}

	return nil
}

func getVectorByDirection(direction rune) vector {

	x := 0
	y := 0

	switch direction {
	case 'v':
		y = 1
	case '^':
		y = -1
	case '<':
		x = -1
	case '>':
		x = 1
	}

	return vector{x, y}
}

func getEffectiveDirection(direction rune, turn rune) rune {
	effective := '*'
	switch string([]rune{direction, turn}) {

	case "^^":
		effective = '^'
	case "^>":
		effective = '>'
	case "^<":
		effective = '<'

	case ">^":
		effective = '>'
	case ">>":
		effective = 'v'
	case "><":
		effective = '^'

	case "v^":
		effective = 'v'
	case "v>":
		effective = '<'
	case "v<":
		effective = '>'

	case "<^":
		effective = '<'
	case "<>":
		effective = '^'
	case "<<":
		effective = 'v'

	default:
		panic(fmt.Sprintf("not expected for: '%c', '%c", direction, turn))
	}

	return effective
}
