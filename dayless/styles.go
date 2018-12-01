package dayless

import "fmt"

func PrintDayHeader(day int) {
	fmt.Printf("Advent of Code 2018 - Day %02d\n", day)
	fmt.Println("================================================================")
	fmt.Println()
}

func PrintStepHeader(step int) {
	switch step {
	case 1:
		fmt.Println("--- Part One ---")
		break
	case 2:
		fmt.Println("--- Part Two ---")
		break
	default:
		fmt.Println("--- Part ??? ---")
	}
}
