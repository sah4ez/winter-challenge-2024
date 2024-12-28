package main

import (
	"math"
)

type Position struct {
	X int
	Y int
}

func (p Position) Equal(pos Position) bool {
	return p.X == pos.X && p.Y == pos.Y
}

func (p Position) Up() Position {
	return Position{X: p.X, Y: p.Y - 1}
}

func (p Position) Down() Position {
	return Position{X: p.X, Y: p.Y + 1}
}

func (p Position) Left() Position {
	return Position{X: p.X - 1, Y: p.Y}
}

func (p Position) Right() Position {
	return Position{X: p.X + 1, Y: p.Y}
}

func (p Position) GetLocality() []Position {
	return []Position{
		p.Up(),
		p.Down(),
		p.Left(),
		p.Right(),
		Position{X: p.X + 1, Y: p.Y + 1},
		Position{X: p.X - 1, Y: p.Y - 1},
		Position{X: p.X + 1, Y: p.Y - 1},
		Position{X: p.X - 1, Y: p.Y + 1},
	}
}

func (p Position) GetRoseLocality() []Position {
	return []Position{
		p.Up(),
		p.Down(),
		p.Left(),
		p.Right(),
	}
}

func (from Position) Distance(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}
