package main

import (
	"math"
	"sort"
)

func (s *State) Dummy(e *Entity) bool {

	if len(s.freePos) == 0 {
		return false
	}

	s.attackPosition = make([]*Entity, 0)

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

	do := func(i int, free *Entity, freePos []*Entity) {
		if free != nil && !free.IsProtein() && s.NearOppoent(free) {
			s.attackPosition = append(s.attackPosition, free)
		}
		if s.MyStock.CanDefend() {
			dirs := free.Pos.Get2RoseLocality()
			for _, dir := range dirs {
				if len(dir) != 2 {
					continue
				}
				step1 := dir[0]
				step2 := dir[1]
				pos1 := s.getByPos(step1)
				pos2 := s.getByPos(step2)
				if pos1 != nil && pos1.IsWall() {
					continue
				}
				if pos1 != nil && pos1.IsProtein() && s.MyStock.GetProduction(pos1.Type) > 0 {
					if pos2 != nil && pos2.IsOpponent() {
						free.NeedDefend = true
						free.DefendEntity = pos2
						continue
					}
				}
				if pos1 == nil {
					if pos2 != nil && pos2.IsOpponent() {
						free.NeedDefend = true
						free.DefendEntity = pos2
						continue
					}
				}
			}
		}
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
					free.NextDistance = math.Abs(newDistance)
					free.CanAttack = true
					freePos[i] = free
				}
				s.scanOppoent[opp.ID()] = opp
			}
		} else {
			proteins := s.GetOrderedProtens()
			for i, p := range proteins {
				p.NextDistance = free.Pos.EucleadDistance(p.Pos)
				proteins[i] = p
			}
			sort.Slice(proteins, func(i, j int) bool {
				return proteins[i].NextDistance < proteins[j].NextDistance
			})
			if len(proteins) > NearestProteins {
				proteins = append(proteins[:0], proteins[:NearestProteins]...)
			}
			for _, protein := range proteins {
				if _, ok := s.eatProtein[protein.ID()]; ok {
					continue
				}
				if s.AroundOpponet(protein) {
					continue
				}
				cost, _ := s.PathScore(free.Pos, protein.Pos)
				// DebugMsg("f>>:", free.ToLog(), protein.ToLog(), cost)
				if cost >= MaxScorePath {
					continue
				}
				// if cost == 0 {
				// cost = free.Pos.EucleadDistance(protein.Pos)
				// }
				// cost := free.Pos.EucleadDistance(protein.Pos)
				if math.Abs(cost) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && cost >= 0) {
					if s.MyStock.NeedCollectProtein(protein.Type) {
						free.NextDistance = math.Abs(cost)
					} else {
						free.NextDistance = math.Abs(cost) / s.MyStock.GetPercent(protein.Type)
					}
					free.Protein = protein
					free.Cost = cost
					freePos[i] = free
				}
			}
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
				free.NextDistance = 0.1
				free.CanAttack = true
				// DebugMsg("exist opponent:", free.ToLog(), e.ToLog())
				freePos[i] = free
			}
		}
	}

	for i, free := range s.freePos {
		do(i, free, s.freePos)
	}
	for i, free := range s.nearNotEatProteins {
		do(i, free, s.nearNotEatProteins)
	}
	zero := s.filterZeroDistance()
	sort.Slice(s.freePos, func(i, j int) bool {
		return s.freePos[i].NextDistance < s.freePos[j].NextDistance
	})

	if len(s.freePos) == 0 {
		s.freePos = append(s.freePos, s.nearNotEatProteins...)
		zero = append(zero, s.filterZeroDistance()...)
		sort.Slice(s.freePos, func(i, j int) bool {
			return s.freePos[i].NextDistance < s.freePos[j].NextDistance
		})
	}
	if len(s.freePos) > 0 {
		min := s.freePos[0]
		maxDistance := make([]struct{}, 0)
		for _, free := range s.freePos {
			if free.NextDistance >= MaxScorePath {
				maxDistance = append(maxDistance, struct{}{})
				DebugMsg("skip cell", free.ToLog())
			}
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
				NeedDefend:   free.NeedDefend,
				OrganRootID:  free.OrganRootID,
			}
			if free.CanAttack {
				s.matrix[free.Pos.Y][free.Pos.X].Type = AttackTypeEntity
				s.nextEntity = append(s.nextEntity, free)
				continue
			}
			free.OrganDir = s.GetHarvesterDir(free)
			if free.NextDistance == 1 {
				free.OrganDir = s.GetHarvesterDir(min)
			}
			s.nextEntity = append(s.nextEntity, free)
		}
		if len(maxDistance) >= len(s.freePos) {
			for _, nearProtein := range s.nearProteins {
				if s.underAttack(nearProtein.Pos) {
					continue
				}
				nearProtein.NextDistance = 0.0
				if _, ok := s.nextHash[nearProtein.ID()]; !ok {
					s.nextHash[nearProtein.ID()] = nearProtein
					s.nextEntity = append(s.nextEntity, nearProtein)
				}
			}
			return false
		}
	}

	if len(s.nextEntity) == 0 {
		s.nextEntity = append(s.nextEntity[:0], zero...)
	}

	return false
}
