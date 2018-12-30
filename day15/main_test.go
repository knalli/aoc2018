package main

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"testing"
)

func assertWithoutBonus(t *testing.T, lines []string, expectedFullRounds int, expectedHitPointsSum int, expectedOutcome int) bool {
	return t.Run("Part1", func(t *testing.T) {
		actualFullRounds, actualHitPointsSum := part1(lines)
		actualOutcome := actualFullRounds * actualHitPointsSum
		if actualFullRounds != expectedFullRounds {
			t.Errorf("Part 1: Total of full rounds is not correct. Expected: %d, actual %d", expectedFullRounds, actualFullRounds)
		}
		if actualHitPointsSum != expectedHitPointsSum {
			t.Errorf("Part 1: Total hit points is not correct. Expected: %d, actual %d", expectedHitPointsSum, actualHitPointsSum)
		}
		if actualOutcome != expectedOutcome {
			t.Fatalf("Part 1: Outcome (solution) is not correct. Expected: %d, actual: %d", expectedOutcome, actualOutcome)
		}
	})
}

func assertWithBonus(t *testing.T, lines []string, expectedBonus int, expectedFullRounds int, expectedHitPointsSum int, expectedOutcome int) bool {
	return t.Run("Part2", func(t *testing.T) {
		actualBonus, actualFullRounds, actualHitPointsSum := part2(lines)
		actualOutcome := actualFullRounds * actualHitPointsSum
		if actualBonus != expectedBonus {
			t.Errorf("Part 2: Bonus is not correct: Expected: %d, actual: %d", expectedBonus, actualBonus)
		}
		if actualFullRounds != expectedFullRounds {
			t.Errorf("Part 2: Total of full rounds is not correct. Expected: %d, actual %d", expectedFullRounds, actualFullRounds)
		}
		if actualHitPointsSum != expectedHitPointsSum {
			t.Errorf("Part 2: Total hit points is not correct. Expected: %d, actual %d", expectedHitPointsSum, actualHitPointsSum)
		}
		if actualOutcome != expectedOutcome {
			t.Fatalf("Part 2: Outcome (solution) is not correct. Expected: %d, actual: %d", expectedOutcome, actualOutcome)
		}
	})
}

func TestSample00(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("sample00.txt")
	assertWithoutBonus(t, lines, 47, 590, 27730)
	assertWithBonus(t, lines, 12, 29, 172, 4988)
}

func TestSample01(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("sample01.txt")
	assertWithoutBonus(t, lines, 37, 982, 36334)
	assertWithBonus(t, lines, 1, 28, 1038, 29064)
}

func TestSample02(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("sample02.txt")
	assertWithoutBonus(t, lines, 46, 859, 39514)
	assertWithBonus(t, lines, 1, 33, 948, 31284)
}

func TestSample03(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("sample03.txt")
	assertWithoutBonus(t, lines, 35, 793, 27755)
	assertWithBonus(t, lines, 12, 37, 94, 3478)
}

func TestSample04(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("sample04.txt")
	assertWithoutBonus(t, lines, 54, 536, 28944)
	assertWithBonus(t, lines, 9, 39, 166, 6474)
}

func TestSample05(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("sample05.txt")
	assertWithoutBonus(t, lines, 20, 937, 18740)
	assertWithBonus(t, lines, 31, 30, 38, 1140)
}

func TestPuzzle(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("puzzle.txt")
	assertWithoutBonus(t, lines, 82, 2624, 215168)
	assertWithBonus(t, lines, 13, 0, 0, 52374)
}

func TestPuzzle2(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("puzzle2.txt")
	assertWithoutBonus(t, lines, 80, 2444, 195520)
	assertWithBonus(t, lines, 42, 0, 0, 52374)
}

func TestPuzzle4(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("puzzle4.txt")
	assertWithoutBonus(t, lines, 67, 2843, 190012)
	assertWithBonus(t, lines, 26, 0, 0, 34364)
}

// https://github.com/ShaneMcC/aoc-2018/blob/master/15/tests/

func TestShaneMoveLeft(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("test_shane_moveLeft.txt")
	assertWithoutBonus(t, lines, 34, 295, 10030)
	assertWithBonus(t, lines, 6, 46, 5, 230)
}

func TestShaneMoveRight(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("test_shane_moveRight.txt")
	assertWithoutBonus(t, lines, 34, 301, 10234)
	assertWithBonus(t, lines, 7, 41, 23, 943)
}

func TestShaneMovement(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("test_shane_movement.txt")
	assertWithoutBonus(t, lines, 18, 1546, 27828)
	assertWithBonus(t, lines, 97, 16, 83, 1328)
}

func TestShaneWall(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("test_shane_wall.txt")
	assertWithoutBonus(t, lines, 38, 486, 18468)
	assertWithBonus(t, lines, 8, 60, 2, 120)
}

func TestShanePuzzle(t *testing.T) {
	lines, _ := dayless.ReadFileToArray("test_shane_puzzle.txt")
	assertWithoutBonus(t, lines, 70, 2781, 197025)
	assertWithBonus(t, lines, 26, 0, 0, 44423)
}
