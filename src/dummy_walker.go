package main

import "math"

func (s *State) Dummy(e *Entity) bool {

	if len(s.freePos) == 0 {
		return false
	}

	for _, protein := range s.proteins {
		if _, ok := s.eatProtein[protein.ID()]; ok {
			continue
		}
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
		if _, ok := s.nextHash[free.ID()]; ok {
			continue
		}
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
