package main

import (
	"de.knallisworld/aoc/aoc2018/day24/lib"
	"de.knallisworld/aoc/aoc2018/dayless"
	"fmt"
	"sort"
	"time"
)

const AocDay = 24
const AocDayName = "day24"
const debug = false
const trace = false

var allGroupTypes = []lib.GroupType{lib.ImmuneSystemGroup, lib.InfectionGroup}

func main() {
	dayless.PrintDayHeader(AocDay)
	defer dayless.TimeTrack(time.Now(), AocDayName)

	lines, _ := dayless.ReadFileToString(AocDayName + "/puzzle.txt")
	master := lib.NewParser().Parse(*lines)

	dayless.PrintStepHeader(1)
	result := fight(master.Clone(), 0)
	if result[lib.ImmuneSystemGroup] > 0 {
		fmt.Printf("🎉 The winning army '%s' have %d units\n", lib.ImmuneSystemGroup.ToString(), result[lib.ImmuneSystemGroup])
	} else {
		fmt.Printf("🎉 The winning army '%s' have %d units\n", lib.InfectionGroup.ToString(), result[lib.InfectionGroup])
	}
	fmt.Println()

	dayless.PrintStepHeader(2)
	bonus := binarySearch(1, 1000000000000, func(bonus int) bool {
		result := fight(master.Clone(), bonus)
		if left, exist := result[lib.ImmuneSystemGroup]; exist && left > 0 {
			return true
		}
		return false
	})
	fmt.Printf("🎉 The winning army '%s' have %d units with bonus=%d\n",
		lib.ImmuneSystemGroup.ToString(),
		fight(master.Clone(), bonus)[lib.ImmuneSystemGroup],
		bonus,
	)
	fmt.Println()
}

func binarySearch(min int, max int, worker func(i int) bool) int {
	b := 1
	for {
		if b*2 < min {
			b *= 2
		} else {
			break
		}
	}
	for b < max-min {
		b *= 2
	}

	for {
		// fmt.Printf("Running for b=%d [min=%d, max=%d]\n", b, min, max)
		for i := min; i <= max; i += b {
			if worker(i) {
				// fmt.Printf("Y for i=%d\n", i)
				min = i - b
				max = i
				if b == 1 {
					return i
				}
				break
			}
		}

		if b == 1 {
			return -1
		} else {
			b /= 2
		}
	}
}

func fight(groups lib.Groups, boost int) map[lib.GroupType]int {

	if boost > 0 {
		for i := range groups {
			group := &groups[i]
			if group.GroupType == lib.ImmuneSystemGroup {
				group.AttackDamage += boost
			}
		}
	}

	if trace {
		for _, g := range groups {
			fmt.Println(g.ToString())
		}
	}

	for {

		if debug {
			fmt.Println("🤺 Fight!")
		}

		// Presenting
		combatEnd := false
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Reference() < groups[j].Reference()
		})
		for _, groupType := range allGroupTypes {
			if debug {
				fmt.Printf("%s:\n", groupType.ToString())
			}
			totalUnits := 0
			for _, group := range groups {
				if group.GroupType == groupType && group.Units > 0 {
					if debug {
						fmt.Printf("Group %d contains %d units [initiative=%d, attack=%s,%d immune=%s, weak=%s]\n", group.Id, group.Units, group.Initiative, group.AttackType.ToString(), group.AttackDamage, group.ImmuneTypes, group.WeakTypes)
					}
					totalUnits += group.Units
				}
			}
			if totalUnits == 0 {
				combatEnd = true
				if debug {
					fmt.Println("✋ No groups remain.")
				}
			}
		}
		if debug {
			fmt.Println()
		}

		if combatEnd {
			break
		}

		// Target selection
		defendersInCombat := make(map[string]string)
		sort.Slice(groups, func(i, j int) bool {
			group1Power := groups[i].EffectivePower()
			group2Power := groups[j].EffectivePower()
			if group1Power != group2Power {
				return group1Power > group2Power
			} else {
				return groups[i].Initiative > groups[j].Initiative
			}
		})
		for idxAttacker := range groups {
			attacker := &groups[idxAttacker]
			// still alive?
			if attacker.Units < 1 {
				continue
			}
			candidates := make([]*lib.Group, 0)
			for idxCandidate := range groups {
				candidate := &groups[idxCandidate]
				// still alive?
				if candidate.Units < 1 {
					continue
				}
				// already in combat?
				if _, exist := defendersInCombat[candidate.Reference()]; exist {
					continue
				}
				if candidate != attacker && candidate.GroupType != attacker.GroupType && candidate.EffectiveDamage(attacker.AttackType, attacker.EffectivePower()) > 0 {
					candidates = append(candidates, candidate)
					if debug {
						fmt.Printf(
							"%s group %d would deal defending group %d %d damage\n",
							attacker.GroupType.ToString(),
							attacker.Id,
							candidate.Id,
							candidate.EffectiveDamage(attacker.AttackType, attacker.EffectivePower()),
						)
					}
				}
			}
			sort.Slice(candidates, func(i, j int) bool {
				candidate1EffectiveDamage := candidates[i].EffectiveDamage(attacker.AttackType, attacker.EffectivePower())
				candidate2EffectiveDamage := candidates[j].EffectiveDamage(attacker.AttackType, attacker.EffectivePower())
				if candidate1EffectiveDamage != candidate2EffectiveDamage {
					return candidate1EffectiveDamage > candidate2EffectiveDamage
				} else {
					candidate1EffectivePower := candidates[i].EffectivePower()
					candidate2EffectivePower := candidates[j].EffectivePower()
					if candidate1EffectivePower != candidate2EffectivePower {
						return candidate1EffectivePower > candidate2EffectivePower
					} else {
						candidate1Initiative := candidates[i].Initiative
						candidate2Initiative := candidates[j].Initiative
						return candidate1Initiative > candidate2Initiative
					}
				}
			})
			defenderChosen := false
			for idxCandidate := range candidates {
				candidate := candidates[idxCandidate]
				if !defenderChosen {
					if _, exist := defendersInCombat[candidate.Reference()]; !exist {
						defendersInCombat[candidate.Reference()] = attacker.Reference()
						defenderChosen = true
					}
				}
			}
		}
		if debug {
			fmt.Println()
		}

		// Attacking
		atLeastOneAttack := false
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Initiative > groups[j].Initiative
		})
		for idxAttacker := range groups {
			attacker := &groups[idxAttacker]
			// still alive?
			if attacker.Units < 1 {
				continue
			}
			for defenderRef, attackerRef := range defendersInCombat {
				defender := groups.FindByReference(defenderRef)
				if attackerRef == attacker.Reference() {
					damage := defender.EffectiveDamage(attacker.AttackType, attacker.EffectivePower())
					killedUnits := damage / defender.HitPointsPerUnit
					defender.Units -= killedUnits
					if defender.Units < 0 {
						defender.Units = 0
					}
					if debug {
						fmt.Printf("⚔️ %s group %d attacks defending group %d, killing %d units\n", attacker.GroupType.ToString(), attacker.Id, defender.Id, killedUnits)
					}
					if killedUnits > 0 {
						atLeastOneAttack = true
					}
					break
				}
			}
		}
		if debug {
			fmt.Println()
		}

		if !atLeastOneAttack {
			fmt.Printf("🤕 Stucked; no units killed. Try next one.\n")
			return nil
		}
	}

	resultMap := make(map[lib.GroupType]int)
	for _, groupType := range allGroupTypes {
		resultMap[groupType] = groups.CountUnitsForGroup(groupType)
	}

	return resultMap
}
