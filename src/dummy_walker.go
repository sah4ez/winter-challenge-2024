package main

import (
	"math"
	"sort"
)

func (s *State) Dummy(e *Entity) bool {

	if len(s.freePos) == 0 {
		return false
	}

	needAttack := false
	for _, attackOverProtein := range s.nearProteins {
		if s.MyStock.NeedCollectProtein(attackOverProtein.Type) {
			DebugMsg("need collect near protein -> ", attackOverProtein.ToLog())
		}
		if s.MyStock.GetPercent(attackOverProtein.Type) < 0.4 {
			continue
		}
		if s.underAttack(attackOverProtein.Pos) {
			continue
		}

		attackOverProtein.NextDistance = 0.0
		attackOverProtein.CanAttack = true
		if _, ok := s.nextHash[attackOverProtein.ID()]; !ok {
			s.nextHash[attackOverProtein.ID()] = attackOverProtein
			s.nextEntity = append(s.nextEntity, attackOverProtein)
			DebugMsg("protein attack -> ", attackOverProtein.ToLog())
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
				newDistance := free.Pos.EucleadDistance(opp.Pos)
				if math.Abs(newDistance) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && newDistance >= 0) {
					free.NextDistance += math.Abs(newDistance)
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
				//distance := free.Pos.EucleadDistance(protein.Pos)
				//if distance > float64(s.w/4) {
				//	// сильно далеко
				//	continue
				//}
				//cost, _ := s.PathScore(free.Pos, protein.Pos)
				//if cost == 0 {
				//	cost = free.Pos.EucleadDistance(protein.Pos)
				//}
				cost := free.Pos.EucleadDistance(protein.Pos)
				if math.Abs(cost) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && cost >= 0) {
					free.NextDistance = math.Abs(cost) / s.MyStock.GetPercent(protein.Type)
					free.Protein = protein
					free.Cost = cost
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

	sort.Slice(s.freePos, func(i, j int) bool {
		return s.freePos[i].NextDistance > s.freePos[j].NextDistance
	})
	min := s.freePos[0]
	for _, free := range s.freePos {
		if _, ok := s.nextHash[free.ID()]; ok {
			continue
		}
		if free.IsProtein() && s.MyStock.NeedCollectProtein(free.Type) {
			free.NextDistance = 0.0
		}

		if free.Type == "" {
			free.Type = FreeTypeEntity
		}
		s.matrix[free.Pos.Y][free.Pos.X] = &Entity{
			Type:         free.Type,
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
