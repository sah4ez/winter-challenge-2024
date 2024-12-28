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
const AttackTypeEntity = "⚔️"
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
const DoNotUseEntityDistance = 9999

func (s *State) Dummy(e *Entity) bool {

	if len(s.freePos) == 0 {
		return false
	}

	for _, attackOverProtein := range s.nearProteins {
		if s.MyStock.GetPercent(attackOverProtein.Type) < 0.4 {
			continue
		}

		attackOverProtein.NextDistance = 0.0
		attackOverProtein.CanAttack = true
		if _, ok := s.nextHash[attackOverProtein.ID()]; !ok {
			s.nextHash[attackOverProtein.ID()] = attackOverProtein
			s.nextEntity = append(s.nextEntity, attackOverProtein)
			DebugMsg("protein attack -> ", attackOverProtein.ToLog())
		}
	}
	if len(s.nextEntity) > 0 {
		return false
	}

	for i, free := range s.freePos {
		for _, protein := range s.proteins {
			if _, ok := s.eatProtein[protein.ID()]; ok {
				continue
			}
			newDistance := free.Pos.Distance(protein.Pos)
			if math.Abs(newDistance) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && newDistance >= 0) {
				free.NextDistance = math.Abs(newDistance) / s.MyStock.GetPercent(protein.Type)
				free.Protein = protein
				s.freePos[i] = free
			}
		}

		dirs := free.Pos.GetLocality()
		for _, pos := range dirs {
			if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
				continue
			}
			e := s.getByPos(pos)
			if e != nil && e.IsOpponent() {
				if _, ok := s.localityOppoent[free.ID()]; ok {
					continue
				}
				s.localityOppoent[e.ID()] = e
				free.NextDistance = 0.0
				free.CanAttack = true
				DebugMsg("exist opponent:", free.ToLog(), e.ToLog())
				s.freePos[i] = free
			}
		}
	}

	min := s.freePos[0]
	for _, free := range s.freePos {
		if _, ok := s.nextHash[free.ID()]; ok {
			continue
		}
		s.matrix[free.Pos.Y][free.Pos.X] = &Entity{
			Type:         FreeTypeEntity,
			Pos:          free.Pos,
			NextDistance: free.NextDistance,
			Owner:        -1,
			CanAttack:    free.CanAttack,
		}
		if free.CanAttack {
			min = free
			s.matrix[free.Pos.Y][free.Pos.X].Type = AttackTypeEntity
			continue
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
	CanAttack    bool
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

func (e *Entity) IsHarvester() bool {
	return e.Type == HarvesterTypeEntity
}

func (e *Entity) IsSporer() bool {
	return e.Type == SporerTypeEntity
}

func (e *Entity) IsTentacle() bool {
	return e.Type == TentacleTypeEntity
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

func (e *Entity) IsFree() bool {
	return e.Type == FreeTypeEntity
}

func (e *Entity) IsEmpty() bool {
	return e.Type == ""
}

func (e *Entity) IsMy() bool {
	return e.Owner == 1 && !e.IsEmpty()
}

func (e *Entity) IsOpponent() bool {
	return e.Owner == 0 && !e.IsEmpty()
}

func (e *Entity) IsNeutral() bool {
	return e.Owner == -1 && !e.IsEmpty()
}

func (e *Entity) ToLog() string {
	return fmt.Sprintf("(%d:%d)%s:%d:%d:%.2f", e.Pos.X, e.Pos.Y, e.Type, e.OrganID, e.Owner, e.NextDistance)
}

func (e *Entity) ID() string {
	return fmt.Sprintf("(%d:%d)", e.Pos.X, e.Pos.Y)
}

func (e *Entity) TentacleAttackPosition() Position {
	if e.OrganDir == DirW {
		return Position{X: e.Pos.X - 1, Y: e.Pos.Y}
	}
	if e.OrganDir == DirE {
		return Position{X: e.Pos.X + 1, Y: e.Pos.Y}
	}
	if e.OrganDir == DirN {
		return Position{X: e.Pos.X, Y: e.Pos.Y - 1}
	}
	if e.OrganDir == DirS {
		return Position{X: e.Pos.X, Y: e.Pos.Y + 1}
	}
	return e.Pos
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

type State struct {
	EntityCount          int
	RequiredActionsCount int

	MyStock       *Stock
	OpponentStock *Stock

	matrix [][]*Entity

	myEntities  []*Entity
	mySporer    []*Entity
	myRoot      []*Entity
	oppEntities []*Entity
	entities    []*Entity
	proteins    []*Entity

	nextEntity      []*Entity
	nextHash        map[string]*Entity
	freePos         []*Entity
	nearProteins    []*Entity
	eatProtein      map[string]*Entity
	localityOppoent map[string]*Entity

	w int
	h int
}

func (s *State) Scan() {
	fmt.Scan(&s.EntityCount)
}

func (s *State) ScanEnties() {
	s.entities = make([]*Entity, 0)
	s.myEntities = make([]*Entity, 0)
	s.mySporer = make([]*Entity, 0)
	s.myRoot = make([]*Entity, 0)
	s.oppEntities = make([]*Entity, 0)
	s.proteins = make([]*Entity, 0)
	s.nextEntity = make([]*Entity, 0)
	s.nextHash = make(map[string]*Entity, 0)
	s.eatProtein = make(map[string]*Entity, 0)
	s.localityOppoent = make(map[string]*Entity, 0)
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
			if e.IsRoot() {
				s.myRoot = append(s.myRoot, e)
			}
		} else if e.IsOpponent() {
			s.oppEntities = append(s.oppEntities, e)
		} else {
			s.entities = append(s.entities, e)
			if e.IsProtein() {
				if s.HasNeigbourHarvester(e) {
					s.eatProtein[e.ID()] = e
				}
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
		_ = s.GetNearProteins()
		s.walk(0, 0, s.Dummy)
	}
	organs := s.AvailableOrang()
	DebugMsg("organs: ", organs, len(s.nextEntity))
	DebugMsg("proteins: ", s.MyStock)
	if len(s.nextEntity) == 0 {
		for i := 0; i < s.RequiredActionsCount; i++ {
			fmt.Println("WAIT") // Write action to stdout
		}
		return
	}

	for _, e := range s.nextEntity[:s.RequiredActionsCount] {
		DebugMsg(">", e.ToLog())
		// if len(s.mySporer) > 0 && organs.HasRoot() && len(s.myRoot)+1 == len(s.mySporer) && e.NextDistance > 1 {
		// sporer := s.mySporer[0]
		// sporer.Protein = &Entity{Pos: sporer.Pos}
		// dir := sporer.OrganDir
		// DebugMsg(">", dir)
		// var needSpore bool
		// if dir == DirE {
		// sporer.Protein.Pos.Y = e.Pos.Y
		// needSpore = sporer.Protein.Pos.X < e.Pos.X
		// } else if dir == DirW {
		// sporer.Protein.Pos.Y = e.Pos.Y
		// needSpore = sporer.Protein.Pos.X > e.Pos.X
		// } else if dir == DirN {
		// sporer.Protein.Pos.Y = e.Pos.Y
		// needSpore = sporer.Protein.Pos.Y < e.Pos.Y
		// } else if dir == DirS {
		// sporer.Protein.Pos.X = e.Pos.X
		// needSpore = sporer.Protein.Pos.Y > e.Pos.Y
		// }
		// DebugMsg(">", needSpore, sporer.Protein.ID(), e.ID(), sporer.ID())
		// if needSpore {
		// fmt.Println(sporer.Spore())
		// continue
		// }
		// }
		if e.CanAttack {
			if organs.HasTentacle() {
				fmt.Println(e.GrowTentacle(s.GetTentacleDir(e)))
				continue
			}
		}
		if s.HasProtein(e) && organs.HasHarvester() {
			if dir := s.GetHarvesterDir(e); dir != "" {
				fmt.Println(e.GrowHarvester(dir))
				continue
			}
		}
		// if len(s.mySporer) == 0 || len(s.myRoot) == len(s.mySporer) {
		// if organs.HasSporer() {
		// fmt.Println(e.GrowSporer(e.OrganDir))
		// continue
		// }
		// }
		if organs.HasBasic() {
			fmt.Println(e.GrowBasic())
			continue
		}
		if organs.HasHarvester() {
			if dir := s.GetHarvesterDir(e); dir != "" {
				fmt.Println(e.GrowHarvester(dir))
				continue
			}
		}
		if organs.HasSporer() {
			fmt.Println(e.GrowSporer(e.OrganDir))
			continue
		}

		if organs.HasTentacle() {
			fmt.Println(e.GrowTentacle(s.GetTentacleDir(e)))
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
				fmt.Fprintf(os.Stderr, " %c(%d;%d;%.2f) ", n.Type[0], j, i, n.NextDistance)
			} else {
				fmt.Fprintf(os.Stderr, " _(%d;%d) ", j, i)
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

func (s *State) GetNearProteins() []*Entity {
	s.nearProteins = make([]*Entity, 0)
	do := func(e *Entity) {
		if e == nil {
			return
		}

		nearMe := false
		nearOpponent := false

		dirs := e.Pos.GetRoseLocality()
		for _, pos := range dirs {
			if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
				continue
			}
			newPos := s.getByPos(pos)
			if newPos == nil {
				continue
			}
			if newPos.IsMy() {
				nearMe = true
				e.OrganID = newPos.OrganID
				e.Owner = newPos.Owner
			}
			if newPos.IsOpponent() {
				nearOpponent = true
			}
		}
		if nearMe && nearOpponent {
			s.nearProteins = append(s.nearProteins, e)
		}
	}

	for _, e := range s.proteins {
		do(e)
	}

	return s.nearProteins
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

			underAttack := false
			{
				newPos := &Entity{Pos: pos, Type: AttackTypeEntity}
				newDirs := newPos.Pos.GetLocality()
				for _, nd := range newDirs {
					if nd.X < 0 || nd.Y < 0 || nd.Y >= s.w || nd.X >= s.h {
						continue
					}
					e := s.getByPos(nd)
					if e != nil && e.IsOpponent() && e.IsTentacle() {
						posAttack := e.TentacleAttackPosition()
						if posAttack.Equal(newPos.Pos) {
							DebugMsg("attack ->", newPos.ToLog())
							underAttack = true
						}
					}
				}
			}

			if underAttack {
				continue
			}

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
		// DebugMsg(">> len free", len(s.freePos))
	}

	return s.freePos
}

func (s *State) GetHarvesterDir(e *Entity) string {
	if e == nil {
		return ""
	}
	up := s.getByPos(e.Pos.Up())
	if up != nil && up.IsProtein() {
		if _, ok := s.eatProtein[up.ID()]; !ok {
			return DirN
		}
	}
	down := s.getByPos(e.Pos.Down())
	if down != nil && down.IsProtein() {
		if _, ok := s.eatProtein[down.ID()]; !ok {
			return DirS
		}
	}
	left := s.getByPos(e.Pos.Left())
	if left != nil && left.IsProtein() {
		if _, ok := s.eatProtein[left.ID()]; !ok {
			return DirW
		}
	}
	right := s.getByPos(e.Pos.Right())
	if right != nil && right.IsProtein() {
		if _, ok := s.eatProtein[right.ID()]; !ok {
			return DirE
		}
	}

	return ""
}

func (s *State) HasNeigbourHarvester(e *Entity) bool {
	if e == nil {
		return false
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
		e := s.getByPos(pos)
		if e != nil && e.IsHarvester() && e.IsMy() {
			return true
		}
	}
	return false
}

func (s *State) HasNeigbourOpponent(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetLocality()
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil && e.IsOpponent() {
			return true
		}
	}
	return false
}

func (s *State) HasProtein(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetLocality()
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil && e.IsProtein() {
			return true
		}
	}
	return false
}

func (s *State) AroundOpponet(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetLocality()
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil && !e.IsWall() || !e.IsOpponent() {
			return false
		}
	}
	return true
}

func (s *State) GetTentacleDir(e *Entity) string {
	if e == nil {
		return ""
	}
	up := s.getByPos(e.Pos.Up())
	if up != nil && up.IsOpponent() {
		return DirN
	}
	down := s.getByPos(e.Pos.Down())
	if down != nil && down.IsOpponent() {
		return DirS
	}
	left := s.getByPos(e.Pos.Left())
	if left != nil && left.IsOpponent() {
		return DirW
	}
	right := s.getByPos(e.Pos.Right())
	if right != nil && right.IsOpponent() {
		return DirE
	}
	up = s.getByPos(e.Pos.Up())
	if up == nil {
		return DirN
	}
	down = s.getByPos(e.Pos.Down())
	if down == nil {
		return DirS
	}
	left = s.getByPos(e.Pos.Left())
	if left == nil {
		return DirW
	}
	right = s.getByPos(e.Pos.Right())
	if right == nil {
		return DirE
	}
	up = s.getByPos(e.Pos.Up())
	if up != nil && up.IsFree() {
		return DirN
	}
	down = s.getByPos(e.Pos.Down())
	if down != nil && down.IsFree() {
		return DirS
	}
	left = s.getByPos(e.Pos.Left())
	if left != nil && left.IsFree() {
		return DirW
	}
	right = s.getByPos(e.Pos.Right())
	if right != nil && right.IsFree() {
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
		entities:        make([]*Entity, 0),
		myEntities:      make([]*Entity, 0),
		mySporer:        make([]*Entity, 0),
		myRoot:          make([]*Entity, 0),
		oppEntities:     make([]*Entity, 0),
		proteins:        make([]*Entity, 0),
		nextEntity:      make([]*Entity, 0),
		nextHash:        make(map[string]*Entity, 0),
		matrix:          make([][]*Entity, 0),
		eatProtein:      make(map[string]*Entity, 0),
		localityOppoent: make(map[string]*Entity, 0),
		w:               w,
		h:               h,
	}

	for i := 0; i < w; i++ {
		s.matrix = append(s.matrix, make([]*Entity, h))
	}

	s.Scan()
	return s
}

type Stock struct {
	A        int
	PercentA float64

	B        int
	PercentB float64

	C        int
	PercentC float64

	D        int
	PercentD float64
}

func (s *Stock) Scan() {
	fmt.Scan(&s.A, &s.B, &s.C, &s.D)
	total := float64(s.A + s.B + s.C + s.D)
	if total > 0.0 {
		s.PercentA = 1.0 - float64(s.A)/total
		s.PercentB = 1.0 - float64(s.B)/total
		s.PercentC = 1.0 - float64(s.C)/total
		s.PercentD = 1.0 - float64(s.D)/total
	}
}

func (s *Stock) GetPercent(protein string) float64 {
	if protein == AProteinTypeEntity {
		return s.PercentA
	}
	if protein == BProteinTypeEntity {
		return s.PercentB
	}
	if protein == CProteinTypeEntity {
		return s.PercentC
	}
	if protein == DProteinTypeEntity {
		return s.PercentD
	}
	return 0.0
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
