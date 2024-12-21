package main

import "fmt"

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
