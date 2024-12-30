package main

import "fmt"

type Game struct {
	Width  int
	Height int
	state  *State

	sporerFrom *Position
	sporerTo   *Position
}

func (g *Game) Scan() {
	// width: columns in the game grid
	// height: rows in the game grid
	fmt.Scan(&g.Width, &g.Height)
}

func (g *Game) State() *State {
	g.state = NewState(g.Width, g.Height)
	// if g.state == nil {
	// g.state = NewState(g.Width, g.Height)
	// } else {
	// g.state.Scan()
	// }
	g.state.ScanEnties()
	return g.state
}

func (g *Game) StartSporer(sporerFrom Position, sporerTo Position) {
	g.sporerFrom = &sporerFrom
	g.sporerTo = &sporerTo
}

func (g *Game) StopSporer() {
	g.sporerFrom = nil
	g.sporerTo = nil
}

func (g *Game) SporerPonits() (from, to Position) {
	return *g.sporerFrom, *g.sporerTo
}

func (g *Game) HasSporer() bool {
	return g.sporerFrom != nil && g.sporerTo != nil
}

func NewGame() *Game {
	g := &Game{}
	g.Scan()
	return g
}
