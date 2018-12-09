package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/deque"
)

const AocDay = 9
const AocDayName = "day09_alternative"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	line, _ := dayless.ReadFileToString(AocDayName + "/puzzle.txt")
	gameParamPlayers, gameParamMaxMarble, _ := extractParams(*line)
	players, lastMarbleWorth, highScore := playTheGame2(gameParamPlayers, gameParamMaxMarble)
	fmt.Printf("%d players; last marble is worth %d points: high score is %d\n", players, lastMarbleWorth, highScore)
	fmt.Println()

	dayless.PrintStepHeader(2)
	players, lastMarbleWorth, highScore = playTheGame2(gameParamPlayers, gameParamMaxMarble*100)
	fmt.Printf("%d players; last marble is worth %d points: high score is %d\n", players, lastMarbleWorth, highScore)
	fmt.Println()
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

func playTheGame2(numPlayers int, lastMarble int) (int, int, int) {

	defer dayless.TimeTrack(time.Now(), "play the game")

	d := deque.New()
	d.Append(0)

	playerScores := make(map[int]int, numPlayers)

	for marble := 1; marble <= lastMarble; marble++ {
		currentPlayer := marble % numPlayers
		if marble%23 == 0 {
			// scoring
			playerScores[currentPlayer] += marble // add current marble value as score
			d.Rotate(-7)
			value, _ := d.Pop()
			playerScores[currentPlayer] += value.(int)
		} else {
			d.Rotate(2)
			d.Append(marble)
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
