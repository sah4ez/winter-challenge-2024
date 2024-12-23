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
