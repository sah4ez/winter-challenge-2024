package main

import (
	"fmt"
	"os"
	"sort"
)

type State struct {
	EntityCount          int
	RequiredActionsCount int

	MyStock       *Stock
	OpponentStock *Stock

	matrix [][]*Entity

	myEntities  []*Entity
	mySporer    []*Entity
	myRoot      []*Entity
	oppEntities []*Entity
	oppRoot     []*Entity
	entities    []*Entity
	proteins    []*Entity

	nextEntity      []*Entity
	nextHash        map[string]*Entity
	freePos         []*Entity
	nearProteins    []*Entity
	eatProtein      map[string]*Entity
	scanOppoent     map[string]*Entity
	localityOppoent map[string]*Entity

	proteinsClusters Clusters
	opponentClusters Clusters
	myClusters       Clusters

	w int
	h int
}

func (s *State) Scan() {
	fmt.Scan(&s.EntityCount)
}

func (s *State) ScanEnties() {
	s.entities = make([]*Entity, 0)
	s.myEntities = make([]*Entity, 0)
	s.mySporer = make([]*Entity, 0)
	s.myRoot = make([]*Entity, 0)
	s.oppEntities = make([]*Entity, 0)
	s.proteins = make([]*Entity, 0)
	s.oppRoot = make([]*Entity, 0)

	s.nextEntity = make([]*Entity, 0)
	s.nextHash = make(map[string]*Entity, 0)
	s.eatProtein = make(map[string]*Entity, 0)
	s.scanOppoent = make(map[string]*Entity, 0)
	s.localityOppoent = make(map[string]*Entity, 0)
	s.matrix = make([][]*Entity, 0)
	for i := 0; i < s.w; i++ {
		s.matrix = append(s.matrix, make([]*Entity, s.h))
	}
	var proteinsObs Observations
	var opponentObs Observations
	var myObs Observations

	for i := 0; i < s.EntityCount; i++ {
		e := NewEntity()
		e.State = s
		if e.IsMy() {
			s.myEntities = append(s.myEntities, e)
			myObs = append(myObs, e.Pos.ToCoordinates())
			if e.IsSporer() {
				s.mySporer = append(s.mySporer, e)
			}
			if e.IsRoot() {
				s.myRoot = append(s.myRoot, e)
			}
		} else if e.IsOpponent() {
			opponentObs = append(opponentObs, e.Pos.ToCoordinates())
			s.oppEntities = append(s.oppEntities, e)
			if e.IsRoot() {
				s.oppRoot = append(s.oppRoot, e)
			}
		} else {
			s.entities = append(s.entities, e)
			if e.IsProtein() {
				s.proteins = append(s.proteins, e)
			}
		}
		s.matrix[e.Pos.Y][e.Pos.X] = e
	}

	for _, e := range s.proteins {
		if s.HasNeigbourHarvester(e) {
			s.eatProtein[e.ID()] = e
		}
	}

	{
		sort.Slice(s.proteins, func(i, j int) bool {
			return s.proteins[i].OrganID > s.proteins[j].OrganID
		})
		for _, p := range s.proteins {
			proteinsObs = append(proteinsObs, p.Pos.ToCoordinates())
		}
		km := NewKmenas()
		k := len(s.myRoot)
		if k%2 != 0 {
			k += 1
		}
		s.proteinsClusters, _ = km.Partition(proteinsObs, k)
	}
	// {
	// km := NewKmenas()
	// s.opponentClusters, _ = km.Partition(opponentObs, len(s.oppRoot))
	// }
	// {
	// km := NewKmenas()
	// s.myClusters, _ = km.Partition(proteinsObs, len(s.myRoot))
	// }
	markCoordinates := func(cc Clusters) {
		for _, c := range cc {
			p := FromCoordinates(c.Center.Coordinates())
			if s.matrix[p.Y][p.X] != nil {
				s.matrix[p.Y][p.X].ClusterCenter = true
			} else {
				s.matrix[p.Y][p.X] = &Entity{
					Pos:           p,
					ClusterCenter: true,
					// Type:          FreeTypeEntity,
					Owner: -1,
				}
			}
		}
	}
	markCoordinates(s.proteinsClusters)
	// markCoordinates(s.opponentClusters)
	// markCoordinates(s.myClusters)
}

func (s *State) GetOrderedProtens() []*Entity {
	hashProteins := make(map[string][]*Entity, 0)
	for _, p := range s.proteins {
		if _, ok := hashProteins[p.Type]; !ok {
			hashProteins[p.Type] = []*Entity{p}
			continue
		}
		hashProteins[p.Type] = append(hashProteins[p.Type], p)
	}
	result := make([]*Entity, 0)
	order := s.MyStock.GetOrderByCountAsc()
	for k := range hashProteins {
		if s.MyStock.GetProduction(k) >= len(s.myRoot) {
			delete(hashProteins, k)
		}
	}

	for _, o := range order {
		if _, ok := hashProteins[o]; ok {
			result = append(result, hashProteins[o]...)
		}
	}

	return result
}

func (s *State) ScanStocks() {
	s.MyStock = NewStock()
	s.OpponentStock = NewStock()

	for _, ep := range s.eatProtein {
		s.MyStock.IncByType(ep.Type)
	}
	DebugMsg("eat protens:", s.MyStock.StockProduction())
}

func (s *State) ScanReqActions() {
	// requiredActionsCount: your number of organisms, output an action for each one in any order
	fmt.Scan(&s.RequiredActionsCount)
}

func (s *State) DoAction(g *Game) {
	organs := s.AvailableOrang()
	_ = s.GetFreePos(organs)
	_ = s.GetNearProteins()
	s.Dummy(nil)
	for _, e := range s.nextEntity {
		DebugMsg(">>>", e.ToLog())
	}
	DebugMsg("organs: ", organs, len(s.nextEntity))
	DebugMsg("proteins: ", s.MyStock, s.MyStock.GetOrderByCountAsc())
	if len(s.nextEntity) == 0 {
		for i := 0; i < s.RequiredActionsCount; i++ {
			fmt.Println("WAIT") // Write action to stdout
		}
		return
	}

	rootUsed := make(map[int]struct{}, 0)

	for i := 0; i < s.RequiredActionsCount; i++ {
		e := s.first()
		if e == nil {
			fmt.Println("WAIT") // Write action to stdout
			continue
		}
		if _, ok := rootUsed[e.OrganRootID]; ok {
			used := true
			for used {
				e = s.first()
				if e == nil {
					break
				}
				DebugMsg(">", e.ToLog(), rootUsed)
				if _, ok := rootUsed[e.OrganRootID]; !ok {
					used = false
				}
			}
		}
		if e == nil {
			fmt.Println("WAIT") // Write action to stdout
			continue
		}
		rootUsed[e.OrganRootID] = struct{}{}
		if g.HasSporerPoints() {
			DebugMsg("has sporer points")
		}

		if e.CanAttack {
			if organs.HasTentacle() {
				fmt.Println(e.GrowTentacle(s.GetTentacleDir2(e)))
				continue
			}
		}
		if len(s.mySporer) > 0 && organs.HasRoot() && g.HasSporerPoints() {
			from, to := g.SporerPonits()
			g.StopSporer()
			// DebugMsg("sporer stop", from.ToLog(), "->", to.ToLog())
			sporer := s.getByPos(from)
			if sporer != nil {
				dir := sporer.OrganDir
				shift := true
				var cancelShift func()
				total := 0
				for shift {
					if total == 8 {
						break
					}
					switch dir {
					case DirE:
						from.X += 1
						if from.X >= to.X {
							shift = false
						}
						cancelShift = func() {
							from.X -= 1
						}
					case DirW:
						from.X -= 1
						if from.X <= to.X {
							shift = false
						}
						cancelShift = func() {
							from.X += 1
						}
					case DirN:
						from.Y -= 1
						if from.Y <= to.Y {
							shift = false
						}
						cancelShift = func() {
							from.Y += 1
						}
					case DirS:
						from.Y += 1
						if from.Y >= to.Y {
							shift = false
						}
						cancelShift = func() {
							from.Y -= 1
						}
					}
					pos := s.getByPos(from)
					if pos != nil && pos.IsWall() {
						break
					}
					total += 1
				}
				DebugMsg("sporer condition:", from.ToLog())
				pos := s.getByPos(from)
				if pos == nil {
					sporer.SporeTo = from
					fmt.Println(sporer.Spore())
					continue
				}
				DebugMsg("sporer condition:", from.ToLog(), pos.ToLog())
				if pos.IsEmpty() || pos.IsFree() {
					sporer.SporeTo = from
					fmt.Println(sporer.Spore())
					continue
				}
				if pos.IsWall() && cancelShift != nil {
					cancelShift()
					sporer.SporeTo = from
					fmt.Println(sporer.Spore())
					continue
				}
			}
		}

		if s.MyStock.CanAttack() {
			if organs.HasTentacle() {
				fmt.Println(e.GrowTentacle(s.GetTentacleDir2(e)))
				continue
			}
		}

		if s.HasProtein(e) && organs.HasHarvester() {
			if dir := s.GetHarvesterDir(e); dir != "" {
				fmt.Println(e.GrowHarvester(dir))
				continue
			}
		}
		if len(s.myRoot) < len(s.oppRoot) || len(s.myRoot) >= len(s.mySporer) {
			if organs.HasSporer() && organs.HasRoot() && s.MyStock.D >= 2 {
				clusterID := s.proteinsClusters.Nearest(e.Pos.ToCoordinates())
				if len(s.proteinsClusters) > 0 {
					cluster := s.proteinsClusters[clusterID]
					clusterCenter := FromCoordinates(cluster.Center.Coordinates())
					centerEntites := s.getByPos(clusterCenter)
					dir := s.GetSporerDir(e, clusterCenter)
					fakeSporer := &Entity{Pos: e.Pos, OrganDir: dir}
					placeBeforeSporer := fakeSporer.SporerFirstCellPosition()
					ee := s.getByPos(placeBeforeSporer)
					if s.InMatrix(placeBeforeSporer) &&
						!centerEntites.IsMy() &&
						!centerEntites.IsOpponent() &&
						(ee == nil || ee.IsFree() || ee.IsEmpty()) {
						g.StartSporer(e.Pos, clusterCenter)
						DebugMsg("sporer start:", e.ToLog(), "cluster center", clusterCenter.ToLog())
						fmt.Println(e.GrowSporer(dir))
						continue
					} else {
						DebugMsg("not free", e.ToLog(), "cluster center", clusterCenter.ToLog())
					}
				}
			}
		}
		// как-то криво работает, упал на 20 позиций
		//if s.canStartAttack(e.Pos) && organsattack.HasTentacle() {
		//	newE := *e
		//	newE.OrganDir = s.GetTentacleDir2(e)
		//	pos := e.TentacleAttackPosition()
		//	attackPos := s.getByPos(pos)
		//	if !attackPos.IsMy() {
		//		fmt.Println(e.GrowTentacle(newE.OrganDir))
		//		continue
		//	}
		//}
		if organs.HasBasic() {
			fmt.Println(e.GrowBasic())
			continue
		}
		if organs.HasHarvester() {
			if dir := s.GetHarvesterDir(e); dir != "" {
				fmt.Println(e.GrowHarvester(dir))
				continue
			}
		}

		if organs.HasTentacle() {
			fmt.Println(e.GrowTentacle(s.GetTentacleDir2(e)))
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
		fmt.Println("WAIT") // Write action to stdout
	}
}

func (s *State) Debug(full bool) {
	DebugMsg("my", *s.MyStock)
	DebugMsg("opp", *s.OpponentStock)
	for i, k := range s.matrix {
		for j, n := range k {
			if n != nil {
				var clusterCenter = ""
				if n.ClusterCenter {
					clusterCenter = ClusterCenter
				}
				nType := n.Type
				if nType == "" {
					nType = " "
				}
				if full {
					fmt.Fprintf(os.Stderr, " %c%s(%d;%d;%.2f) ", nType[0], clusterCenter, j, i, n.Cost)
				} else {
					if clusterCenter == "" {
						fmt.Fprintf(os.Stderr, " %c ", n.Type[0])
					} else {
						fmt.Fprintf(os.Stderr, " %s ", clusterCenter)
					}
				}
			} else {
				if full {
					fmt.Fprintf(os.Stderr, " _(%d;%d) ", j, i)
				} else {
					fmt.Fprintf(os.Stderr, " _ ")
				}

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

func (s *State) setByPos(e *Entity) {
	if e == nil {
		return
	}
	p := e.Pos
	if p.Y >= len(s.matrix) || p.Y < 0 {
		return
	}
	row := s.matrix[p.Y]
	if p.X >= len(row) || p.X < 0 {
		return
	}
	s.matrix[p.Y][p.X] = e
}

func (s *State) get(x, y int) (e *Entity) {
	p := Position{X: x, Y: y}
	if p.Y >= len(s.matrix) || p.Y < 0 {
		return nil
	}
	row := s.matrix[p.Y]
	if p.X >= len(row) || p.X < 0 {
		return nil
	}
	e = row[p.X]
	if e == nil {
		e = &Entity{Pos: p, Owner: -1, State: s, Cost: 1}
	}
	return e
}

func (s *State) Width() int {
	return s.w
}

func (s *State) Height() int {
	return s.h
}

func (s *State) GetNearProteins() []*Entity {
	s.nearProteins = make([]*Entity, 0)
	do := func(e *Entity) {
		if e == nil {
			return
		}

		nearMe := false
		nearOpponent := false

		dirs := e.Pos.GetRoseLocality()
		for _, pos := range dirs {
			if !s.InMatrix(pos) {
				continue
			}
			newPos := s.getByPos(pos)
			if newPos == nil {
				continue
			}
			if newPos.IsMy() {
				nearMe = true
				e.OrganID = newPos.OrganID
				e.Owner = newPos.Owner
			}
			if newPos.IsOpponent() {
				nearOpponent = true
			}
		}
		if nearMe && nearOpponent {
			s.nearProteins = append(s.nearProteins, e)
		}
	}

	for _, e := range s.proteins {
		do(e)
	}

	return s.nearProteins
}

func (s *State) underAttack(pos Position) bool {
	newPos := &Entity{Pos: pos, Type: AttackTypeEntity}
	newDirs := newPos.Pos.GetLocality()
	for _, nd := range newDirs {
		if !s.InMatrix(nd) {
			continue
		}
		e := s.getByPos(nd)
		if e != nil && e.IsOpponent() && e.IsTentacle() {
			posAttack := e.TentacleAttackPosition()
			if posAttack.Equal(newPos.Pos) {
				// DebugMsg("attack ->", newPos.ToLog())
				return true
			}
		}
	}
	return false
}

func (s *State) canStartAttack(pos Position) bool {
	newPos := &Entity{Pos: pos, Type: AttackTypeEntity}
	newDirs := newPos.Pos.GetLocality()
	posHash := make(map[string]Position, 0)
	for _, nd := range newDirs {
		if !s.InMatrix(nd) {
			continue
		}
		for _, nnd := range nd.GetLocality() {
			if !s.InMatrix(nnd) {
				continue
			}
			if _, ok := posHash[nnd.ID()]; ok {
				continue
			}
			posHash[nd.ID()] = nnd

			e := s.getByPos(nnd)
			if e != nil {
				DebugMsg("start attack 2nd ->", pos.ToLog(), nnd.ToLog(), e.ToLog())
			}
			if e != nil && e.IsOpponent() {
				return true
			}
		}
	}
	return false
}

func (s *State) GetFreePos(organs Organs) []*Entity {
	hash := make(map[string]*Entity, 0)
	s.freePos = make([]*Entity, 0)
	do := func(e *Entity, useProtein bool) {
		if e == nil {
			return
		}
		dirs := e.Pos.GetRoseLocality()
		for _, pos := range dirs {
			if !s.InMatrix(pos) {
				continue
			}
			newPos := s.getByPos(pos)

			underAttack := s.underAttack(pos)

			if underAttack {
				continue
			}
			// хуже работает
			if newPos != nil && newPos.IsProtein() && !organs.HasHarvester() {
				typeProtein := newPos.Type
				if s.MyStock.NeedCollectProtein(newPos.Type) {
					newPos = &Entity{Pos: pos}
					newPos.Type = typeProtein
					newPos.OrganID = e.OrganID
					newPos.OrganParentID = e.OrganParentID
					newPos.OrganRootID = e.OrganRootID
					newPos.CanSpaces = true
					if _, ok := hash[newPos.ID()]; !ok {
						s.freePos = append(s.freePos, newPos)
						hash[newPos.ID()] = newPos
					}
				}
			} else if newPos != nil && newPos.IsProtein() && s.ArroundMy(newPos) {
				// typeProtein := newPos.Type
				newPos = &Entity{Pos: pos}
				newPos.OrganID = e.OrganID
				// newPos.Type = typeProtein
				newPos.OrganParentID = e.OrganParentID
				newPos.OrganRootID = e.OrganRootID
				if _, ok := hash[newPos.ID()]; !ok {
					s.freePos = append(s.freePos, newPos)
					hash[newPos.ID()] = newPos
				}
			}
			if newPos == nil || (useProtein && newPos.IsProtein()) {
				newPos = &Entity{Pos: pos}
				newPos.OrganID = e.OrganID
				newPos.OrganParentID = e.OrganParentID
				newPos.OrganRootID = e.OrganRootID
				if _, ok := hash[newPos.ID()]; !ok {
					s.freePos = append(s.freePos, newPos)
					hash[newPos.ID()] = newPos
				}
			}
		}
	}

	for _, e := range s.myEntities {
		do(e, false)
	}
	if len(s.freePos) == 0 {
		for _, e := range s.myEntities {
			do(e, true)
		}
		// DebugMsg(">> len free", len(s.freePos))
	}

	return s.freePos
}

func (s *State) GetHarvesterDir(e *Entity) string {
	if e == nil {
		return ""
	}
	hash := make(map[string]*Entity, 0)
	dirs := e.Pos.GetRoseLocality()
	for _, dir := range dirs {
		if !s.InMatrix(dir) {
			continue
		}
		p := s.getByPos(dir)
		if p == nil {
			continue
		}
		if p.IsProtein() {
			if _, ok := hash[p.Type]; !ok {
				hash[p.Type] = p
			}
		}
	}
	orderedStock := s.MyStock.GetOrderByCountAsc()
	for _, tp := range orderedStock {
		if v, ok := hash[tp]; ok {
			degree := PointToAngle(e.Pos, v.Pos)
			return AngleToDir(degree)
		}
	}
	return ""
}

func (s *State) HasNeigbourHarvester(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetRoseLocality()
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil && e.IsHarvester() && e.IsMy() {
			return true
		}
	}
	return false
}

func (s *State) HasNeigbourOpponent(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetLocality()
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil && e.IsOpponent() {
			return true
		}
	}
	return false
}

func (s *State) HasProtein(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetRoseLocality()
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil && e.IsProtein() {
			return true
		}
	}
	return false
}

func (s *State) AroundOpponet(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetLocality()
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil && (!e.IsWall() || !e.IsOpponent()) {
			return false
		}
	}
	return true
}

func (s *State) GetTentacleDir2(e *Entity) string {
	if e == nil {
		return ""
	}
	dirs := e.Pos.GetLocality()
	for _, dir := range dirs {
		if !s.InMatrix(dir) {
			continue
		}
		pos := s.getByPos(dir)
		if pos != nil && pos.IsOpponent() {
			degree := PointToAngle(e.Pos, dir)
			dir := ""
			total := 8
			for total >= 0 && dir == "" {
				dir = AngleToDir(degree)
				tentacle := *e
				tentacle.OrganDir = dir
				attackPosition := tentacle.TentacleAttackPosition()
				total -= 1
				if !s.InMatrix(attackPosition) {
					degree += 45
					dir = ""
					continue
				}
				if attackEntity := s.getByPos(attackPosition); attackEntity != nil && (attackEntity.IsWall() || attackEntity.IsMy()) {
					degree += 45
					dir = ""
					continue
				}
				return dir
			}
		}
	}
	result := make([]Entity, 0)
	for _, opp := range s.oppEntities {
		if s.FreeOppEntites(opp) {
			distance := e.Pos.EucleadDistance(opp.Pos)
			newOpp := *opp
			newOpp.NextDistance = distance
			result = append(result, newOpp)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].NextDistance < result[j].NextDistance
	})
	if len(result) > 0 {
		degree := PointToAngle(e.Pos, result[0].Pos)
		dir := ""
		total := 8
		for total >= 0 && dir == "" {
			dir = AngleToDir(degree)
			tentacle := *e
			tentacle.OrganDir = dir
			attackPosition := tentacle.TentacleAttackPosition()
			total -= 1
			if !s.InMatrix(attackPosition) {
				degree += 45
				dir = ""
				continue
			}
			if attackEntity := s.getByPos(attackPosition); attackEntity != nil && (attackEntity.IsWall() || attackEntity.IsMy()) {
				degree += 45
				dir = ""
				continue
			}
			return dir
		}
		return DirN
	}
	return DirN
}

func (s *State) GetSporerDir(from *Entity, to Position) string {
	if from == nil {
		return ""
	}
	degree := PointToAngle(from.Pos, to)
	return AngleToDir(degree)
}

func (s *State) first() *Entity {

	if len(s.nextEntity) == 0 {
		return nil
	}
	result := s.nextEntity[0]
	s.nextEntity = append(s.nextEntity[:0], s.nextEntity[1:]...)
	return result
}

func (s *State) filterZeroDistance() []*Entity {

	zero := make([]*Entity, 0)
	if len(s.freePos) == 0 {
		return zero
	}
	filtered := make([]*Entity, 0)
	for _, e := range s.freePos {
		if e == nil {
			continue
		}
		if e.NextDistance <= 0.0 {
			zero = append(zero, e)
			continue
		}
		filtered = append(filtered, e)
	}
	s.freePos = append(s.freePos[:0], filtered...)
	return zero
}

func (s *State) FreeEntites(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetRoseLocality()
	full := []struct{}{}
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil {
			if e.IsProtein() || e.IsOpponent() {
				full = append(full, struct{}{})
			}
		} else {
			full = append(full, struct{}{})
		}
	}
	return len(full) > 0
}

func (s *State) ArroundMy(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetRoseLocality()
	oneMy := []struct{}{}
	full := []struct{}{}
	total := 0
	for _, pos := range dirs {
		if !s.InMatrix(pos) {
			continue
		}
		total += 1
		e := s.getByPos(pos)
		if e != nil {
			if e.IsMy() && !e.IsHarvester() {
				oneMy = append(oneMy, struct{}{})
			}
			if e.IsWall() {
				full = append(full, struct{}{})
			}
		}
	}
	return len(full)+len(oneMy) == total && len(oneMy) >= 1
}

func (s *State) FreeOppEntites(e *Entity) bool {
	if e == nil {
		return false
	}

	dirs := e.Pos.GetRoseLocality()
	full := []struct{}{}
	for _, pos := range dirs {
		if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
			continue
		}
		e := s.getByPos(pos)
		if e != nil {
			if e.IsProtein() {
				full = append(full, struct{}{})
			}
		} else {
			full = append(full, struct{}{})
		}
	}
	return len(full) > 0
}

func (s *State) initMatrix() {
	s.matrix = make([][]*Entity, 0)
	for i := 0; i < s.w; i++ {
		s.matrix = append(s.matrix, make([]*Entity, s.h))
	}
}

func NewState(h, w int) *State {
	s := &State{
		entities:        make([]*Entity, 0),
		myEntities:      make([]*Entity, 0),
		mySporer:        make([]*Entity, 0),
		myRoot:          make([]*Entity, 0),
		oppEntities:     make([]*Entity, 0),
		proteins:        make([]*Entity, 0),
		nextEntity:      make([]*Entity, 0),
		nextHash:        make(map[string]*Entity, 0),
		matrix:          make([][]*Entity, 0),
		eatProtein:      make(map[string]*Entity, 0),
		localityOppoent: make(map[string]*Entity, 0),
		w:               w,
		h:               h,
	}

	for i := 0; i < w; i++ {
		s.matrix = append(s.matrix, make([]*Entity, h))
	}

	s.Scan()
	return s
}
