package main

import "math"

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
				free.Protein = protein
				freePos[i] = free
			}
		}
	}

	min := freePos[0]
	for _, free := range freePos {
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

	s.nextEntity = append(s.nextEntity, min)

	return false
}
