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
	oppEntities []*Entity
	entities    []*Entity
	proteins    []*Entity

	nextEntity []*Entity
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

	for i := 0; i < s.EntityCount; i++ {
		e := NewEntity()
		if e.Owner == 1 {
			s.myEntities = append(s.myEntities, e)
			if e.IsSporer() {
				s.mySporer = append(s.mySporer, e)
			}
		} else if e.Owner == 0 {
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
		s.walk(0, 0, s.Dummy)
	}
	organs := s.AvailableOrang()
	DebugMsg("organs: ", organs, len(s.nextEntity))

	for _, e := range s.nextEntity {
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
	freePos := make([]*Entity, 0)
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
				freePos = append(freePos, newPos)
			}
		}
	}

	for _, e := range s.myEntities {
		do(e, false)
	}
	// if len(freePos) == 0 {
	// for _, e := range s.myEntities {
	// freePos = append(freePos, e)
	// }
	// }

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
