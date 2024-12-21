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
