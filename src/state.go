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
				proteinsObs = append(proteinsObs, e.Pos.ToCoordinates())
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
		km := NewKmenas()
		s.proteinsClusters, _ = km.Partition(proteinsObs, 5)
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
					Type:          FreeTypeEntity,
					Owner:         -1,
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

	for _, o := range order {
		result = append(result, hashProteins[o]...)
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
	for i := 0; i < s.RequiredActionsCount; i++ {
		_ = s.GetFreePos()
		_ = s.GetNearProteins()
		s.walk(0, 0, s.Dummy)
	}
	organs := s.AvailableOrang()
	for _, e := range s.nextEntity {
		DebugMsg(">>>", e.ToLog())
	}
	DebugMsg("organs: ", organs, len(s.nextEntity))
	DebugMsg("proteins: ", s.MyStock)
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

		if len(s.mySporer) > 0 && organs.HasRoot() && g.HasSporerPoints() {
			from, to := g.SporerPonits()
			g.StopSporer()
			// DebugMsg("sporer stop", from.ToLog(), "->", to.ToLog())
			sporer := s.getByPos(from)
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
					if from.Y <= to.Y {
						shift = false
					}
					cancelShift = func() {
						from.Y -= 1
					}
				}
				total += 1
			}
			// DebugMsg("sporer condition:", from.ToLog())
			pos := s.getByPos(from)
			if pos == nil {
				sporer.SporeTo = from
				fmt.Println(sporer.Spore())
				continue
			}
			// DebugMsg("sporer condition:", from.ToLog(), pos.ToLog())
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
		if e.CanAttack {
			if organs.HasTentacle() {
				fmt.Println(e.GrowTentacle(s.GetTentacleDir2(e)))
				continue
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
		if len(s.mySporer) == 0 || len(s.myRoot) == len(s.mySporer) {
			if organs.HasSporer() && organs.HasRoot() && s.MyStock.D >= 2 {
				clusterID := s.proteinsClusters.Nearest(e.Pos.ToCoordinates())
				if len(s.proteinsClusters) > 0 {
					cluster := s.proteinsClusters[clusterID]
					clusterCenter := FromCoordinates(cluster.Center.Coordinates())
					centerEntites := s.getByPos(clusterCenter)
					if !centerEntites.IsMy() && !centerEntites.IsOpponent() {
						g.StartSporer(e.Pos, clusterCenter)
						DebugMsg("sporer start:", e.ToLog(), "cluster center", clusterCenter.ToLog())
						fmt.Println(e.GrowSporer(s.GetSporerDir(e, clusterCenter)))
						continue
					} else {
						DebugMsg("not free", e.ToLog(), "cluster center", clusterCenter.ToLog())
					}
				}
			}
		}
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
				if full {
					fmt.Fprintf(os.Stderr, " %c%s(%d;%d;%.2f) ", n.Type[0], clusterCenter, j, i, n.NextDistance)
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
			if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
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

func (s *State) GetFreePos() []*Entity {
	s.freePos = make([]*Entity, 0)
	do := func(e *Entity, useProtein bool) {
		if e == nil {
			return
		}
		dirs := e.Pos.GetRoseLocality()
		for _, pos := range dirs {
			if pos.X < 0 || pos.Y < 0 || pos.Y >= s.w || pos.X >= s.h {
				continue
			}
			newPos := s.getByPos(pos)

			underAttack := false
			{
				newPos := &Entity{Pos: pos, Type: AttackTypeEntity}
				newDirs := newPos.Pos.GetLocality()
				for _, nd := range newDirs {
					if nd.X < 0 || nd.Y < 0 || nd.Y >= s.w || nd.X >= s.h {
						continue
					}
					e := s.getByPos(nd)
					if e != nil && e.IsOpponent() && e.IsTentacle() {
						posAttack := e.TentacleAttackPosition()
						if posAttack.Equal(newPos.Pos) {
							// DebugMsg("attack ->", newPos.ToLog())
							underAttack = true
						}
					}
				}
			}

			if underAttack {
				continue
			}
			// хуже работает
			if newPos != nil && newPos.IsProtein() {
				if s.MyStock.NeedCollectProtein(newPos.Type) {
					newPos = &Entity{Pos: pos}
					newPos.OrganID = e.OrganID
					newPos.OrganParentID = e.OrganParentID
					newPos.OrganRootID = e.OrganRootID
					s.freePos = append(s.freePos, newPos)
				}
			}

			if newPos == nil || (useProtein && newPos.IsProtein()) {
				newPos = &Entity{Pos: pos}
				newPos.OrganID = e.OrganID
				newPos.OrganParentID = e.OrganParentID
				newPos.OrganRootID = e.OrganRootID
				s.freePos = append(s.freePos, newPos)
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
	up := s.getByPos(e.Pos.Up())
	if up != nil && up.IsProtein() {
		if _, ok := s.eatProtein[up.ID()]; !ok {
			return DirN
		}
	}
	down := s.getByPos(e.Pos.Down())
	if down != nil && down.IsProtein() {
		if _, ok := s.eatProtein[down.ID()]; !ok {
			return DirS
		}
	}
	left := s.getByPos(e.Pos.Left())
	if left != nil && left.IsProtein() {
		if _, ok := s.eatProtein[left.ID()]; !ok {
			return DirW
		}
	}
	right := s.getByPos(e.Pos.Right())
	if right != nil && right.IsProtein() {
		if _, ok := s.eatProtein[right.ID()]; !ok {
			return DirE
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
		if e != nil && !e.IsWall() || !e.IsOpponent() {
			return false
		}
	}
	return true
}

func (s *State) GetTentacleDir(e *Entity) string {
	if e == nil {
		return ""
	}
	up := s.getByPos(e.Pos.Up())
	if up != nil && up.IsOpponent() {
		return DirN
	}
	down := s.getByPos(e.Pos.Down())
	if down != nil && down.IsOpponent() {
		return DirS
	}
	left := s.getByPos(e.Pos.Left())
	if left != nil && left.IsOpponent() {
		return DirW
	}
	right := s.getByPos(e.Pos.Right())
	if right != nil && right.IsOpponent() {
		return DirE
	}
	up = s.getByPos(e.Pos.Up())
	if up == nil {
		return DirN
	}
	down = s.getByPos(e.Pos.Down())
	if down == nil {
		return DirS
	}
	left = s.getByPos(e.Pos.Left())
	if left == nil {
		return DirW
	}
	right = s.getByPos(e.Pos.Right())
	if right == nil {
		return DirE
	}
	up = s.getByPos(e.Pos.Up())
	if up != nil && up.IsFree() {
		return DirN
	}
	down = s.getByPos(e.Pos.Down())
	if down != nil && down.IsFree() {
		return DirS
	}
	left = s.getByPos(e.Pos.Left())
	if left != nil && left.IsFree() {
		return DirW
	}
	right = s.getByPos(e.Pos.Right())
	if right != nil && right.IsFree() {
		return DirE
	}

	return ""
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
			return AngleToDir(degree)
		}
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
