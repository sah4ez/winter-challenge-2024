package main

import "math"

func (s *State) Dummy(e *Entity) bool {

	if len(s.freePos) == 0 {
		return false
	}

	needAttack := false
	for _, attackOverProtein := range s.nearProteins {
		if s.MyStock.GetPercent(attackOverProtein.Type) < 0.4 {
			continue
		}

		attackOverProtein.NextDistance = 0.0
		attackOverProtein.CanAttack = true
		if _, ok := s.nextHash[attackOverProtein.ID()]; !ok {
			s.nextHash[attackOverProtein.ID()] = attackOverProtein
			s.nextEntity = append(s.nextEntity, attackOverProtein)
			// DebugMsg("protein attack -> ", attackOverProtein.ToLog())
			needAttack = true
		}
	}
	if needAttack {
		return false
	}

	for i, free := range s.freePos {
		if s.MyStock.CanAttack() {
			for _, opp := range s.oppEntities {
				if _, ok := s.scanOppoent[opp.ID()]; ok {
					continue
				}
				if !s.FreeOppEntites(opp) {
					continue
				}
				newDistance := free.Pos.Distance(opp.Pos)
				if math.Abs(newDistance) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && newDistance >= 0) {
					free.NextDistance = math.Abs(newDistance)
					free.CanAttack = true
					s.freePos[i] = free
				}
				s.scanOppoent[opp.ID()] = opp
			}
		} else {
			for _, protein := range s.GetOrderedProtens() {
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
		}
		if !s.FreeEntites(free) {
			free.NextDistance += 30
		}

		dirs := free.Pos.GetLocality()
		for _, pos := range dirs {
			if !s.InMatrix(pos) {
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
				// DebugMsg("exist opponent:", free.ToLog(), e.ToLog())
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
			OrganRootID:  free.OrganRootID,
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
			}
		}
	}

	if _, ok := s.nextHash[min.ID()]; !ok {
		s.nextHash[min.ID()] = min
		s.nextEntity = append(s.nextEntity, min)
	}

	return false
}
