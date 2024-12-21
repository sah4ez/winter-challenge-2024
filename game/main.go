package main

import (
	"fmt"
	"os"
)

const BasicType = "BASIC"
const GrowCmd = "GROW"
const WaitCmd = "WAIT"
const WallTypeEntity = "WALL"
const RootTypeEntity = "ROOT"
const BasicTypeEntity = "BASIC"
const AProteinTypeEntity = "A"

type Entity struct {
	X             int
	Y             int
	Type          string
	Owner         int
	OrganID       int
	OrganDir      string
	OrganParentID int
	OrganRootID   int
}

func (e *Entity) Scan() {
	fmt.Scan(&e.X, &e.Y, &e.Type, &e.Owner, &e.OrganID, &e.OrganDir, &e.OrganParentID, &e.OrganRootID)
}

func (e *Entity) Grow(x, y int, typeOrgan string) string {
	return fmt.Sprintf("%s %d %d %s", GrowCmd, x, y, typeOrgan)
}

func NewEntity() *Entity {

	e := &Entity{}
	e.Scan()
	return e
}

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

/**
 * Grow and multiply your organisms to end up larger than your opponent.
 **/

func main() {
	game := NewGame()
	for {
		state := game.State()

		state.ScanStocks()
		state.ScanReqActions()

		state.Debug()
		state.DoAction()
	}
}

type State struct {
	EntityCount          int
	RequiredActionsCount int

	MyStock       *Stock
	OpponentStock *Stock

	matrix [][]*Entity

	myEntities  []*Entity
	oppEntities []*Entity
	entities    []*Entity
}

func (s *State) Scan() {
	fmt.Scan(&s.EntityCount)
}

func (s *State) ScanEnties() {
	s.entities = make([]*Entity, 0)
	s.myEntities = make([]*Entity, 0)
	s.oppEntities = make([]*Entity, 0)

	for i := 0; i < s.EntityCount; i++ {
		e := NewEntity()
		if e.Owner == 1 {
			s.myEntities = append(s.myEntities, e)
		} else if e.Owner == 0 {
			s.oppEntities = append(s.oppEntities, e)
		} else {
			s.entities = append(s.entities, e)
		}
		s.matrix[e.Y][e.X] = e
	}
}

func (s *State) ScanStocks() {
	s.MyStock = NewStock()
	s.OpponentStock = NewStock()
}

func (s *State) ScanReqActions() {
	// requiredActionsCount: your number of organisms, output an action for each one in any order
	fmt.Scan(&s.RequiredActionsCount)
}

func (s *State) DoAction() {
	for i := 0; i < s.RequiredActionsCount; i++ {

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Println("WAIT") // Write action to stdout
	}
}

func (s *State) Debug() {
	for i, k := range s.matrix {
		for j, n := range k {
			if n != nil {
				fmt.Fprintf(os.Stderr, " %c(%d, %d) ", n.Type[0], i, j)
			} else {
				fmt.Fprintf(os.Stderr, " _(%d, %d) ", i, j)
			}
		}
		fmt.Fprint(os.Stderr, "\n")
	}
}

func NewState(h, w int) *State {
	s := &State{
		entities:    make([]*Entity, 0),
		myEntities:  make([]*Entity, 0),
		oppEntities: make([]*Entity, 0),
		matrix:      make([][]*Entity, 0),
	}

	for i := 0; i < w; i++ {
		s.matrix = append(s.matrix, make([]*Entity, h))
	}

	s.Scan()
	return s
}

type Stock struct {
	A int
	B int
	C int
	D int
}

func (s *Stock) Scan() {
	fmt.Scan(&s.A, &s.B, &s.C, &s.D)
}

func NewStock() *Stock {
	s := &Stock{}
	s.Scan()
	return s
}

func DebugMsg(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
}
