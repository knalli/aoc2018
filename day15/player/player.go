package player

import (
	"de.knallisworld/aoc/aoc2018/day15/grid"
	"fmt"
)

type Player struct {
	PType       int32
	AttackPower int
	HitPoints   int
	X           int
	Y           int
}

type Players struct {
	values []Player
}

func New(pType int32, x int, y int) Player {
	return NewWithDetails(pType, 3, 200, x, y)
}

func NewWithDetails(pType int32, attackPower int, hitPoints int, x int, y int) Player {
	return Player{
		PType:       pType,
		AttackPower: attackPower,
		HitPoints:   hitPoints,
		X:           x,
		Y:           y,
	}
}

func (p Player) Point() grid.Point {
	return grid.Point{X: p.X, Y: p.Y}
}

func (p Player) IsAlive() bool {
	return p.HitPoints > 0
}

func (p Player) Attack(target *Player) {
	target.HitPoints -= p.AttackPower
}

func NewPlayers(players []Player) Players {
	return Players{values: players}
}

func (p Player) ToString() string {
	return fmt.Sprintf("%c[%s,HP=%d]", p.PType, p.Point().ToString(), p.HitPoints)
}

func (p Players) ToString() string {
	result := ""
	for _, p := range p.values {
		if len(result) > 0 {
			result += "\n"
		}
		result += p.ToString()
	}
	return result
}

func (p Players) Do(f func(player Player)) {
	for _, player := range p.values {
		f(player)
	}
}

func (p Players) GetByXY(x int, y int) *Player {
	for i := range p.values {
		p := &p.values[i]
		if p.X == x && p.Y == y {
			return p
		}
	}
	return nil
}
