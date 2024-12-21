package main

import "fmt"

type Game struct {
	Width  int
	Height int
	state  *State
}

func (g *Game) Scan() {
	// width: columns in the game grid
	// height: rows in the game grid
	fmt.Scan(&g.Width, &g.Height)
}

func (g *Game) State() *State {
	if g.state == nil {
		g.state = NewState(g.Width, g.Height)
	} else {
		g.state.Scan()
	}
	g.state.ScanEnties()
	DebugMsg("load enties")
	return g.state
}

func NewGame() *Game {
	g := &Game{}
	g.Scan()
	return g
}
