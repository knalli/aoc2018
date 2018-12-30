package main

import (
	lGrid "de.knallisworld/aoc/aoc2018/day15/grid"
	lPlayer "de.knallisworld/aoc/aoc2018/day15/player"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"github.com/golang-collections/collections/set"
	"math"
	"sort"
	"time"
)

const AocDay = 15
const AocDayName = "day15"
const trace = true

const SYM_ELF = 'E'
const SYM_GOBLIN = 'G'
const SYM_WALL = '#'
const SYM_CARVE = '.'

func traceln(a ...interface{}) {
	if trace {
		fmt.Println(a...)
	}
}

func tracef(format string, a ...interface{}) {
	if trace {
		fmt.Printf(format, a...)
	}
}

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	dayless.PrintStepHeader(1)
	lines, _ := dayless.ReadFileToArray(AocDayName + "/puzzle.txt")
	fullRounds, hitPointsSum := part1(lines)
	fmt.Printf("🎉 The outcome of this battle is: %d * %d = %d\n", fullRounds, hitPointsSum, fullRounds*hitPointsSum)
	fmt.Println()

	dayless.PrintStepHeader(2)
	bonus, fullRounds, hitPointsSum := part2(lines)
	fmt.Printf("🎉 With an additional bonus for the elves of %d,  the outcome of this battle is: %d * %d = %d\n", bonus, fullRounds, hitPointsSum, fullRounds*hitPointsSum)
	fmt.Println()
}

func part1(lines []string) (fullRounds int, hitPointsSum int) {
	grid, players := initializeGame(lines)
	tracef("Initially: \n%s\n", grid.ToStringWithBorder())
	for fullRounds = 1; ; fullRounds++ {
		traceln("=====================")
		tracef("🎰 Round %d...\n", fullRounds)
		traceln("=====================")
		stillRunning := playRound(&grid, players)
		if !stillRunning {
			tracef("👉 Round %d has been aborted (so not counting as 'full')\n", fullRounds)
			fullRounds-- // last round was not ended
			break
		}
		tracef("\n  👉 After %d rounds:\n%s\n%s\n\n", fullRounds, grid.ToStringWithBorder(), players.ToString())
	}
	hitPointsSum = players.Filter(func(player *lPlayer.Player) bool {
		return player.IsAlive()
	}).Sum(func(player *lPlayer.Player) int {
		return player.HitPoints
	})
	fmt.Printf("Combat ends after %d full rounds\n", fullRounds)
	if players.IsEmpty(func(player *lPlayer.Player) bool {
		return player.Group == SYM_ELF && player.IsAlive()
	}) {
		fmt.Printf("Goblins win with %d total hit points left\n", hitPointsSum)
	} else {
		fmt.Printf("Elves win with %d total hit points left\n", hitPointsSum)
	}
	return fullRounds, hitPointsSum
}

func part2(lines []string) (minAttackBonus int, fullRounds int, hitPointsSum int) {

	minAttackBonus = dayless.BinarySearch(1, math.MaxInt32, func(i int) bool {
		grid, players := initializeGame(lines)
		tracef("😀 Buffing all elves with additional attack power of %d\n", i)
		// Buff elves
		players.Do(func(player *lPlayer.Player) {
			if player.Group == SYM_ELF {
				player.AttackPower = player.AttackPower + i
			}
		})
		tracef("Initially: \n%s\n", grid.ToStringWithBorder())
		for fullRounds = 1; ; fullRounds++ {
			traceln("=====================")
			tracef("🎰 Round %d...\n", fullRounds)
			traceln("=====================")
			stillRunning := playRound(&grid, players)
			if !stillRunning {
				tracef("👉 Round %d has been aborted (so not counting as 'full')\n", fullRounds)
				fullRounds-- // last round was not ended
				break
			}
			tracef("\n  👉 After %d rounds:\n%s\n%s\n\n", fullRounds, grid.ToStringWithBorder(), players.ToString())

			if !players.IsEmpty(func(p *lPlayer.Player) bool {
				return p.Group == SYM_ELF && !p.IsAlive()
			}) {
				return false
			}
		}
		return players.IsEmpty(func(p *lPlayer.Player) bool {
			return p.Group == SYM_ELF && !p.IsAlive()
		})
	}, trace)
	tracef("  👉 Found additional attack power bonus = %d is required Elves will win\n", minAttackBonus)

	grid, players := initializeGame(lines)
	// Buff elves
	players.Do(func(player *lPlayer.Player) {
		if player.Group == SYM_ELF {
			player.AttackPower = player.AttackPower + minAttackBonus
		}
	})
	tracef("Initially: \n%s\n", grid.ToStringWithBorder())
	for fullRounds = 1; ; fullRounds++ {
		traceln("=====================")
		tracef("🎰 Round %d...\n", fullRounds)
		traceln("=====================")
		stillRunning := playRound(&grid, players)
		if !stillRunning {
			tracef("👉 Round %d has been aborted (so not counting as 'full')\n", fullRounds)
			fullRounds-- // last round was not ended
			break
		}
		tracef("\n  👉 After %d rounds:\n%s\n%s\n\n", fullRounds, grid.ToStringWithBorder(), players.ToString())
	}
	hitPointsSum = players.Filter(func(player *lPlayer.Player) bool {
		return player.IsAlive()
	}).Sum(func(player *lPlayer.Player) int {
		return player.HitPoints
	})
	fmt.Printf("Combat ends after %d full rounds\n", fullRounds)
	if players.IsEmpty(func(player *lPlayer.Player) bool {
		return player.Group == SYM_ELF && player.IsAlive()
	}) {
		fmt.Printf("Goblins win with %d total hit points left\n", hitPointsSum)
	} else {
		fmt.Printf("Elves win with %d total hit points left\n", hitPointsSum)
	}
	return minAttackBonus, fullRounds, hitPointsSum
}

func part2b(lines []string) (minAttackBonus int, fullRounds int, hitPointsSum int) {

	for minAttackBonus = 1; ; minAttackBonus++ {
		grid, players := initializeGame(lines)
		// Buff elves
		players.Do(func(player *lPlayer.Player) {
			if player.Group == SYM_ELF {
				player.AttackPower = player.AttackPower + minAttackBonus
			}
		})
		tracef("Initially: \n%s\n", grid.ToStringWithBorder())
		aborted := false
		for fullRounds = 1; ; fullRounds++ {
			traceln("=====================")
			tracef("🎰 Round %d...\n", fullRounds)
			traceln("=====================")
			stillRunning := playRound(&grid, players)
			if !stillRunning {
				tracef("👉 Round %d has been aborted (so not counting as 'full')\n", fullRounds)
				fullRounds-- // last round was not ended
				break
			}
			tracef("\n  👉 After %d rounds:\n%s\n%s\n\n", fullRounds, grid.ToStringWithBorder(), players.ToString())

			if !players.IsEmpty(func(p *lPlayer.Player) bool {
				return p.Group == SYM_ELF && !p.IsAlive()
			}) {
				aborted = true
				break
			}
		}
		if !aborted && players.IsEmpty(func(p *lPlayer.Player) bool {
			return p.Group == SYM_ELF && !p.IsAlive()
		}) {
			break // found
		}
	}
	tracef("  👉 Found additional attack power bonus = %d is required Elves will win\n", minAttackBonus)

	grid, players := initializeGame(lines)
	// Buff elves
	players.Do(func(player *lPlayer.Player) {
		if player.Group == SYM_ELF {
			player.AttackPower = player.AttackPower + minAttackBonus
		}
	})
	tracef("Initially: \n%s\n", grid.ToStringWithBorder())
	for fullRounds = 1; ; fullRounds++ {
		traceln("=====================")
		tracef("🎰 Round %d...\n", fullRounds)
		traceln("=====================")
		stillRunning := playRound(&grid, players)
		if !stillRunning {
			tracef("👉 Round %d has been aborted (so not counting as 'full')\n", fullRounds)
			fullRounds-- // last round was not ended
			break
		}
		tracef("\n  👉 After %d rounds:\n%s\n%s\n\n", fullRounds, grid.ToStringWithBorder(), players.ToString())
	}
	hitPointsSum = players.Filter(func(player *lPlayer.Player) bool {
		return player.IsAlive()
	}).Sum(func(player *lPlayer.Player) int {
		return player.HitPoints
	})
	return minAttackBonus, fullRounds, hitPointsSum
}

func initializeGame(lines []string) (lGrid.Grid, lPlayer.Players) {
	grid := lGrid.NewByStrings(lines)
	players := make(lPlayer.Players, 0)
	grid.Walk(func(p lGrid.Point, value rune) {
		switch value {
		case SYM_ELF, SYM_GOBLIN:
			player := lPlayer.New(value, p.X, p.Y)
			players = append(players, &player)
		}
	})
	return grid, players
}

func playRound(grid *lGrid.Grid, players lPlayer.Players) bool {

	tracef("👉 There are %d active players on the field\n", len(players.Filter(func(player *lPlayer.Player) bool {
		return player.IsAlive()
	})))

	// "units take their turns within a round is the reading order of their starting positions in that round, regardless
	// of the type of unit or whether other units have moved after the round started"
	sort.Slice(players, func(i, j int) bool {
		return players[i].Point().GetGridPosition(grid) < players[j].Point().GetGridPosition(grid)
	})

	for _, unit := range players {

		if !unit.IsAlive() {
			// dead units don't play anymore
			continue
		}

		tracef("🤖 %s is playing...\n", unit.ToString())

		// Identify all possible targets
		availableTargets := searchPossibleTargets(unit, grid, players)
		// no targets left -> round ends -> combat ends
		if len(availableTargets) == 0 {
			traceln("🖐 Combat will be end (no available targets left within round)")
			return false
		}

		// Do only move if no one is near...
		if len(resolveNearTargets(grid, unit, players)) == 0 {

			// Identify all open squares that are in range of each target (technical the target to move)
			availablePOIs := resolveAllAdjacentPointsByType(availableTargets, grid, SYM_CARVE)
			tracef("  👉️ Available POIs: Length = %d, Set = [%s]\n", len(availablePOIs), availablePOIs.ToString())
			if len(availablePOIs) == 0 {
				continue // next unit
			}

			var nextPoint *lGrid.Point
			for _, p := range availablePOIs {
				if grid.IsAdjacent(unit.Point(), p) {
					// Simple case: One of the available is actually an adjacent of the unit; go for it.
					nextPoint = &p
					break
				}
			}

			if nextPoint == nil {

				// Okay: Not so easy, we have to search for the next one...

				path := grid.GetShortestPathMulti(
					unit.Point(),
					availablePOIs,
					func(value rune) bool {
						return value == SYM_CARVE
					},
				)
				if len(path) == 0 {
					tracef("  🚫 No path found for %s\n", unit.ToString())
					continue // next one
				}

				// path = start -> next -> … -> chosen/goal (not target)
				goalPoint := path[len(path)-1]
				nextPoint = &path[1]
				if found := grid.ForEachAdjacentBreakable(unit.Point().X, unit.Point().Y, func(unitAdjacentPoint lGrid.Point, value rune) bool {
					if value != SYM_CARVE {
						return true // next
					}
					localPath := grid.GetShortestPathMulti(
						unitAdjacentPoint,
						[]lGrid.Point{goalPoint},
						func(value rune) bool {
							return value == SYM_CARVE
						},
					)
					if len(localPath) == len(path)-1 {
						// if the path is one lower than the actual shortest path, we have found the (first) field of the shortest path
						return false
					}
					return true // next
				}); found != nil {
					nextPoint = found
				}
			}

			// Move player
			move(unit, *nextPoint, grid)
		}

		// Check if already near a target
		nearTargets := resolveNearTargets(grid, unit, players)
		// ... and attack?
		if len(nearTargets) > 0 {
			target := nearTargets[0]
			attack(unit, target, grid)
		}
	}

	return true // next round
}

func resolveAllAdjacentPointsByType(targets lPlayer.Players, grid *lGrid.Grid, needle rune) (result lGrid.Points) {
	s := set.New()
	for _, target := range targets {
		grid.ForEachAdjacent(target.X, target.Y, func(p lGrid.Point, value rune) {
			if value == needle {
				s.Insert(p)
			}
		})
	}
	if s.Len() > 0 {
		s.Do(func(i interface{}) {
			p := i.(lGrid.Point)
			result = append(result, p)
		})
		// ensure set has not re-ordered the priority of points => "reading order"
		sort.Slice(result, func(i, j int) bool {
			return result[i].GetGridPosition(grid) < result[j].GetGridPosition(grid)
		})
	}
	return result
}

func resolveNearTargets(grid *lGrid.Grid, unit *lPlayer.Player, players lPlayer.Players) (result lPlayer.Players) {
	// Check if already near a target
	grid.ForEachAdjacent(unit.X, unit.Y, func(p lGrid.Point, value rune) {
		if (unit.Group == SYM_ELF && value == SYM_GOBLIN) || (unit.Group == SYM_GOBLIN && value == SYM_ELF) {
			target := players.GetByXY(p.X, p.Y)
			if target == nil {
				panic("failed assertion as player must exist at this coordinate")
			}
			if !target.IsAlive() {
				panic("oops, target is a ghost")
			}
			result = append(result, target)
		}
	})
	// "[…] the adjacent target with the fewest hit points is selected;
	// in a tie, the adjacent target with the fewest hit points which is first in reading order is selected."
	sort.Slice(result, func(i, j int) bool {
		player1 := result[i]
		player2 := result[j]
		if player1.HitPoints != player2.HitPoints {
			return player1.HitPoints < player2.HitPoints
		} else {
			return player1.Point().GetGridPosition(grid) < player2.Point().GetGridPosition(grid)
		}
	})
	return result
}

func searchPossibleTargets(unit *lPlayer.Player, grid *lGrid.Grid, players lPlayer.Players) (result lPlayer.Players) {
	grid.Walk(func(p lGrid.Point, value rune) {
		if (unit.Group == SYM_ELF && value == SYM_GOBLIN) || unit.Group == SYM_GOBLIN && value == SYM_ELF {
			// find in players
			target := players.GetByXY(p.X, p.Y)
			if target == nil {
				panic("assertion failed as player must exist at this coordinate")
			}
			if !target.IsAlive() {
				panic("oops, target is a ghost")
			}
			result = append(result, target)
		}
	})
	return result
}

func move(unit *lPlayer.Player, point lGrid.Point, grid *lGrid.Grid) {
	if !unit.IsAlive() {
		panic("oops, a ghost wants to moving")
	}
	tracef("  🚶️ %s is moving to %s\n", unit.ToString(), point.ToString())
	if grid.GetValueAtXY(point) != SYM_CARVE {
		panic("oops, the final move position is already taken")
	}
	grid.PutValueAtXY(unit.Point(), SYM_CARVE)
	grid.PutValueAtXY(point, unit.Group)
	unit.X = point.X
	unit.Y = point.Y
}

func attack(attacker *lPlayer.Player, target *lPlayer.Player, grid *lGrid.Grid) {
	if !attacker.IsAlive() {
		panic("oops, a ghost wants to be an attacker")
	}
	if !target.IsAlive() {
		panic("oops, a ghost wants to be an defender")
	}

	tracef("  ⚔️ %s is attacking %s\n", attacker.ToString(), target.ToString())
	attacker.Attack(target)

	if !target.IsAlive() {
		tracef("  ☠️ The defender %s has died. RIP\n", target.ToString())
		grid.PutValueAtXY(target.Point(), SYM_CARVE)
		// take out of game
		target.X = 0
		target.Y = 0
		target.HitPoints = 0
	}
}
