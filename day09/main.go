package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const AocDay = 9
const AocDayName = "day09"
const debug = false

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	line, _ := dayless.ReadFileToString(AocDayName + "/puzzle.txt")
	gameParamPlayers, gameParamMaxMarble, _ := extractParams(*line)
	players, lastMarbleWorth, highScore := playTheGame(gameParamPlayers, gameParamMaxMarble)
	fmt.Printf("%d players; last marble is worth %d points: high score is %d\n", players, lastMarbleWorth, highScore)
	fmt.Println()

	dayless.PrintStepHeader(2)
	players, lastMarbleWorth, highScore = playTheGame(gameParamPlayers, gameParamMaxMarble*100)
	fmt.Printf("%d players; last marble is worth %d points: high score is %d\n", players, lastMarbleWorth, highScore)
	fmt.Println()
}

type link struct {
	left  *link
	right *link
	value int
}

func extractParams(line string) (int, int, error) {
	parts := strings.Split(line, " ")
	players, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	lastMarble, err := strconv.Atoi(parts[6])
	if err != nil {
		return 0, 0, err
	}
	return players, lastMarble, nil
}

func playTheGame(numPlayers int, lastMarble int) (int, int, int) {

	defer dayless.TimeTrack(time.Now(), "play the game")

	// start
	zero := link{value: 0}
	zero.left = &zero
	zero.right = &zero
	current := &zero

	playerScores := make(map[int]int, numPlayers)

	for marble := 1; marble <= lastMarble; marble++ {
		currentPlayer := marble % numPlayers
		if marble%23 == 0 {
			// scoring
			playerScores[currentPlayer-1] += marble              // add current marble value as score
			target := current.left.left.left.left.left.left.left // removing marble 7*left
			playerScores[currentPlayer-1] += target.value
			// aka remove item in link list
			before := target.left
			after := target.right
			before.right = after
			after.left = before
			current = after
		} else {
			target := current.right.right // placing new marble 2*right
			// aka add item in link list
			before := target.left
			placed := &link{left: before, right: target, value: marble}
			before.right = placed
			target.left = placed
			current = placed
		}

		if debug {
			fmt.Printf("[%d] {%6d}", currentPlayer, playerScores[currentPlayer-1])
			printRow(&zero)
			fmt.Println()
		}
	}

	maxScore := 0
	for i := 0; i < numPlayers; i++ {
		if maxScore < playerScores[i] {
			maxScore = playerScores[i]
		}
	}

	return numPlayers, lastMarble, maxScore
}

func printRow(start *link) {
	fmt.Printf(" %2d", start.value)
	for j := start.right; j != start; j = j.right {
		fmt.Printf(" %2d", j.value)
	}
}
