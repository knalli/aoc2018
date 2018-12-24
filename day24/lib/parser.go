package lib

import (
	"de.knallisworld/aoc/aoc2018/dayless"
	"regexp"
	"strings"
)

// Immune System:
// 17 units each with 5390 hit points (weak to radiation, bludgeoning) with
//  an attack that does 4507 fire damage at initiative 2
// 989 units each with 1274 hit points (immune to fire; weak to bludgeoning,
//  slashing) with an attack that does 25 slashing damage at initiative 3

type GroupParser struct {
	pattern *regexp.Regexp
}

func NewParser() GroupParser {
	return GroupParser{
		pattern: regexp.MustCompile("(\\d+) units\\s+each\\s+with\\s+(\\d+)\\s+hit\\s+points(\\s+\\((.*)\\))?\\s+with\\s+an\\s+attack\\s+that\\s+does\\s+(\\d+)\\s+(cold|fire|bludgeoning|radiation|slashing)\\s+damage\\s+at\\s+initiative (\\d+)"),
	}
}

func (p GroupParser) Parse(lines string) (result Groups) {
	// split by groups
	for _, part := range strings.Split(lines, "\n\n") {
		for _, g := range p.parseGroup(part) {
			result = append(result, g)
		}
	}
	return result
}

func (p GroupParser) parseGroup(str string) (result Groups) {

	groupType := ParseGroupType(str[0:strings.Index(str, ":")])

	id := 1
	for _, m := range p.pattern.FindAllStringSubmatch(str, 1000) {
		immuneTypes := make([]AttackType, 0)
		weakTypes := make([]AttackType, 0)
		for _, modStr := range strings.Split(m[4], ";") {
			if strings.Contains(modStr, "immune to") {
				immuneTypes = p.parseAttackTypes(modStr)
			}
			if strings.Contains(modStr, "weak to") {
				weakTypes = p.parseAttackTypes(modStr)
			}
		}
		result = append(result, Group{
			GroupType:        groupType,
			Id:               id,
			Units:            dayless.ParseInt(m[1]),
			HitPointsPerUnit: dayless.ParseInt(m[2]),
			ImmuneTypes:      immuneTypes,
			WeakTypes:        weakTypes,
			AttackDamage:     dayless.ParseInt(m[5]),
			AttackType:       ParseAttackType(m[6]),
			Initiative:       dayless.ParseInt(m[7]),
		})
		id++
	}

	return result
}

func (p GroupParser) parseAttackTypes(str string) (result []AttackType) {
	for _, t := range []AttackType{BludgeoningAttack, ColdAttack, FireAttack, RadiationAttack, SlashingAttack} {
		if strings.Contains(str, t.ToString()) {
			result = append(result, t)
		}
	}
	return result
}
