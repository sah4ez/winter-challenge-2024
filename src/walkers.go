package main

import "math"

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
