package main

import (
	"fmt"
	"os"
)

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

	nextEntity []*Entity
	nextHash   map[string]*Entity
	freePos    []*Entity
	eatProtein map[string]*Entity

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
		s.walk(0, 0, s.Dummy)
	}
	organs := s.AvailableOrang()
	DebugMsg("organs: ", organs, len(s.nextEntity))
	DebugMsg("proteins: ", s.eatProtein)

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
		if e.NextDistance <= 1 && organs.HasHarvester() {
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
		if s.HasNeigbourOpponent(e) {
			if organs.HasTentacle() {
				fmt.Println(e.GrowTentacle(s.GetTentacleDir(e)))
				continue
			}
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
		if e != nil && e.IsOpponent() {
			return true
		}
	}
	return false
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
		myRoot:      make([]*Entity, 0),
		oppEntities: make([]*Entity, 0),
		proteins:    make([]*Entity, 0),
		nextEntity:  make([]*Entity, 0),
		nextHash:    make(map[string]*Entity, 0),
		matrix:      make([][]*Entity, 0),
		eatProtein:  make(map[string]*Entity, 0),
		w:           w,
		h:           h,
	}

	for i := 0; i < w; i++ {
		s.matrix = append(s.matrix, make([]*Entity, h))
	}

	s.Scan()
	return s
}
