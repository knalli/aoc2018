package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const AocDay = 4
const AocDayName = "day04"

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, err := dayless.ReadFileToArray(AocDayName + "/Puzzle.txt")
	shifts, err := buildShiftMatrix(lines)
	if err != nil {
		panic(err)
	}
	mostAsleepGuard, asleepMinutes, topMinute := findMostSleepingGuard(shifts)
	fmt.Printf("Guard #%d spent %d minutes asleep, mostly at minute %d\n", mostAsleepGuard, asleepMinutes, topMinute)
	fmt.Printf("$guard * $$topMinute = %d\n", mostAsleepGuard*topMinute)
	fmt.Println()

	dayless.PrintStepHeader(2)
	mostFrequentlyGuard, mostFrequentlyMinute, _ := findMostFrequentlyAsleepOnSameMinute(shifts)
	fmt.Printf("$guard * $topMinuteTimes = %d\n", mostFrequentlyGuard*mostFrequentlyMinute)
	fmt.Println()
}

type shift struct {
	day     string
	guard   int
	asleeps [60]bool
}

func buildShiftMatrix(lines []string) (map[string]shift, error) {

	result := make(map[string]shift)

	sort.Strings(lines)
	reLine := regexp.MustCompile("\\[(\\d+-\\d+-\\d+) (\\d+):(\\d+)] (.*)")
	reGuardShift := regexp.MustCompile("Guard #(\\d+) begins shift")

	sleepingStart := -1
	for _, line := range lines {
		lineMatch := reLine.FindStringSubmatch(line)
		layout := "2006-01-02"
		date, err := time.Parse(layout, lineMatch[1])
		if err != nil {
			return nil, err
		}
		timeHour := lineMatch[2]
		timeMinute, err := strconv.Atoi(lineMatch[3])
		if err != nil {
			return nil, err
		}
		text := lineMatch[4]

		guardShiftMatch := reGuardShift.FindStringSubmatch(text)
		day := date.Format(layout)
		if len(guardShiftMatch) > 1 {
			id, err := strconv.Atoi(guardShiftMatch[1])
			if err != nil {
				return nil, err
			}
			if timeHour == "23" {
				// next day
				x := date.AddDate(0, 0, 1).Format(layout)
				result[x] = shift{guard: id}
			} else {
				result[day] = shift{guard: id}
			}
			sleepingStart = -1
		}

		if timeHour != "00" {
			// Because all asleep/awake times are during the midnight hour (00:00 - 00:59), only the minute portion (00 - 59) is relevant for those events.
			continue
		}

		if text == "falls asleep" {
			sleepingStart = timeMinute
		} else if text == "wakes up" {
			if sleepingStart == -1 {
				panic("error: wake up after no asleep")
			}
			sleeps := [60]bool{}
			for k, v := range result[day].asleeps {
				sleeps[k] = v
			}
			for i := sleepingStart; i < timeMinute; i++ {
				sleeps[i] = true
			}
			if result[day].guard == 0 {
				panic("Invalid guard #0")
			}
			result[day] = shift{guard: result[day].guard, asleeps: sleeps}
		} else {
			continue
		}
	}

	return result, nil
}

func findMostSleepingGuard(shifts map[string]shift) (guard int, minutes int, top int) {

	asleeps := make(map[int]int)
	for _, shift := range shifts {
		asleep := 0
		for _, b := range shift.asleeps {
			if b {
				asleep++
			}
		}
		asleeps[shift.guard] += asleep
	}

	maxAsleep := 0
	maxAsleepGuard := 0
	for guard, asleep := range asleeps {
		if asleep > maxAsleep {
			maxAsleep = asleep
			maxAsleepGuard = guard
		}
	}

	topMinute, _ := findMostTotalAsleepMinutes(maxAsleepGuard, shifts)
	return maxAsleepGuard, maxAsleep, topMinute
}

func findMostTotalAsleepMinutes(guard int, shifts map[string]shift) (int, int) {
	m := buildNumberOfMinutesByGuard(shifts, guard)

	max := 0
	idx := 0
	for k, v := range m {
		if v > max {
			max = v
			idx = k
		}
	}

	return idx, max
}

func buildNumberOfMinutesByGuard(shifts map[string]shift, guard int) [60]int {
	m := [60]int{}
	for _, shift := range shifts {
		if shift.guard != guard {
			continue
		}
		for i, b := range shift.asleeps {
			if b {
				m[i]++
			}
		}
	}
	return m
}

func findMostFrequentlyAsleepOnSameMinute(shifts map[string]shift) (guard int, minute int, times int) {

	maxAsleepMinute := 0
	maxAsleepMinuteTimes := 0
	maxAsleepMinuteGuard := 0

	set := make(map[int]struct{})
	mm := make(map[int][60]int)

	for _, shift := range shifts {
		if _, ok := set[shift.guard]; !ok {
			set[shift.guard] = struct{}{}
			mm[shift.guard] = buildNumberOfMinutesByGuard(shifts, shift.guard)

			minute, times := findMostTotalAsleepMinutes(shift.guard, shifts)
			if times > maxAsleepMinuteTimes {
				maxAsleepMinute = minute
				maxAsleepMinuteTimes = times
				maxAsleepMinuteGuard = shift.guard
			}
		}
	}

	return maxAsleepMinuteGuard, maxAsleepMinute, maxAsleepMinuteTimes

}
