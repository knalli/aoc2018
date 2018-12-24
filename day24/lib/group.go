package lib

import (
	"fmt"
)

type AttackType string

func (at AttackType) ToString() string {
	switch at {
	case BludgeoningAttack:
		return "bludgeoning"
	case ColdAttack:
		return "cold"
	case FireAttack:
		return "fire"
	case RadiationAttack:
		return "radiation"
	case SlashingAttack:
		return "slashing"
	default:
		panic("found invalid AttackType")
	}
}

func ParseAttackType(str string) AttackType {
	switch str {
	case "bludgeoning":
		return BludgeoningAttack
	case "cold":
		return ColdAttack
	case "fire":
		return FireAttack
	case "radiation":
		return RadiationAttack
	case "slashing":
		return SlashingAttack
	default:
		panic("parse invalid AttackType")
	}
}

const (
	BludgeoningAttack AttackType = "bludgeoning"
	ColdAttack        AttackType = "cold"
	FireAttack        AttackType = "fire"
	RadiationAttack   AttackType = "radiation"
	SlashingAttack    AttackType = "slashing"
)

type GroupType string

const (
	ImmuneSystemGroup GroupType = "Immune System"
	InfectionGroup    GroupType = "Infection"
)

func (gt GroupType) ToString() string {
	switch gt {
	case ImmuneSystemGroup:
		return "Immune System"
	case InfectionGroup:
		return "Infection"
	default:
		panic("found invalid GroupType")
	}
}

func ParseGroupType(str string) GroupType {
	switch str {
	case "Immune System":
		return ImmuneSystemGroup
	case "Infection":
		return InfectionGroup
	default:
		panic("parse invalid GroupType")
	}
}

type Group struct {
	GroupType        GroupType
	Id               int
	Units            int
	HitPointsPerUnit int
	AttackDamage     int
	AttackType       AttackType
	WeakTypes        []AttackType
	ImmuneTypes      []AttackType
	Initiative       int
}

func (p Group) Reference() string {
	return fmt.Sprintf("%s_%d", p.GroupType.ToString(), p.Id)
}

func (p Group) EffectivePower() int {
	return p.Units * p.AttackDamage
}

func (p Group) EffectiveDamage(attackType AttackType, attackDamage int) int {
	if p.IsImmuneFor(attackType) {
		return 0
	} else if p.IsWeakFor(attackType) {
		return 2 * attackDamage
	} else {
		return attackDamage
	}
}

func (p Group) IsWeakFor(t AttackType) bool {
	return containsAttackType(p.WeakTypes, t)
}

func (p Group) IsImmuneFor(t AttackType) bool {
	return containsAttackType(p.ImmuneTypes, t)
}

func (p Group) weakTypesAsString() string {
	if len(p.WeakTypes) == 0 {
		return ""
	}
	str := ""
	for i, t := range p.WeakTypes {
		if i > 0 {
			str += ", "
		}
		str += t.ToString()
	}
	return fmt.Sprintf("weak to %s", str)
}

func (p Group) immuneTypesAsString() string {
	if len(p.ImmuneTypes) == 0 {
		return ""
	}
	str := ""
	for i, t := range p.ImmuneTypes {
		if i > 0 {
			str += ", "
		}
		str += t.ToString()
	}
	return fmt.Sprintf("immune to %s", str)
}

func (p Group) typesAsString() string {
	immune := p.immuneTypesAsString()
	weak := p.weakTypesAsString()
	if len(immune) > 0 && len(weak) > 0 {
		return fmt.Sprintf(" (%s; %s)", immune, weak)
	} else if len(immune) > 0 {
		return fmt.Sprintf(" (%s)", immune)
	} else if len(weak) > 0 {
		return fmt.Sprintf(" (%s)", weak)
	} else {
		return ""
	}
}

func (p Group) ToString() string {
	// 17 units each with 5390 hit points (weak to radiation, bludgeoning) with
	// an attack that does 4507 fire damage at initiative 2

	return fmt.Sprintf(
		"%s: %d units each with %d hit points%s with an attack that does %d %s damage at initiative %d",
		p.GroupType.ToString(),
		p.Units,
		p.HitPointsPerUnit,
		p.typesAsString(),
		p.AttackDamage,
		p.AttackType.ToString(),
		p.Initiative,
	)
}

func (g Group) Clone() Group {
	return Group{
		GroupType:        g.GroupType,
		Id:               g.Id,
		Units:            g.Units,
		HitPointsPerUnit: g.HitPointsPerUnit,
		ImmuneTypes:      g.ImmuneTypes,
		WeakTypes:        g.WeakTypes,
		AttackDamage:     g.AttackDamage,
		AttackType:       g.AttackType,
		Initiative:       g.Initiative,
	}
}

type Groups []Group

func (g Groups) FindByReference(reference string) *Group {
	for i := range g {
		if g[i].Reference() == reference {
			return &g[i]
		}
	}
	return nil
}

func (g Groups) CountUnitsExceptGroup(groupType GroupType) (total int) {
	for i := range g {
		if g[i].GroupType != groupType {
			total += g[i].Units
		}
	}
	return total
}

func (g Groups) CountUnitsForGroup(groupType GroupType) (total int) {
	for i := range g {
		if g[i].GroupType == groupType {
			total += g[i].Units
		}
	}
	return total
}

func (g Groups) Clone() Groups {
	result := make(Groups, 0)
	for _, group := range g {
		result = append(result, group.Clone())
	}
	return result
}

func containsAttackType(array []AttackType, search AttackType) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}
