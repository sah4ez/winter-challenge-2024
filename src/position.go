package main

import (
	"fmt"
	"math"
)

type Position struct {
	X int
	Y int

	Parent *Position
}

func (p Position) ID() string {
	return fmt.Sprintf("(%d:%d)", p.X, p.Y)
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

func (from Position) EucleadDistance(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}

func (p Position) ToCoordinates() Coordinates {
	return Coordinates{
		float64(p.X),
		float64(p.Y),
	}
}

func (p Position) ToLog() string {
	depth := 0
	parent := p.Parent
	for parent != nil {
		parent = parent.Parent
		depth += 1
	}
	return fmt.Sprintf("(%d:%d)%d", p.X, p.Y, depth)
}

func FromCoordinates(c Coordinates) Position {
	if len(c) != 2 {
		panic("invalid coordinates")
	}
	return Position{X: int(c[0]), Y: int(c[1])}
}

func (s *State) InMatrix(p Position) bool {
	out := p.X < 0 || p.Y < 0 || p.Y >= s.w || p.X >= s.h
	return !out
}

func NewPos(x, y int) Position {
	return Position{X: x, Y: y}
}
