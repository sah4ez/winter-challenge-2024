package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
	// "testing"
)

func (s *State) PathFind(start, dest *Entity, canAttack bool) *Path {

	openNodes := minHeap{}
	heap.Push(&openNodes, &Node{Entity: dest, Cost: dest.Cost})

	checkedNodes := make([]*Entity, 0)

	hasBeenAdded := func(entity *Entity) bool {

		for _, c := range checkedNodes {
			if entity.Pos.Equal(c.Pos) {
				return true
			}
		}
		return false

	}

	path := &Path{
		Entities: make([]*Entity, 0),
	}

	if canAttack {
		if !start.Walkable() || !dest.IsOpponent() {
			return nil
		}
	} else {
		if !start.Walkable() || !dest.Walkable() {
			return nil
		}
	}

	for {

		// If the list of openNodes (nodes to check) is at 0, then we've checked all Nodes, and so the function can quit.
		if len(openNodes) == 0 {
			break
		}

		node := heap.Pop(&openNodes).(*Node)

		// If we've reached the start, then we've constructed our Path going from the destination to the start; we just have
		// to loop through each Node and go up, adding it and its parents recursively to the path.
		if node.Entity.Pos.Equal(start.Pos) {

			var t = node
			for true {
				path.Entities = append(path.Entities, t.Entity)
				t = t.Parent
				if t == nil {
					break
				}
			}

			break
		}

		// Otherwise, we add the current node's neighbors to the list of cells to check, and list of cells that have already been
		// checked (so we don't get nodes being checked multiple times).
		if node.Entity.Pos.X > 0 {
			if s.InMatrix(NewPos(node.Entity.Pos.X-1, node.Entity.Pos.Y)) {
				c := s.get(node.Entity.Pos.X-1, node.Entity.Pos.Y)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Entity.Walkable() && !hasBeenAdded(n.Entity) {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Entity)
				}
			}
		}
		if node.Entity.Pos.X < s.Height()-1 {
			if s.InMatrix(NewPos(node.Entity.Pos.X+1, node.Entity.Pos.Y)) {
				c := s.get(node.Entity.Pos.X+1, node.Entity.Pos.Y)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Entity.Walkable() && !hasBeenAdded(n.Entity) {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Entity)
				}
			}
		}

		if node.Entity.Pos.Y > 0 {
			if s.InMatrix(NewPos(node.Entity.Pos.X, node.Entity.Pos.Y-1)) {
				c := s.get(node.Entity.Pos.X, node.Entity.Pos.Y-1)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Entity.Walkable() && !hasBeenAdded(n.Entity) {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Entity)
				}
			}
		}
		if node.Entity.Pos.Y < s.Width()-1 {
			if s.InMatrix(NewPos(node.Entity.Pos.X, node.Entity.Pos.Y+1)) {
				c := s.get(node.Entity.Pos.X, node.Entity.Pos.Y+1)
				n := &Node{c, node, c.Cost + node.Cost}
				if n.Entity.Walkable() && !hasBeenAdded(n.Entity) {
					heap.Push(&openNodes, n)
					checkedNodes = append(checkedNodes, n.Entity)
				}
			}
		}
	}

	return path

}

// A Path is a struct that represents a path, or sequence of Cells from point A to point B. The Cells list is the list of Cells contained in the Path,
// and the CurrentIndex value represents the current step on the Path. Using Path.Next() and Path.Prev() advances and walks back the Path by one step.
type Path struct {
	Entities     []*Entity
	CurrentIndex int
}

// TotalCost returns the total cost of the Path (i.e. is the sum of all of the Cells in the Path).
func (p *Path) TotalCost() float64 {

	cost := 0.0
	for _, cell := range p.Entities {
		cost += cell.Cost
	}
	return cost

}

// Reverse reverses the Cells in the Path.
func (p *Path) Reverse() {

	np := []*Entity{}

	for c := len(p.Entities) - 1; c >= 0; c-- {
		np = append(np, p.Entities[c])
	}

	p.Entities = np

}

// Restart restarts the Path, so that calling path.Current() will now return the first Cell in the Path.
func (p *Path) Restart() {
	p.CurrentIndex = 0
}

// Current returns the current Cell in the Path.
func (p *Path) Current() *Entity {
	return p.Entities[p.CurrentIndex]

}

// Next returns the next cell in the path. If the Path is at the end, Next() returns nil.
func (p *Path) Next() *Entity {

	if p.CurrentIndex < len(p.Entities)-1 {
		return p.Entities[p.CurrentIndex+1]
	}
	return nil

}

// Advance advances the path by one cell.
func (p *Path) Advance() {

	p.CurrentIndex++
	if p.CurrentIndex >= len(p.Entities) {
		p.CurrentIndex = len(p.Entities) - 1
	}

}

// Prev returns the previous cell in the path. If the Path is at the start, Prev() returns nil.
func (p *Path) Prev() *Entity {

	if p.CurrentIndex > 0 {
		return p.Entities[p.CurrentIndex-1]
	}
	return nil

}

// Same returns if the Path shares the exact same cells as the other specified Path.
func (p *Path) Same(otherPath *Path) bool {

	if p == nil || otherPath == nil || len(p.Entities) != len(otherPath.Entities) {
		return false
	}

	for i := range p.Entities {
		if len(otherPath.Entities) <= i || p.Entities[i] != otherPath.Entities[i] {
			return false
		}
	}

	return true

}

// Length returns the length of the Path (how many Cells are in the Path).
func (p *Path) Length() int {
	return len(p.Entities)
}

// Get returns the Cell of the specified index in the Path. If the index is outside of the
// length of the Path, it returns -1.
func (p *Path) Get(index int) *Entity {
	if index < len(p.Entities) {
		return p.Entities[index]
	}
	return nil
}

// Index returns the index of the specified Cell in the Path. If the Cell isn't contained
// in the Path, it returns -1.
func (p *Path) Index(cell *Entity) int {
	for i, c := range p.Entities {
		if c == cell {
			return i
		}
	}
	return -1
}

// SetIndex sets the index of the Path, allowing you to safely manually manipulate the Path
// as necessary. If the index exceeds the bounds of the Path, it will be clamped.
func (p *Path) SetIndex(index int) {

	if index >= len(p.Entities) {
		p.CurrentIndex = len(p.Entities) - 1
	} else if index < 0 {
		p.CurrentIndex = 0
	} else {
		p.CurrentIndex = index
	}

}

func (p *Path) Print() {

	pos := []string{}
	for _, e := range p.Entities {
		pos = append(pos, e.ID())
	}
	DebugMsg("path:", strings.Join(pos, "->"))
}

// AtStart returns if the Path's current index is 0, the first Cell in the Path.
func (p *Path) AtStart() bool {
	return p.CurrentIndex == 0
}

// AtEnd returns if the Path's current index is the last Cell in the Path.
func (p *Path) AtEnd() bool {
	return p.CurrentIndex >= len(p.Entities)-1
}

// Node represents the node a path, it contains the cell it represents.
// Also contains other information such as the parent and the cost.
type Node struct {
	Entity *Entity
	Parent *Node
	Cost   float64
}

type minHeap []*Node

func (mH minHeap) Len() int           { return len(mH) }
func (mH minHeap) Less(i, j int) bool { return mH[i].Cost < mH[j].Cost }
func (mH minHeap) Swap(i, j int)      { mH[i], mH[j] = mH[j], mH[i] }
func (mH *minHeap) Pop() interface{} {
	old := *mH
	n := len(old)
	x := old[n-1]
	*mH = old[0 : n-1]
	return x
}

func (mH *minHeap) Push(x interface{}) {
	*mH = append(*mH, x.(*Node))
}

// A Cluster which data points gravitate around
type Cluster struct {
	Center       Coordinates
	Observations Observations
}

// Clusters is a slice of clusters
type Clusters []Cluster

// New sets up a new set of clusters and randomly seeds their initial positions
func NewCluster(k int, dataset Observations) (Clusters, error) {
	var c Clusters
	if len(dataset) == 0 || len(dataset[0].Coordinates()) == 0 {
		return c, fmt.Errorf("there must be at least one dimension in the data set")
	}
	if k == 0 {
		return c, fmt.Errorf("k must be greater than 0")
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < k; i++ {
		var p Coordinates
		for j := 0; j < len(dataset[0].Coordinates()); j++ {
			p = append(p, rand.Float64())
		}

		c = append(c, Cluster{
			Center: p,
		})
	}
	return c, nil
}

// Append adds an observation to the Cluster
func (c *Cluster) Append(point Observation) {
	c.Observations = append(c.Observations, point)
}

// Nearest returns the index of the cluster nearest to point
func (c Clusters) Nearest(point Observation) int {
	var ci int
	dist := -1.0

	// Find the nearest cluster for this data point
	for i, cluster := range c {
		d := point.Distance(cluster.Center)
		if dist < 0 || d < dist {
			dist = d
			ci = i
		}
	}

	return ci
}

// Neighbour returns the neighbouring cluster of a point along with the average distance to its points
func (c Clusters) Neighbour(point Observation, fromCluster int) (int, float64) {
	var d float64
	nc := -1

	for i, cluster := range c {
		if i == fromCluster {
			continue
		}

		cd := AverageDistance(point, cluster.Observations)
		if nc < 0 || cd < d {
			nc = i
			d = cd
		}
	}

	return nc, d
}

// Recenter recenters a cluster
func (c *Cluster) Recenter() {
	center, err := c.Observations.Center()
	if err != nil {
		return
	}

	c.Center = center
}

// Recenter recenters all clusters
func (c Clusters) Recenter() {
	for i := 0; i < len(c); i++ {
		c[i].Recenter()
	}
}

// Reset clears all point assignments
func (c Clusters) Reset() {
	for i := 0; i < len(c); i++ {
		c[i].Observations = Observations{}
	}
}

// PointsInDimension returns all coordinates in a given dimension
func (c Cluster) PointsInDimension(n int) Coordinates {
	var v []float64
	for _, p := range c.Observations {
		v = append(v, p.Coordinates()[n])
	}
	return v
}

// CentersInDimension returns all cluster centroids' coordinates in a given
// dimension
func (c Clusters) CentersInDimension(n int) Coordinates {
	var v []float64
	for _, cl := range c {
		v = append(v, cl.Center[n])
	}
	return v
}

const BasicType = "BASIC"
const GrowCmd = "GROW"
const WaitCmd = "WAIT"
const SporeCmd = "SPORE"
const WallTypeEntity = "WALL"
const RootTypeEntity = "ROOT"
const BasicTypeEntity = "BASIC"
const FreeTypeEntity = "Free"
const AttackTypeEntity = "⚔️"
const HarvesterTypeEntity = "HARVESTER"
const TentacleTypeEntity = "TENTACLE"
const SporerTypeEntity = "SPORER"
const DirW = "W"
const DirS = "S"
const DirN = "N"
const DirE = "E"
const AProteinTypeEntity = "A"
const BProteinTypeEntity = "B"
const CProteinTypeEntity = "C"
const DProteinTypeEntity = "D"
const DoNotUseEntityDistance = 9999
const MaxScorePath = 9999
const ClusterCenter = "✨"
const MaxDepthWalking = 100
const InitScore = 1000
const NearestProteins = 10
const NearestOpponent = 10
const MaxRoot = 4

func (s *State) Dummy(e *Entity) bool {

	if len(s.freePos) == 0 {
		return false
	}

	s.attackPosition = make([]*Entity, 0)
	s.scanOppoent = make(map[string]*Entity, 0)

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
				canAttack := true
				cost, _ := s.PathScore(free.Pos, opp.Pos, canAttack)
				// DebugMsg("f>>:", free.ToLog(), protein.ToLog(), cost)
				if cost >= MaxScorePath {
					continue
				}
				newDistance := free.Pos.EucleadDistance(opp.Pos)
				if math.Abs(newDistance) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && newDistance >= 0) {
					free.NextDistance = math.Abs(newDistance)
					free.CanAttack = canAttack
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
				canAttack := false
				cost, _ := s.PathScore(free.Pos, protein.Pos, canAttack)
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
				canAttack := true
				cost, _ := s.PathScore(free.Pos, pos, canAttack)
				// DebugMsg("d>>:", free.ToLog(), pos.ToLog(), cost)
				if cost > 0 && cost < MaxScorePath {
					s.localityOppoent[e.ID()] = e
					free.NextDistance = 0.1
					free.CanAttack = canAttack
					// DebugMsg("exist opponent:", free.ToLog(), e.ToLog())
					freePos[i] = free
				}
			}
		}
	}

	doAttack := func(i int, free *Entity, freePos []*Entity) {
		opponets := s.GetOrderedOpponent()
		for _, opp := range opponets {
			opp.NextDistance = free.Pos.EucleadDistance(opp.Pos)
		}
		sort.Slice(opponets, func(i, j int) bool {
			return opponets[i].NextDistance < opponets[j].NextDistance
		})
		if len(opponets) > NearestOpponent {
			opponets = append(opponets[:0], opponets[:NearestOpponent]...)
		}
		for _, opp := range opponets {
			if _, ok := s.scanOppoent[opp.ID()]; ok {
				continue
			}
			if !s.FreeOppEntites(opp) {
				continue
			}
			canAttack := true
			cost, _ := s.PathScore(free.Pos, opp.Pos, canAttack)
			// DebugMsg("d>>:", free.ToLog(), opp.ToLog(), cost)
			if cost >= MaxScorePath {
				continue
			}
			if math.Abs(cost) <= math.Abs(free.NextDistance) || (free.NextDistance == 0 && cost >= 0) {
				free.NextDistance = math.Abs(cost)
				free.CanAttack = true
				free.Cost = cost
				freePos[i] = free
			}
			s.scanOppoent[opp.ID()] = opp
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
		s.freePos = s.GetFreePosToAttack()
		for i, free := range s.freePos {
			doAttack(i, free, s.freePos)
		}
		zero = s.filterZeroDistance()
		sort.Slice(s.freePos, func(i, j int) bool {
			return s.freePos[i].NextDistance < s.freePos[j].NextDistance
		})
		DebugMsg("attack lenght", len(zero), len(s.freePos))
	}

	if len(s.freePos) == 0 {
		DebugMsg("emtpy steps...")
		s.freePos = append(s.freePos, zero...)
		for _, p := range s.nearNotEatProteins {
			if _, ok := s.eatProtein[p.Pos.ID()]; ok {
				continue
			}
			s.freePos = append(s.freePos, p)
		}
		zero = append(zero, s.filterZeroDistance()...)
		sort.Slice(s.freePos, func(i, j int) bool {
			return s.freePos[i].NextDistance < s.freePos[j].NextDistance
		})
	}

	fillNextEntity := func(freePos []*Entity) {
		maxDistance := make([]struct{}, 0)
		for _, free := range freePos {
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
			s.nextEntity = append(s.nextEntity, free)
		}
		if len(maxDistance) >= len(freePos) {
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
		}
	}

	if len(s.freePos) > 0 {
		fillNextEntity(s.freePos)
	}

	if len(s.nextEntity) == 0 {
		delta := s.MyStock.Score() - s.OpponentStock.Score()
		if delta < 1 {
			delta = 1
		}
		if len(s.myEntities)+delta == len(s.oppEntities) && s.MyStock.Score() > s.OpponentStock.Score() {
			s.freePos = append(s.freePos, zero...)
			s.freePos = append(s.freePos, s.proteins...)
			for i, free := range s.freePos {
				do(i, free, s.freePos)
			}
			sort.Slice(s.freePos, func(i, j int) bool {
				return s.freePos[i].NextDistance < s.freePos[j].NextDistance
			})
			fillNextEntity(s.freePos)
		} else if len(s.myEntities) == len(s.oppEntities) && s.MyStock.Score() <= s.OpponentStock.Score() {
			s.freePos = s.proteins
			for i, free := range s.freePos {
				do(i, free, s.freePos)
			}
			sort.Slice(s.freePos, func(i, j int) bool {
				return s.freePos[i].NextDistance < s.freePos[j].NextDistance
			})
			fillNextEntity(s.freePos)
		}
		DebugMsg("2 emtpy steps...", len(s.nextEntity))
		if len(s.nextEntity) == 0 {
			s.nextEntity = append(s.nextEntity, zero...)
		}
	}

	return false
}

type Entity struct {
	Pos           Position
	Type          string
	Owner         int
	OrganID       int
	OrganDir      string
	OrganParentID int
	OrganRootID   int

	NextDistance  float64
	Cost          float64
	Protein       *Entity
	CanAttack     bool
	NeedDefend    bool
	DefendEntity  *Entity
	ClusterCenter bool
	SporeTo       Position
	CanSpaces     bool
	State         *State
}

func (e *Entity) Scan() {
	fmt.Scan(&e.Pos.X, &e.Pos.Y, &e.Type, &e.Owner, &e.OrganID, &e.OrganDir, &e.OrganParentID, &e.OrganRootID)
	if e.IsWall() {
		e.Cost = 100.0
	}

	if e.IsOpponent() || e.IsMy() {
		e.Cost = 90.0
	}
	if e.IsProtein() {
		e.Cost = 50.0
	}
	if e.IsEmpty() {
		e.Cost = 10.0
	}
}

func (e *Entity) GrowBasic() string {
	return fmt.Sprintf("%s %d %d %d %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, BasicTypeEntity)
}

func (e *Entity) GrowHarvester(direction string) string {
	return fmt.Sprintf("%s %d %d %d %s %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, HarvesterTypeEntity, direction)
}

func (e *Entity) GrowTentacle(direction string) string {
	return fmt.Sprintf("%s %d %d %d %s %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, TentacleTypeEntity, direction)
}

func (e *Entity) GrowSporer(direction string) string {
	return fmt.Sprintf("%s %d %d %d %s %s", GrowCmd, e.OrganID, e.Pos.X, e.Pos.Y, SporerTypeEntity, direction)
}

func (e *Entity) Spore() string {
	return fmt.Sprintf("%s %d %d %d", SporeCmd, e.OrganID, e.SporeTo.X, e.SporeTo.Y)
}

func (e *Entity) IsWall() bool {
	return e.Type == WallTypeEntity
}

func (e *Entity) IsRoot() bool {
	return e.Type == RootTypeEntity
}

func (e *Entity) IsHarvester() bool {
	return e.Type == HarvesterTypeEntity
}

func (e *Entity) IsSporer() bool {
	return e.Type == SporerTypeEntity
}

func (e *Entity) IsTentacle() bool {
	return e.Type == TentacleTypeEntity
}

func (e *Entity) IsAProtein() bool {
	return e.Type == AProteinTypeEntity
}

func (e *Entity) IsBProtein() bool {
	return e.Type == BProteinTypeEntity
}

func (e *Entity) IsCProtein() bool {
	return e.Type == CProteinTypeEntity
}

func (e *Entity) IsDProtein() bool {
	return e.Type == DProteinTypeEntity
}

func (e *Entity) IsProtein() bool {
	return e.IsAProtein() || e.IsBProtein() || e.IsCProtein() || e.IsDProtein()
}

func (e *Entity) IsBasic() bool {
	return e.Type == BasicTypeEntity
}

func (e *Entity) IsFree() bool {
	return e.Type == FreeTypeEntity
}

func (e *Entity) IsEmpty() bool {
	return e.Type == ""
}

func (e *Entity) IsMy() bool {
	return e.Owner == 1 && !e.IsEmpty()
}

func (e *Entity) IsOpponent() bool {
	return e.Owner == 0 && !e.IsEmpty()
}

func (e *Entity) IsNeutral() bool {
	return e.Owner == -1 && !e.IsEmpty()
}

func (e *Entity) ToLog() string {
	return fmt.Sprintf("(%d:%d)%s:%d:%d:%.2f", e.Pos.X, e.Pos.Y, e.Type, e.OrganRootID, e.Owner, e.NextDistance)
}

func (e *Entity) ID() string {
	return fmt.Sprintf("(%d:%d)", e.Pos.X, e.Pos.Y)
}

func (e *Entity) TentacleAttackPosition() Position {
	if e.OrganDir == DirW {
		return Position{X: e.Pos.X - 1, Y: e.Pos.Y}
	}
	if e.OrganDir == DirE {
		return Position{X: e.Pos.X + 1, Y: e.Pos.Y}
	}
	if e.OrganDir == DirN {
		return Position{X: e.Pos.X, Y: e.Pos.Y - 1}
	}
	if e.OrganDir == DirS {
		return Position{X: e.Pos.X, Y: e.Pos.Y + 1}
	}
	return e.Pos
}

func (e *Entity) SporerFirstCellPosition() Position {
	if e.OrganDir == DirW {
		return Position{X: e.Pos.X - 1, Y: e.Pos.Y}
	}
	if e.OrganDir == DirE {
		return Position{X: e.Pos.X + 1, Y: e.Pos.Y}
	}
	if e.OrganDir == DirN {
		return Position{X: e.Pos.X, Y: e.Pos.Y - 1}
	}
	if e.OrganDir == DirS {
		return Position{X: e.Pos.X, Y: e.Pos.Y + 1}
	}
	return e.Pos
}

func NewEntity() *Entity {

	e := &Entity{}
	e.Scan()
	return e
}

func (e *Entity) Walkable() bool {
	return e.IsEmpty() || e.IsFree() || e.IsProtein()
}

func NewWall(x, y int) *Entity {
	return &Entity{
		Pos:   Position{X: x, Y: y},
		Owner: -1,
		Type:  WallTypeEntity,
	}
}

func NewProteinA(x, y int) *Entity {
	return &Entity{
		Pos:   Position{X: x, Y: y},
		Owner: -1,
		Type:  AProteinTypeEntity,
	}
}

func NewEntityMy(x, y int, etype string) *Entity {
	return &Entity{
		Pos:   Position{X: x, Y: y},
		Owner: 1,
		Type:  etype,
	}
}

func NewEntityOpp(x, y int, etype string) *Entity {
	return &Entity{
		Pos:   Position{X: x, Y: y},
		Owner: 0,
		Type:  etype,
	}
}

type Game struct {
	Width  int
	Height int
	state  *State

	sporerFrom *Position
	sporerTo   *Position
}

func (g *Game) Scan() {
	// width: columns in the game grid
	// height: rows in the game grid
	fmt.Scan(&g.Width, &g.Height)
}

func (g *Game) State() *State {
	g.state = NewState(g.Width, g.Height)
	// if g.state == nil {
	// g.state = NewState(g.Width, g.Height)
	// } else {
	// g.state.Scan()
	// }
	g.state.ScanEnties()
	return g.state
}

func (g *Game) StartSporer(sporerFrom Position, sporerTo Position) {
	g.sporerFrom = &sporerFrom
	g.sporerTo = &sporerTo
}

func (g *Game) StopSporer() {
	g.sporerFrom = nil
	g.sporerTo = nil
}

func (g *Game) SporerPonits() (from, to Position) {
	return *g.sporerFrom, *g.sporerTo
}

func (g *Game) HasSporerPoints() bool {
	return g.sporerFrom != nil && g.sporerTo != nil
}

func NewGame() *Game {
	g := &Game{}
	g.Scan()
	return g
}

// Package kmeans implements the k-means clustering algorithm
// See: https://en.wikipedia.org/wiki/K-means_clustering

// Kmeans configuration/option struct
type Kmeans struct {
	// when a plotter is set, Plot gets called after each iteration
	plotter Plotter
	// deltaThreshold (in percent between 0.0 and 0.1) aborts processing if
	// less than n% of data points shifted clusters in the last iteration
	deltaThreshold float64
	// iterationThreshold aborts processing when the specified amount of
	// algorithm iterations was reached
	iterationThreshold int
}

// The Plotter interface lets you implement your own plotters
type Plotter interface {
	Plot(cc Clusters, iteration int) error
}

// NewWithOptions returns a Kmeans configuration struct with custom settings
func NewWithOptions(deltaThreshold float64, plotter Plotter) (Kmeans, error) {
	if deltaThreshold <= 0.0 || deltaThreshold >= 1.0 {
		return Kmeans{}, fmt.Errorf("threshold is out of bounds (must be >0.0 and <1.0, in percent)")
	}

	return Kmeans{
		plotter:            plotter,
		deltaThreshold:     deltaThreshold,
		iterationThreshold: 96,
	}, nil
}

// New returns a Kmeans configuration struct with default settings
func NewKmenas() Kmeans {
	m, _ := NewWithOptions(0.01, nil)
	return m
}

// Partition executes the k-means algorithm on the given dataset and
// partitions it into k clusters
func (m Kmeans) Partition(dataset Observations, k int) (Clusters, error) {
	if k > len(dataset) {
		return Clusters{}, fmt.Errorf("the size of the data set must at least equal k")
	}

	cc, err := NewCluster(k, dataset)
	if err != nil {
		return cc, err
	}

	points := make([]int, len(dataset))
	changes := 1

	for i := 0; changes > 0; i++ {
		changes = 0
		cc.Reset()

		for p, point := range dataset {
			ci := cc.Nearest(point)
			cc[ci].Append(point)
			if points[p] != ci {
				points[p] = ci
				changes++
			}
		}

		for ci := 0; ci < len(cc); ci++ {
			if len(cc[ci].Observations) == 0 {
				// During the iterations, if any of the cluster centers has no
				// data points associated with it, assign a random data point
				// to it.
				// Also see: http://user.ceng.metu.edu.tr/~tcan/ceng465_f1314/Schedule/KMeansEmpty.html
				var ri int
				for {
					// find a cluster with at least two data points, otherwise
					// we're just emptying one cluster to fill another
					ri = rand.Intn(len(dataset)) //nolint:gosec // rand.Intn is good enough for this
					if len(cc[points[ri]].Observations) > 1 {
						break
					}
				}
				cc[ci].Append(dataset[ri])
				points[ri] = ci

				// Ensure that we always see at least one more iteration after
				// randomly assigning a data point to a cluster
				changes = len(dataset)
			}
		}

		if changes > 0 {
			cc.Recenter()
		}
		if m.plotter != nil {
			err := m.plotter.Plot(cc, i)
			if err != nil {
				return nil, fmt.Errorf("failed to plot chart: %s", err)
			}
		}
		if i == m.iterationThreshold ||
			changes < int(float64(len(dataset))*m.deltaThreshold) {
			// fmt.Println("Aborting:", changes, int(float64(len(dataset))*m.TerminationThreshold))
			break
		}
	}

	return cc, nil
}

/**
 * Grow and multiply your organisms to end up larger than your opponent.
 **/

func main() {
	game := NewGame()
	// step := 0
	for {
		state := game.State()

		state.ScanStocks()
		state.ScanReqActions()

		state.DoAction(game)
		// full := true
		// state.Debug(full)
		// DebugMsg("step: ", step)
		// step += 1
	}
}

// Coordinates is a slice of float64
type Coordinates []float64

// Observation is a data point (float64 between 0.0 and 1.0) in n dimensions
type Observation interface {
	Coordinates() Coordinates
	Distance(point Coordinates) float64
}

// Observations is a slice of observations
type Observations []Observation

// Coordinates implements the Observation interface for a plain set of float64
// coordinates
func (c Coordinates) Coordinates() Coordinates {
	return Coordinates(c)
}

// Distance returns the euclidean distance between two coordinates
func (c Coordinates) Distance(p2 Coordinates) float64 {
	var r float64
	for i, v := range c {
		r += math.Pow(v-p2[i], 2)
	}
	return r
}

// Center returns the center coordinates of a set of Observations
func (c Observations) Center() (Coordinates, error) {
	var l = len(c)
	if l == 0 {
		return nil, fmt.Errorf("there is no mean for an empty set of points")
	}

	cc := make([]float64, len(c[0].Coordinates()))
	for _, point := range c {
		for j, v := range point.Coordinates() {
			cc[j] += v
		}
	}

	var mean Coordinates
	for _, v := range cc {
		mean = append(mean, v/float64(l))
	}
	return mean, nil
}

// AverageDistance returns the average distance between o and all observations
func AverageDistance(o Observation, observations Observations) float64 {
	var d float64
	var l int

	for _, observation := range observations {
		dist := o.Distance(observation.Coordinates())
		if dist == 0 {
			continue
		}

		l++
		d += dist
	}

	if l == 0 {
		return 0
	}
	return d / float64(l)
}

type Organs map[string]struct{}

/*
Organ		A	B	C	D
BASIC		1	0	0	0
HARVESTER	0	0	1	1
TENTACLE	0	1	1	0
SPORER		0	1	0	1
ROOT		1	1	1	1
*/
func (s *State) AvailableOrang() Organs {
	result := make(map[string]struct{}, 0)

	if s.MyStock.A > 0 {
		result[BasicTypeEntity] = struct{}{}
	}
	if s.MyStock.C > 0 && s.MyStock.D > 0 {
		result[HarvesterTypeEntity] = struct{}{}
	}
	if s.MyStock.B > 0 && s.MyStock.C > 0 {
		result[TentacleTypeEntity] = struct{}{}
	}
	if s.MyStock.B > 1 && s.MyStock.D > 1 {
		result[SporerTypeEntity] = struct{}{}
	}
	if s.MyStock.A > 1 && s.MyStock.B > 1 &&
		s.MyStock.C > 1 && s.MyStock.D > 1 {
		result[RootTypeEntity] = struct{}{}
	}

	return result
}

func (o Organs) HasBasic() bool {
	_, ok := o[BasicTypeEntity]
	return ok
}

func (o Organs) HasTentacle() bool {
	_, ok := o[TentacleTypeEntity]
	return ok
}

func (o Organs) HasRoot() bool {
	_, ok := o[RootTypeEntity]
	return ok
}

func (o Organs) HasSporer() bool {
	_, ok := o[SporerTypeEntity]
	return ok
}

func (o Organs) HasHarvester() bool {
	_, ok := o[HarvesterTypeEntity]
	return ok
}

func (s *State) PathScore(from Position, to Position, canAttack bool) (float64, bool) {
	fromEntity := s.get(from.X, from.Y)
	toEntity := s.get(to.X, to.Y)
	// DebugMsg(">>>", fromEntity.ToLog(), toEntity.ToLog())

	path := s.PathFind(fromEntity, toEntity, canAttack)
	if path == nil {
		return 666, false
	}
	score := path.TotalCost()
	found := score > 0
	if score == 0 {
		// DebugMsg(">>>", fromEntity, toEntity)
		// path.Print()
		score = MaxScorePath
	}

	return score, found
}

//func (s *State) PathScore2(from Position, to Position, depth int, hash map[string]struct{}, score int) (int, bool) {
//	find := false
//	if depth == MaxDepthWalking {
//		return 0, find
//	}
//	if score == InitScore {
//		return InitScore, find
//	}
//	if hash == nil {
//		hash = make(map[string]struct{}, 0)
//	}
//	hash[from.ID()] = struct{}{}
//
//	depth += 1
//	dirs := from.GetRoseLocality()
//	for i, dir := range dirs {
//		if !s.InMatrix(dir) {
//			continue
//		}
//		if _, ok := hash[dir.ID()]; ok {
//			continue
//		}
//		DebugMsg(">>", score, depth, from.ToLog(), to.ToLog(), dir.ToLog())
//		e := s.getByPos(dir)
//		if e != nil && (e.IsMy() || e.IsOpponent() || e.IsWall()) {
//			continue
//		}
//		if dir.Parent != nil {
//			continue
//		}
//		dir.Parent = &from
//		dirs[i] = dir
//		if dir.Equal(to) {
//			DebugMsg("FIND>", score, depth, from.ToLog(), to.ToLog(), dir.ToLog())
//			find = true
//			return score + 10, find
//		}
//	}
//
//	score += 10
//	for _, dir := range dirs {
//		if !s.InMatrix(dir) {
//			continue
//		}
//		if _, ok := hash[dir.ID()]; ok {
//			continue
//		}
//		// hash[dir.ID()] = struct{}{}
//		DebugMsg(">>", score, depth, from.ToLog(), to.ToLog(), dir.ToLog())
//		e := s.getByPos(dir)
//		if e != nil && (e.IsMy() || e.IsOpponent() || e.IsWall()) {
//			continue
//		}
//		score, find = s.PathScore(dir, to, depth, hash, score)
//		if find {
//			break
//		}
//	}
//	return score, find
//}

func TestPathScore(t *testing.T) {
	testCases := []struct {
		name      string
		s         State
		from      Position
		to        Position
		exp       float64
		expFind   bool
		canAttack bool
		fillFn    func(*State)
	}{
		{
			name: "one step right",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(0, 1),
			exp:     20.0,
			expFind: true,
		},
		{
			name: "ten step right",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(0, 3),
			exp:     40.0,
			expFind: true,
		},
		{
			name: "1 step by diagonal",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(1, 1),
			exp:     30.0,
			expFind: true,
		},
		{
			name: "2 step by diagonal",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(2, 2),
			exp:     50.0,
			expFind: true,
		},
		{
			name: "3 step by diagonal",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(3, 3),
			exp:     70.0,
			expFind: true,
		},
		{
			name: "9 step by diagonal",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(9, 9),
			exp:     190.0,
			expFind: true,
		},
		{
			name: "9 step by diagonal and wall",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(0, 9),
			exp:     280.0,
			expFind: true,
			fillFn: func(s *State) {
				if s == nil {
					return
				}
				s.setByPos(NewWall(0, 5))
				s.setByPos(NewWall(1, 5))
				s.setByPos(NewWall(2, 5))
				s.setByPos(NewWall(3, 5))
				s.setByPos(NewWall(4, 5))
				s.setByPos(NewWall(5, 5))
				s.setByPos(NewWall(6, 5))
				s.setByPos(NewWall(7, 5))
				s.setByPos(NewWall(8, 5))
			},
		},
		{
			name: "no path",
			s: State{
				w: 3,
				h: 3,
			},
			from:    NewPos(0, 0),
			to:      NewPos(0, 2),
			exp:     MaxScorePath,
			expFind: false,
			fillFn: func(s *State) {
				if s == nil {
					return
				}
				s.setByPos(NewWall(0, 1))
				s.setByPos(NewWall(1, 1))
				s.setByPos(NewWall(2, 1))
			},
		},
		{
			name: "path through protein",
			s: State{
				w: 3,
				h: 3,
			},
			from:    NewPos(0, 0),
			to:      NewPos(0, 2),
			exp:     60.0,
			expFind: true,
			fillFn: func(s *State) {
				if s == nil {
					return
				}
				s.setByPos(NewWall(0, 1))
				s.setByPos(NewWall(1, 1))
				s.setByPos(NewProteinA(2, 1))
			},
		},
		{
			name: "reverse path through protein",
			s: State{
				w: 21,
				h: 10,
			},
			from:    NewPos(8, 17),
			to:      NewPos(6, 14),
			exp:     50.0,
			expFind: true,
			fillFn: func(s *State) {
				if s == nil {
					return
				}
				s.setByPos(NewWall(10, 18))
				s.setByPos(NewWall(9, 15))
				s.setByPos(NewWall(10, 14))
				s.setByPos(NewProteinA(9, 19))
				s.setByPos(NewProteinA(8, 19))
				s.setByPos(NewProteinA(7, 18))
				s.setByPos(NewProteinA(6, 14))
				s.setByPos(NewEntityMy(8, 19, RootTypeEntity))
				s.setByPos(NewEntityMy(8, 18, HarvesterTypeEntity))
				s.setByPos(NewEntityMy(0, 18, HarvesterTypeEntity))
			},
		},
		{
			name: "step to dead end",
			s: State{
				w: 18,
				h: 8,
			},
			from:    NewPos(3, 12),
			to:      NewPos(2, 11),
			exp:     MaxScorePath,
			expFind: false,
			fillFn: func(s *State) {
				if s == nil {
					return
				}
				s.setByPos(NewWall(2, 13))
				s.setByPos(NewWall(2, 12))
				s.setByPos(NewWall(3, 11))
				s.setByPos(NewWall(3, 11))
				s.setByPos(NewWall(4, 11))
				s.setByPos(NewWall(4, 12))
				s.setByPos(NewProteinA(2, 11))
				s.setByPos(NewEntityMy(3, 13, BasicType))
				s.setByPos(NewEntityMy(4, 13, BasicType))
			},
		},
		{
			name: "one step up",
			s: State{
				w: 18,
				h: 9,
			},
			from:    NewPos(4, 17),
			to:      NewPos(3, 17),
			exp:     10.0,
			expFind: true,
			fillFn: func(s *State) {
				if s == nil {
					return
				}
				s.setByPos(NewWall(2, 15))
				s.setByPos(NewWall(3, 15))
				s.setByPos(NewWall(4, 15))
				s.setByPos(NewWall(5, 15))
				s.setByPos(NewWall(6, 15))
				s.setByPos(NewWall(7, 15))
				s.setByPos(NewProteinA(3, 17))
				s.setByPos(NewEntityMy(6, 16, RootTypeEntity))
				s.setByPos(NewEntityMy(5, 17, RootTypeEntity))
				s.setByPos(NewEntityMy(6, 17, BasicType))
			},
		},
		{
			name: "to attack",
			s: State{
				w: 3,
				h: 3,
			},
			from:      NewPos(0, 0),
			to:        NewPos(0, 2),
			exp:       60.0,
			expFind:   true,
			canAttack: true,
			fillFn: func(s *State) {
				if s == nil {
					return
				}
				s.setByPos(NewWall(0, 1))
				s.setByPos(NewWall(1, 1))
				s.setByPos(NewEntityOpp(0, 2, RootTypeEntity))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+":"+tc.from.ToLog()+"->"+tc.to.ToLog(), func(t *testing.T) {
			tc.s.initMatrix()
			if tc.fillFn != nil {
				tc.fillFn(&tc.s)
			}
			act, find := tc.s.PathScore(tc.from, tc.to, tc.canAttack)
			if act != tc.exp {
				t.Error("unexpected score", act, tc.exp)
			}
			if find != tc.expFind {
				t.Error("unexpected find result", find, tc.expFind)
			}
		})
	}
}

type Position struct {
	X int
	Y int

	Parent *Position
}

func (p Position) ID() string {
	return fmt.Sprintf("(%d:%d)", p.X, p.Y)
}

func (p Position) Equal(pos Position) bool {
	return p.X == pos.X && p.Y == pos.Y
}

func (p Position) Shift(dir string) Position {
	if dir == DirS {
		return p.Down()
	}
	if dir == DirE {
		return p.Right()
	}
	if dir == DirW {
		return p.Left()
	}
	if dir == DirN {
		return p.Up()
	}
	return Position{}
}

func (p Position) Up() Position {
	return Position{X: p.X, Y: p.Y - 1}
}

func (p Position) Down() Position {
	return Position{X: p.X, Y: p.Y + 1}
}

func (p Position) Left() Position {
	return Position{X: p.X - 1, Y: p.Y}
}

func (p Position) Right() Position {
	return Position{X: p.X + 1, Y: p.Y}
}

func (p Position) Up2() Position {
	return Position{X: p.X, Y: p.Y - 2}
}

func (p Position) Down2() Position {
	return Position{X: p.X, Y: p.Y + 2}
}

func (p Position) Left2() Position {
	return Position{X: p.X - 2, Y: p.Y}
}

func (p Position) Right2() Position {
	return Position{X: p.X + 2, Y: p.Y}
}

func (p Position) GetLocality() []Position {
	return []Position{
		p.Up(),
		p.Down(),
		p.Left(),
		p.Right(),
		Position{X: p.X + 1, Y: p.Y + 1},
		Position{X: p.X - 1, Y: p.Y - 1},
		Position{X: p.X + 1, Y: p.Y - 1},
		Position{X: p.X - 1, Y: p.Y + 1},
	}
}

func (p Position) GetRoseLocality() []Position {
	return []Position{
		p.Up(),
		p.Down(),
		p.Left(),
		p.Right(),
	}
}

func (p Position) Get2RoseLocality() [][]Position {
	return [][]Position{
		[]Position{p.Up(), p.Up2()},
		[]Position{p.Down(), p.Down2()},
		[]Position{p.Left(), p.Left2()},
		[]Position{p.Right(), p.Right2()},
	}
}

func (from Position) EucleadDistance(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}

func (p Position) ToCoordinates() Coordinates {
	return Coordinates{
		float64(p.X),
		float64(p.Y),
	}
}

func (p Position) ToLog() string {
	depth := 0
	parent := p.Parent
	for parent != nil {
		parent = parent.Parent
		depth += 1
	}
	return fmt.Sprintf("(%d:%d)%d", p.X, p.Y, depth)
}

func FromCoordinates(c Coordinates) Position {
	if len(c) != 2 {
		panic("invalid coordinates")
	}
	return Position{X: int(c[0]), Y: int(c[1])}
}

func (s *State) InMatrix(p Position) bool {
	out := p.X < 0 || p.Y < 0 || p.Y >= s.w || p.X >= s.h
	return !out
}

func NewPos(x, y int) Position {
	return Position{X: x, Y: y}
}

type State struct {
	EntityCount          int
	RequiredActionsCount int

	MyStock       *Stock
	OpponentStock *Stock

	matrix [][]*Entity

	myEntities []*Entity
	mySporer   []*Entity
	myRoot     []*Entity

	myEntitiesByRoot map[int][]*Entity
	oppEntities      []*Entity
	oppRoot          []*Entity
	entities         []*Entity
	proteins         []*Entity

	nextEntity         []*Entity
	nextHash           map[string]*Entity
	freePos            []*Entity
	nearProteins       []*Entity
	nearNotEatProteins []*Entity
	eatProtein         map[string]*Entity
	scanOppoent        map[string]*Entity
	localityOppoent    map[string]*Entity
	attackPosition     []*Entity

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
	s.myEntitiesByRoot = make(map[int][]*Entity, 0)
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
			if _, ok := s.myEntitiesByRoot[e.OrganRootID]; ok {
				s.myEntitiesByRoot[e.OrganRootID] = append(s.myEntitiesByRoot[e.OrganRootID], e)
			} else {
				s.myEntitiesByRoot[e.OrganRootID] = []*Entity{e}
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
		k := len(s.GetRootWithFreePos())
		if k%2 != 0 {
			k += 1
		}
		s.proteinsClusters, _ = km.Partition(proteinsObs, k)
	}
	{
		km := NewKmenas()
		s.opponentClusters, _ = km.Partition(opponentObs, len(s.oppRoot))
	}
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
	markCoordinates(s.opponentClusters)
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
		if s.MyStock.GetProduction(k) > len(s.GetRootWithFreePos()) {
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

func (s *State) GetOrderedOpponent() []*Entity {
	hashOppoent := make(map[string][]*Entity, 0)
	for _, p := range s.oppEntities {
		if _, ok := hashOppoent[p.Type]; !ok {
			hashOppoent[p.Type] = []*Entity{p}
			continue
		}
		hashOppoent[p.Type] = append(hashOppoent[p.Type], p)
	}
	result := make([]*Entity, 0)
	order := []string{
		HarvesterTypeEntity,
		BasicTypeEntity,
		SporerTypeEntity,
		RootTypeEntity,
		TentacleTypeEntity,
	}

	for _, o := range order {
		if _, ok := hashOppoent[o]; ok {
			result = append(result, hashOppoent[o]...)
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

	if len(s.attackPosition) > 0 {
		for _, e := range s.attackPosition {
			if e != nil {
				DebugMsg("Under attack:", e.Pos.ToLog())
			}
		}
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
			DebugMsg("can attack")
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

		if s.MyStock.TryAttack() && len(s.GetRootWithFreePos()) <= len(s.mySporer) {
			if organs.HasTentacle() {
				fmt.Println(e.GrowTentacle(s.GetTentacleDir2(e)))
				continue
			}
		}

		if s.MyStock.CanDefend() {
			if e.NeedDefend && organs.HasTentacle() && e.DefendEntity != nil {
				degree := PointToAngle(e.Pos, e.DefendEntity.Pos)
				fmt.Println(e.GrowTentacle(AngleToDir(degree)))
				continue
			}
		}

		if s.HasProtein(e) && organs.HasHarvester() {
			if dir := s.GetHarvesterDir(e); dir != "" {
				fmt.Println(e.GrowHarvester(dir))
				continue
			}
		}
		if len(s.GetRootWithFreePos()) < len(s.oppRoot) {
			if organs.HasSporer() && organs.HasRoot() && s.MyStock.D >= 3 {
				cluters := s.proteinsClusters
				clusterID := cluters.Nearest(e.Pos.ToCoordinates())
				if len(cluters) > 0 {
					cluster := cluters[clusterID]
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
		e = &Entity{Pos: p, Owner: -1, State: s, Cost: 10}
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
	s.nearNotEatProteins = make([]*Entity, 0)
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
				// e.Owner = newPos.Owner
			}
			if newPos.IsOpponent() {
				nearOpponent = true
			}
		}
		if nearMe && !nearOpponent {
			s.nearNotEatProteins = append(s.nearNotEatProteins, e)
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
					newPos = &Entity{Pos: pos, Owner: -1}
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
				newPos = &Entity{Pos: pos, Owner: -1}
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
				newPos = &Entity{Pos: pos, Owner: -1}
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
		do(e, true)
	}
	if len(s.freePos) == 0 {
		for _, e := range s.myEntities {
			do(e, true)
		}
		// DebugMsg(">> len free", len(s.freePos))
	}

	return s.freePos
}

func (s *State) GetFreePosToAttack() []*Entity {
	hash := make(map[string]*Entity, 0)
	s.freePos = make([]*Entity, 0)
	do := func(e *Entity) {
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
			if newPos == nil {
				newPos = &Entity{Pos: pos, Owner: -1}
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
		do(e)
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
		ee := s.getByPos(pos)
		if ee != nil && ee.IsHarvester() && ee.IsMy() {
			havsterEatPos := ee.Pos.Shift(ee.OrganDir)
			if e.Pos.Equal(havsterEatPos) {
				return true
			}
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
		if !s.InMatrix(pos) {
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

func (s *State) NearOppoent(e *Entity) bool {
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
			if e.IsOpponent() && e.IsTentacle() {
				attackPos := e.TentacleAttackPosition()
				if e.Pos.Equal(attackPos) {
					full = append(full, struct{}{})
				}
			}
		}
	}
	return len(full) > 0
}

func (s *State) GetRootWithFreePos() []*Entity {
	result := make([]*Entity, 0)
	for _, entities := range s.myEntitiesByRoot {
		if len(entities) == 0 {
			continue
		}
		i := 0
		findRoot := false
		hasFreeSpace := false
		for !findRoot || !hasFreeSpace {
			entity := entities[i]
			if entity.IsRoot() {
				findRoot = true
			}
			dirs := entity.Pos.GetRoseLocality()
			for _, dir := range dirs {
				if !s.InMatrix(dir) {
					continue
				}
				e := s.getByPos(dir)
				if e == nil {
					hasFreeSpace = true
				} else if e != nil && e.IsProtein() {
					hasFreeSpace = true
				}
			}
			i += 1
			if i == len(entities) {
				break
			}
		}
		if findRoot && hasFreeSpace {
			result = append(result, entities...)
		}
	}

	return result
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

type Stock struct {
	A        int
	PercentA float64

	B        int
	PercentB float64

	C        int
	PercentC float64

	D        int
	PercentD float64

	APerStep int
	BPerStep int
	CPerStep int
	DPerStep int
}

type SigleProtein struct {
	Type  string
	Count int
}

func (s *Stock) Scan() {
	fmt.Scan(&s.A, &s.B, &s.C, &s.D)
	total := float64(s.A + s.B + s.C + s.D)
	if total > 0.0 {
		s.PercentA = 1.0 - float64(s.A)/total
		s.PercentB = 1.0 - float64(s.B)/total
		s.PercentC = 1.0 - float64(s.C)/total
		s.PercentD = 1.0 - float64(s.D)/total
	}
}

func (s *Stock) GetPercent(protein string) float64 {
	if protein == AProteinTypeEntity {
		return s.PercentA
	}
	if protein == BProteinTypeEntity {
		return s.PercentB
	}
	if protein == CProteinTypeEntity {
		return s.PercentC
	}
	if protein == DProteinTypeEntity {
		return s.PercentD
	}
	return 0.0
}

func (s *Stock) GetProduction(protein string) int {
	if protein == AProteinTypeEntity {
		return s.APerStep
	}
	if protein == BProteinTypeEntity {
		return s.BPerStep
	}
	if protein == CProteinTypeEntity {
		return s.CPerStep
	}
	if protein == DProteinTypeEntity {
		return s.DPerStep
	}
	return -1
}

func (s *Stock) IncByType(protein string) int {
	if protein == AProteinTypeEntity {
		s.APerStep += 1
		return s.APerStep
	}
	if protein == BProteinTypeEntity {
		s.BPerStep += 1
		return s.BPerStep
	}
	if protein == CProteinTypeEntity {
		s.CPerStep += 1
		return s.CPerStep
	}
	if protein == DProteinTypeEntity {
		s.DPerStep += 1
		return s.DPerStep
	}
	return 0.0
}

func (s *Stock) NeedCollectProtein(protein string) bool {
	if protein == BProteinTypeEntity {
		return s.BPerStep == 0
	}
	if protein == CProteinTypeEntity {
		return s.CPerStep == 0
	}
	if protein == DProteinTypeEntity {
		return s.DPerStep == 0
	}
	if protein == AProteinTypeEntity {
		return s.APerStep == 0
	}
	return false
}

func (s *Stock) GetOrderByCountAsc() []string {
	proteins := []SigleProtein{
		{Type: AProteinTypeEntity, Count: s.A*s.APerStep + 1},
		{Type: BProteinTypeEntity, Count: s.B * s.BPerStep},
		{Type: CProteinTypeEntity, Count: s.C * s.CPerStep},
		{Type: DProteinTypeEntity, Count: s.D * s.DPerStep},
	}

	sort.Slice(proteins, func(i, j int) bool {
		return proteins[i].Count < proteins[j].Count
	})
	result := make([]string, 0)
	for _, pp := range proteins {
		result = append(result, pp.Type)
	}
	return result
}

func (s *Stock) StockProduction() string {
	return fmt.Sprintf("A:%d B:%d C:%d D:%d", s.APerStep, s.BPerStep, s.CPerStep, s.DPerStep)
}

func (s *Stock) CanAttack() bool {
	return s.APerStep >= 2 && s.BPerStep >= 2 && s.CPerStep >= 2 && s.DPerStep >= 2
}

func (s *Stock) TryAttack() bool {
	return s.APerStep >= 0 && s.BPerStep >= 1 && s.CPerStep >= 1 && s.DPerStep >= 0
}

func (s *Stock) CanDefend() bool {
	return s.APerStep >= 0 && s.BPerStep >= 1 && s.CPerStep >= 1 && s.DPerStep >= 0
}

func (s *Stock) String() string {
	return fmt.Sprintf("A:%d %.2f %d B:%d %.2f %d C:%d %.2f %d D:%d %.2f %d",
		s.APerStep, s.PercentA, s.A,
		s.BPerStep, s.PercentB, s.B,
		s.CPerStep, s.PercentC, s.C,
		s.DPerStep, s.PercentD, s.D,
	)
}

func (s *Stock) Score() int {
	return s.A + s.B + s.C + s.D
}

func NewStock() *Stock {
	s := &Stock{}
	s.Scan()
	return s
}

func TestGetOrderByCount(t *testing.T) {
	testCases := []struct {
		name  string
		stock Stock
		exp   string
	}{
		{
			name: "dcba",
			stock: Stock{
				A: 10,
				B: 8,
				C: 4,
				D: 2,
			},
			exp: "DCBA",
		},
		{
			name: "abcd",
			stock: Stock{
				A: 2,
				B: 4,
				C: 8,
				D: 10,
			},
			exp: "ABCD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actArr := tc.stock.GetOrderByCountAsc()
			act := strings.Join(actArr, "")
			if act != tc.exp {
				t.Error("not equal", act, tc.exp)
			}
		})
	}

}

func DebugMsg(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
}

func PointToAngle(from, to Position) int {

	ax, ay, bx, by := float64(from.X), float64(from.Y), float64(to.X), float64(to.Y)

	res := math.Atan2(bx-ax, by-ay)
	degree := res * 180 / math.Pi

	// DebugMsg("angle", res, res*180/math.Pi)
	return int(math.Round(degree))
}

func AngleToDir(degree int) string {

	if -45 <= degree && degree <= 45 {
		return DirS
	}
	if 45 <= degree && degree <= 135 {
		return DirE
	}
	if -135 <= degree && degree <= -45 {
		return DirW
	}

	return DirN
}

//)
//
///*
//
//(0;0) (1;0) (2;0) (3;0) (4;0) (5;0)
//(0;1) (1;1) (2;1) (3;1) (4;1) (5;1)
//(0;2) (1;2) (2;2) (3;2) (4;2) (5;2)
//(0;3) (1;3) (2;3) (3;3) (4;3) (5;3)
//(0;4) (1;4) (2;4) (3;4) (4;4) (5;4)
//(0;5) (1;5) (2;5) (3;5) (4;5) (5;5)
//
//*/
//
//func TestPointToAngle(t *testing.T) {
//	testCases := []struct {
//		name string
//		from Position
//		to   Position
//		exp  int
//	}{
//		{
//			Position{0, 0},
//			Position{-3, 3},
//			45,
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			act := PointToAngle(tc.from, tc.to)
//			if act != tc.exp {
//				t.Error("not equal", act, tc.exp, tc.from.ToLog(), tc.to.ToLog())
//			}
//		})
//	}
//}

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
			if n != nil && n.IsMy() && s.FreeEntites(n) {
				if !fn(n) {
					break
				}
			}
		}
	}
}
