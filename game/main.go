package main

import (
	"fmt"
	"math"
	"os"
)

const BasicType = "BASIC"
const GrowCmd = "GROW"
const WaitCmd = "WAIT"
const SporeCmd = "SPORE"
const WallTypeEntity = "WALL"
const RootTypeEntity = "ROOT"
const BasicTypeEntity = "BASIC"
const FreeTypeEntity = "Free"
const HarvesterTypeEntity = "HARVESTER"
const TentacleTypeEntity = "TENTACLE"
const SporerTypeEntity = "SPORER"
const DirW = "W"
const DirS = "S"
const DirN = "N"
const DirE = "E"
const AProteinTypeEntity = "A"
const BProteinTypeEntity = "B"
const CProteinTypeEntity = "C"
const DProteinTypeEntity = "D"

func (s *State) Dummy(e *Entity) bool {

	if len(s.freePos) == 0 {
		return false
	}

	for _, protein := range s.proteins {
		for i, free := range s.freePos {
			newDistance := free.Pos.Distance(protein.Pos)
			if math.Abs(newDistance) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && newDistance >= 0) {
				free.NextDistance = math.Abs(newDistance)
				free.Protein = protein
				s.freePos[i] = free
			}
		}
	}

	min := s.freePos[0]
	for _, free := range s.freePos {
		s.matrix[free.Pos.Y][free.Pos.X] = &Entity{
			Type:         FreeTypeEntity,
			Pos:          free.Pos,
			NextDistance: free.NextDistance,
		}
		if min.NextDistance >= free.NextDistance {
			min = free
			min.OrganDir = s.GetHarvesterDir(min)
			if min.NextDistance == 1 {
				min.OrganDir = s.GetHarvesterDir(min)
			} else {
				min.OrganDir = s.GetSporerDir(min, min.Protein)
			}
		}
	}

	if _, ok := s.nextHash[min.ID()]; !ok {
		s.nextHash[min.ID()] = min
		s.nextEntity = append(s.nextEntity, min)
	}

	return false
}

type Entity struct {
	Pos           Position
	Type          string
	Owner         int
	OrganID       int
	OrganDir      string
	OrganParentID int
	OrganRootID   int

	NextDistance float64
	Score        float64
	Protein      *Entity
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

func (e *Entity) GrowTentacle(direction string) string {
	return fmt.Sprintf("%s %d %d %d %s %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, TentacleTypeEntity, direction)
}

func (e *Entity) GrowSporer(direction string) string {
	return fmt.Sprintf("%s %d %d %d %s %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, SporerTypeEntity, direction)
}

func (e *Entity) Spore() string {
	return fmt.Sprintf("%s %d %d %d", SporeCmd, e.OrganID, e.Protein.Pos.X, e.Protein.Pos.Y)
}

func (e *Entity) IsWall() bool {
	return e.Type == WallTypeEntity
}

func (e *Entity) IsRoot() bool {
	return e.Type == RootTypeEntity
}

func (e *Entity) IsSporer() bool {
	return e.Type == SporerTypeEntity
}

func (e *Entity) IsAProtein() bool {
	return e.Type == AProteinTypeEntity
}

func (e *Entity) IsBProtein() bool {
	return e.Type == BProteinTypeEntity
}

func (e *Entity) IsCProtein() bool {
	return e.Type == CProteinTypeEntity
}

func (e *Entity) IsDProtein() bool {
	return e.Type == DProteinTypeEntity
}

func (e *Entity) IsProtein() bool {
	return e.IsAProtein() || e.IsBProtein() || e.IsCProtein() || e.IsDProtein()
}

func (e *Entity) IsBasic() bool {
	return e.Type == BasicTypeEntity
}

func (e *Entity) IsEmpty() bool {
	return e.Type == ""
}

func (e *Entity) IsMy() bool {
	return e.Owner == 1
}

func (e *Entity) IsOpponent() bool {
	return e.Owner == 0
}

func (e *Entity) IsNeutral() bool {
	return e.Owner == -1
}

func (e *Entity) ToLog() string {
	return fmt.Sprintf("(%d:%d)%s:%d", e.Pos.X, e.Pos.Y, e.Type, e.Owner)
}

func (e *Entity) ID() string {
	return fmt.Sprintf("(%d:%d)", e.Pos.X, e.Pos.Y)
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
	// step := 0
	for {
		state := game.State()

		state.ScanStocks()
		state.ScanReqActions()

		state.DoAction()
		// state.Debug()
		// DebugMsg("step: ", step)
		// step += 1
	}
}

type Organs map[string]struct{}

/*
Organ		A	B	C	D
BASIC		1	0	0	0
HARVESTER	0	0	1	1
TENTACLE	0	1	1	0
SPORER		0	1	0	1
ROOT		1	1	1	1
*/
func (s *State) AvailableOrang() Organs {
	result := make(map[string]struct{}, 0)

	if s.MyStock.A > 0 {
		result[BasicTypeEntity] = struct{}{}
	}
	if s.MyStock.C > 0 && s.MyStock.D > 0 {
		result[HarvesterTypeEntity] = struct{}{}
	}
	if s.MyStock.B > 0 && s.MyStock.C > 0 {
		result[TentacleTypeEntity] = struct{}{}
	}
	if s.MyStock.B > 0 && s.MyStock.D > 0 {
		result[SporerTypeEntity] = struct{}{}
	}
	if s.MyStock.A > 0 && s.MyStock.B > 0 &&
		s.MyStock.C > 0 && s.MyStock.D > 0 {
		result[RootTypeEntity] = struct{}{}
	}

	return result
}

func (o Organs) HasBasic() bool {
	_, ok := o[BasicTypeEntity]
	return ok
}

func (o Organs) HasTentacle() bool {
	_, ok := o[TentacleTypeEntity]
	return ok
}

func (o Organs) HasRoot() bool {
	_, ok := o[RootTypeEntity]
	return ok
}

func (o Organs) HasSporer() bool {
	_, ok := o[SporerTypeEntity]
	return ok
}

func (o Organs) HasHarvester() bool {
	_, ok := o[HarvesterTypeEntity]
	return ok
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
	mySporer    []*Entity
	oppEntities []*Entity
	entities    []*Entity
	proteins    []*Entity

	nextEntity []*Entity
	nextHash   map[string]*Entity
	freePos    []*Entity
	w          int
	h          int
}

func (s *State) Scan() {
	fmt.Scan(&s.EntityCount)
}

func (s *State) ScanEnties() {
	s.entities = make([]*Entity, 0)
	s.myEntities = make([]*Entity, 0)
	s.mySporer = make([]*Entity, 0)
	s.oppEntities = make([]*Entity, 0)
	s.proteins = make([]*Entity, 0)
	s.nextEntity = make([]*Entity, 0)
	s.nextHash = make(map[string]*Entity, 0)
	s.matrix = make([][]*Entity, 0)
	for i := 0; i < s.w; i++ {
		s.matrix = append(s.matrix, make([]*Entity, s.h))
	}

	for i := 0; i < s.EntityCount; i++ {
		e := NewEntity()
		if e.IsMy() {
			s.myEntities = append(s.myEntities, e)
			if e.IsSporer() {
				s.mySporer = append(s.mySporer, e)
			}
		} else if e.IsOpponent() {
			s.oppEntities = append(s.oppEntities, e)
		} else {
			s.entities = append(s.entities, e)
			if e.IsProtein() {
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
		_ = s.GetFreePos()
		s.walk(0, 0, s.Dummy)
	}
	organs := s.AvailableOrang()
	DebugMsg("organs: ", organs, len(s.nextEntity))

	for _, e := range s.nextEntity {
		DebugMsg(">", e.ToLog())
		// if len(s.mySporer) > 0 && organs.HasRoot() {
		// dir := s.mySporer[0].OrganDir
		// if dir == DirE {
		// e.Protein.Pos.Y = e.Pos.Y
		// } else if dir == DirW {
		// e.Protein.Pos.Y = e.Pos.Y
		// } else if dir == DirN {
		// e.Protein.Pos.X = e.Pos.X
		// } else if dir == DirS {
		// e.Protein.Pos.X = e.Pos.X
		// }
		// fmt.Println(e.Spore())
		// } else if organs.HasSporer() {
		// fmt.Println(e.GrowSporer(e.OrganDir))
		// } else
		// if organs.HasHarvester() {
		// if e.NextDistance == 1 {
		// fmt.Println(e.GrowHarvester(e.OrganDir))
		// }
		// if organs.HasBasic() {
		// fmt.Println(e.GrowBasic())
		// }
		// } else
		if e.NextDistance == 1 && organs.HasHarvester() {
			fmt.Println(e.GrowHarvester(e.OrganDir))
			continue
		}
		if organs.HasBasic() {
			fmt.Println(e.GrowBasic())
			continue
		}
		if organs.HasHarvester() {
			fmt.Println(e.GrowHarvester(e.OrganDir))
			continue
		}
		if organs.HasSporer() {
			fmt.Println(e.GrowSporer(e.OrganDir))
			continue
		}
		if organs.HasTentacle() {
			fmt.Println(e.GrowTentacle(e.OrganDir))
			continue
		}
		// }
		fmt.Println("WAIT") // Write action to stdout
	}

	if len(s.nextEntity) == 0 {
		for i := 0; i < s.RequiredActionsCount; i++ {
			fmt.Println("WAIT") // Write action to stdout
		}
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
	if p.Y >= len(s.matrix) || p.Y < 0 {
		return nil
	}
	row := s.matrix[p.Y]
	if p.X >= len(row) || p.X < 0 {
		return nil
	}
	return row[p.X]
}

func (s *State) GetFreePos() []*Entity {
	s.freePos = make([]*Entity, 0)
	do := func(e *Entity, useProtein bool) {
		if e == nil {
			return
		}
		dirs := []Position{
			e.Pos.Up(),
			e.Pos.Down(),
			e.Pos.Left(),
			e.Pos.Right(),
		}
		for _, pos := range dirs {
			if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
				continue
			}
			newPos := s.getByPos(pos)
			if newPos == nil || (useProtein && newPos.IsProtein()) {
				newPos = &Entity{Pos: pos}
				newPos.OrganID = e.OrganID
				s.freePos = append(s.freePos, newPos)
			}
		}
	}

	for _, e := range s.myEntities {
		do(e, false)
	}
	if len(s.freePos) == 0 {
		for _, e := range s.myEntities {
			do(e, true)
		}
		DebugMsg(">> len free", len(s.freePos))
	}

	return s.freePos
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

func (s *State) GetSporerDir(from, to *Entity) string {
	if from == nil || to == nil {
		return ""
	}
	dx := to.Pos.X - from.Pos.X
	dy := to.Pos.Y - from.Pos.Y

	if dx == 0 {
		if dy > 0 {
			return DirN
		}
		return DirS
	}
	if dy == 0 {
		if dx > 0 {
			return DirE
		}
		return DirW
	}

	tanF := float64(dy) / float64(dx)
	if tanF < 1 && tanF > -1 {
		return DirE
	}

	if tanF < -1 || tanF > 1 {
		return DirW
	}

	if dy < 0 {
		return DirN
	}

	return DirS
}

func NewState(h, w int) *State {
	s := &State{
		entities:    make([]*Entity, 0),
		myEntities:  make([]*Entity, 0),
		mySporer:    make([]*Entity, 0),
		oppEntities: make([]*Entity, 0),
		proteins:    make([]*Entity, 0),
		nextEntity:  make([]*Entity, 0),
		nextHash:    make(map[string]*Entity, 0),
		matrix:      make([][]*Entity, 0),
		w:           w,
		h:           h,
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
