package main

import (
	"fmt"
	"math"
	"os"
)

const BasicType = "BASIC"
const GrowCmd = "GROW"
const WaitCmd = "WAIT"
const WallTypeEntity = "WALL"
const RootTypeEntity = "ROOT"
const BasicTypeEntity = "BASIC"
const HarvesterTypeEntity = "HARVESTER"
const DirW = "W"
const DirS = "S"
const DirN = "N"
const DirE = "E"
const AProteinTypeEntity = "A"

type Entity struct {
	Pos           Position
	Type          string
	Owner         int
	OrganID       int
	OrganDir      string
	OrganParentID int
	OrganRootID   int

	NextDistance float64
}

func (e *Entity) Scan() {
	fmt.Scan(&e.Pos.X, &e.Pos.Y, &e.Type, &e.Owner, &e.OrganID, &e.OrganDir, &e.OrganParentID, &e.OrganRootID)
}

func (e *Entity) GrowBasic() string {
	return fmt.Sprintf("%s %d %d %d %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, BasicTypeEntity)
}

func (e *Entity) GrowHarvester(direction string) string {
	return fmt.Sprintf("%s %d %d %d %s %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, HarvesterTypeEntity, direction)
}

func (e *Entity) IsWall() bool {
	return e.Type == WallTypeEntity
}

func (e *Entity) IsRoot() bool {
	return e.Type == RootTypeEntity
}

func (e *Entity) IsAProtein() bool {
	return e.Type == AProteinTypeEntity
}

func (e *Entity) IsProtein() bool {
	return e.IsAProtein()
}

func (e *Entity) IsBasic() bool {
	return e.Type == BasicTypeEntity
}

func (e *Entity) IsEmpty() bool {
	return e.Type == ""
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
	g.state = NewState(g.Width, g.Height)
	// if g.state == nil {
	// g.state = NewState(g.Width, g.Height)
	// } else {
	// g.state.Scan()
	// }
	g.state.ScanEnties()
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
	step := 0
	for {
		state := game.State()

		state.ScanStocks()
		state.ScanReqActions()

		state.DoAction()
		state.Debug()
		DebugMsg("step: ", step)
		step += 1
	}
}

type Position struct {
	X int
	Y int
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

func (from Position) Distance(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
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
	proteins    []*Entity

	nextEntity []*Entity
}

func (s *State) Scan() {
	fmt.Scan(&s.EntityCount)
}

func (s *State) ScanEnties() {
	s.entities = make([]*Entity, 0)
	s.myEntities = make([]*Entity, 0)
	s.oppEntities = make([]*Entity, 0)
	s.proteins = make([]*Entity, 0)
	s.nextEntity = make([]*Entity, 0)

	for i := 0; i < s.EntityCount; i++ {
		e := NewEntity()
		if e.Owner == 1 {
			s.myEntities = append(s.myEntities, e)
		} else if e.Owner == 0 {
			s.oppEntities = append(s.oppEntities, e)
		} else {
			s.entities = append(s.entities, e)
			if e.IsAProtein() {
				s.proteins = append(s.proteins, e)
			}
		}
		s.matrix[e.Pos.Y][e.Pos.X] = e
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

		s.walk(0, 0, s.Dummy)
		// fmt.Fprintln(os.Stderr, "Debug messages...")
		// fmt.Println("WAIT") // Write action to stdout
	}
	total := s.RequiredActionsCount
	for _, e := range s.nextEntity {
		if s.MyStock.A == 0 {
			break
		}
		if e.NextDistance <= 1.0 {
			if s.MyStock.C == 1 && s.MyStock.D == 1 && e.OrganDir != "" {
				fmt.Println(e.GrowHarvester(e.OrganDir))
			} else {
				fmt.Println(e.GrowBasic())
			}
		} else {
			fmt.Println(e.GrowBasic())
		}
		total = total - 1
		if total == 0 {
			break
		}
	}

	for i := 0; i < total; i++ {
		fmt.Println("WAIT") // Write action to stdout
	}
}

func (s *State) Debug() {
	DebugMsg("my", *s.MyStock)
	DebugMsg("opp", *s.OpponentStock)
	for i, k := range s.matrix {
		for j, n := range k {
			if n != nil {
				fmt.Fprintf(os.Stderr, " %c(%d;%d;%.2f) ", n.Type[0], i, j, n.NextDistance)
			} else {
				fmt.Fprintf(os.Stderr, " _(%d;%d) ", i, j)
			}
		}
		fmt.Fprint(os.Stderr, "\n")
	}
}

func (s *State) getByPos(p Position) (e *Entity) {
	if p.Y >= len(s.matrix) {
		return nil
	}
	row := s.matrix[p.Y]
	if p.X >= len(row) {
		return nil
	}
	return row[p.X]
}

func (s *State) GetFreePos() []*Entity {
	freePos := make([]*Entity, 0)
	do := func(e *Entity) {
		if e == nil {
			return
		}
		up := s.getByPos(e.Pos.Up())
		if up == nil || up.IsProtein() {
			up = &Entity{Pos: e.Pos.Up()}
			up.OrganID = e.OrganID
			freePos = append(freePos, up)
		}
		down := s.getByPos(e.Pos.Down())
		if down == nil || down.IsProtein() {
			down = &Entity{Pos: e.Pos.Down()}
			down.OrganID = e.OrganID
			freePos = append(freePos, down)
		}
		left := s.getByPos(e.Pos.Left())
		if left == nil || left.IsProtein() {
			left = &Entity{Pos: e.Pos.Left()}
			left.OrganID = e.OrganID
			freePos = append(freePos, left)
		}
		right := s.getByPos(e.Pos.Right())
		if right == nil || right.IsProtein() {
			right = &Entity{Pos: e.Pos.Right()}
			right.OrganID = e.OrganID
			freePos = append(freePos, right)
		}
	}

	for _, e := range s.myEntities {
		do(e)
	}

	return freePos
}

func (s *State) GetHarvesterDir(e *Entity) string {
	if e == nil {
		return ""
	}
	up := s.getByPos(e.Pos.Up())
	if up != nil && up.IsProtein() {
		return DirN
	}
	down := s.getByPos(e.Pos.Down())
	if down != nil && down.IsProtein() {
		return DirS
	}
	left := s.getByPos(e.Pos.Left())
	if left != nil && left.IsProtein() {
		return DirW
	}
	right := s.getByPos(e.Pos.Right())
	if right != nil && right.IsProtein() {
		return DirE
	}

	return ""
}

func NewState(h, w int) *State {
	s := &State{
		entities:    make([]*Entity, 0),
		myEntities:  make([]*Entity, 0),
		oppEntities: make([]*Entity, 0),
		proteins:    make([]*Entity, 0),
		nextEntity:  make([]*Entity, 0),
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

func (s *State) walk(x, y int, fn func(e *Entity) bool) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= len(s.matrix) || len(s.matrix) == 0 {
		return
	}
	if y >= len(s.matrix[0]) {
		return
	}

	for _, k := range s.matrix[x:] {
		if k == nil {
			continue
		}
		for _, n := range k[y:] {
			if !fn(n) {
				break
			}
		}
	}
}

func (s *State) Dummy(e *Entity) bool {

	freePos := s.GetFreePos()
	if len(freePos) == 0 {
		return false
	}

	for _, protein := range s.proteins {
		for i, free := range freePos {
			newDistance := free.Pos.Distance(protein.Pos)
			if math.Abs(newDistance) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && newDistance >= 0) {
				free.NextDistance = math.Abs(newDistance)
				freePos[i] = free
			}
			// DebugMsg("dist", newDistance, free.NextDistance)
		}
	}

	min := freePos[0]
	for _, free := range freePos {
		s.matrix[free.Pos.Y][free.Pos.X] = &Entity{Type: "F", Pos: free.Pos, NextDistance: free.NextDistance}
		if min.NextDistance >= free.NextDistance {
			min = free
			if min.NextDistance == 1 {
				min.OrganDir = s.GetHarvesterDir(min)
			}
		}
	}

	s.nextEntity = append(s.nextEntity, min)

	return false
}
