package player

import (
	"de.knallisworld/aoc/aoc2018/day15/grid"
	"fmt"
)

type Player struct {
	Group       rune
	AttackPower int
	HitPoints   int
	X           int
	Y           int
}

func New(pType rune, x int, y int) Player {
	return NewWithDetails(pType, 3, 200, x, y)
}

func NewWithDetails(group rune, attackPower int, hitPoints int, x int, y int) Player {
	return Player{
		Group:       group,
		AttackPower: attackPower,
		HitPoints:   hitPoints,
		X:           x,
		Y:           y,
	}
}

func (p *Player) Point() grid.Point {
	return grid.Point{X: p.X, Y: p.Y}
}

func (p *Player) IsAlive() bool {
	return p.HitPoints > 0
}

func (p *Player) Attack(target *Player) {
	target.HitPoints -= p.AttackPower
}

func (p Player) ToString() string {
	return fmt.Sprintf("%c[%s,HP=%d,AP=%d]", p.Group, p.Point().ToString(), p.HitPoints, p.AttackPower)
}

type Players []*Player

func (p Players) ToString() string {
	result := ""
	for _, p := range p {
		if len(result) > 0 {
			result += "\n"
		}
		result += p.ToString()
	}
	return result
}

func (p Players) Do(f func(player *Player)) {
	for i := range p {
		f(p[i])
	}
}

func (p Players) IsEmpty(f func(player *Player) bool) bool {
	return len(p.Filter(f)) == 0
}

func (p Players) Filter(f func(player *Player) bool) (result Players) {
	for i := range p {
		if f(p[i]) {
			result = append(result, p[i])
		}
	}
	return result
}

func (p Players) Sum(f func(player *Player) int) (result int) {
	for i := range p {
		result += f(p[i])
	}
	return result
}

func (p Players) GetByXY(x int, y int) *Player {
	for i := range p {
		p := p[i]
		if p.X == x && p.Y == y {
			return p
		}
	}
	return nil
}
