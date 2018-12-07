package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"errors"
	"fmt"
	"regexp"
	"time"
)

const AocDay = 7
const AocDayName = "day07"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	instructions, _ := readInstructions(lines)
	rInstructions := reduceInstructions(instructions)
	fmt.Println(instructions)
	fmt.Println(rInstructions)
	fmt.Printf("Order of instructions: %s\n", getOrderOfSteps(rInstructions))
	fmt.Println()

	dayless.PrintStepHeader(2)
	seconds, order := getOrderOfStepsParallel(rInstructions, 5)
	fmt.Printf("Order of instructions (%d seconds): %s\n", seconds, order)
	fmt.Println()
}

type instruction struct {
	requirement string
	step        string
}

func readInstructions(lines []string) ([]instruction, error) {
	// Step Z must be finished before step U can begin.
	pattern := regexp.MustCompile("Step ([A-Z]) must be finished before step ([A-Z]) can begin\\.")

	result := make([]instruction, len(lines))

	for i, line := range lines {
		matches := pattern.FindStringSubmatch(line)
		if len(matches) != 3 {
			return nil, errors.New("failed to match input line")
		}
		result[i] = instruction{requirement: matches[1], step: matches[2]}
	}

	return result, nil
}

func reduceInstructions(instructions []instruction) map[string][]string {
	result := make(map[string][]string)

	for _, instruction := range instructions {
		result[instruction.step] = append(result[instruction.step], instruction.requirement)
	}

	return result
}

func getOrderOfSteps(instructions map[string][]string) string {

	stepsReady := make([]string, 0)

	// START
	for i := 'A'; i <= 'Z'; i++ {
		step := string(i)
		if len(instructions[step]) == 0 {
			stepsReady = append(stepsReady, step)
			break
		}
	}

	return getOrderOfSteps2(instructions, stepsReady)
}

func getOrderOfSteps2(instructions map[string][]string, stepsReady []string) string {

	stepsLeft := false

	for i := 'A'; i <= 'Z'; i++ {
		step := string(i)

		// already visited?
		if contains(stepsReady, step) {
			continue
		} else {
			// at least one step not visited
			stepsLeft = true
		}

		requiredSteps := instructions[step]
		allRequiredStepsAreReady := true
		for _, requiredStep := range requiredSteps {
			if !contains(stepsReady, requiredStep) {
				allRequiredStepsAreReady = false
				break
			}
		}

		if allRequiredStepsAreReady {
			stepsReady = append(stepsReady, step)
			break
		}
	}

	if stepsLeft {
		return getOrderOfSteps2(instructions, stepsReady)
	} else {
		// build final string
		result := ""
		for _, l := range stepsReady {
			result += l
		}
		return result
	}
}

func contains(array []string, search string) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}

type worker struct {
	step        string
	availableAt int
}

func getOrderOfStepsParallel(instructions map[string][]string, workerSize int) (int, string) {

	workers := make([]worker, workerSize)
	stepsPicked := make([]string, 0)
	stepsReady := make([]string, 0)

	// START
	for i := 'A'; i <= 'Z'; i++ {
		step := string(i)
		if len(instructions[step]) == 0 {
			stepsPicked = append(stepsPicked, step)
			workers[0].step = step
			workers[0].availableAt = int(i - 4) // A==65 => 60+1
			break
		}
	}

	fmt.Print("Second ")
	for i := range workers {
		fmt.Printf("w%d ", i)
	}
	fmt.Println("Done")

	second := 0
	for ; ; second++ {

		numberOfBusyWorkers := 0
		for i := range workers {
			worker := &workers[i]

			// worker still busy?
			if worker.availableAt > second {
				numberOfBusyWorkers++
				continue
			}

			// any worker ready now?
			if worker.availableAt == second && worker.step != "" {
				stepsReady = append(stepsReady, worker.step)
				worker.step = ""
			}
		}

		if numberOfBusyWorkers < len(workers) {
			// at last one worker is free for work

			for i := 'A'; i <= 'Z'; i++ {
				step := string(i)

				// already visited?
				// if contains(stepsReady, step) {
				// 	continue
				// }
				// already picked?
				if contains(stepsPicked, step) {
					continue
				}

				requiredSteps := instructions[step]
				allRequiredStepsAreReady := true
				for _, requiredStep := range requiredSteps {
					if !contains(stepsReady, requiredStep) {
						allRequiredStepsAreReady = false
						break
					}
				}

				if allRequiredStepsAreReady {
					picked := false
					for w := range workers {
						worker := &workers[w]
						if worker.availableAt <= second {
							worker.step = step
							worker.availableAt = second + int(i-4)
							picked = true
							break
						}
					}
					if picked {
						stepsPicked = append(stepsPicked, step)
						numberOfBusyWorkers++
					} else {
						// looks like all workers are busy -> break loop
						break
					}
				}
			}

		}

		fmt.Printf("%6d ", second)
		for _, worker := range workers {
			step := worker.step
			if step == "" {
				step = "."
			}
			fmt.Printf("%2s ", step)
		}
		for _, l := range stepsReady {
			fmt.Printf(l)
		}
		fmt.Println()

		if numberOfBusyWorkers == 0 {
			break
		}

	}

	// build final string
	result := ""
	for _, l := range stepsReady {
		result += l
	}
	return second, result
}
