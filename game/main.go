package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
	// "testing"
)

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
const ClusterCenter = "✨"

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
			for _, protein := range s.proteins {
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

type Entity struct {
	Pos           Position
	Type          string
	Owner         int
	OrganID       int
	OrganDir      string
	OrganParentID int
	OrganRootID   int

	NextDistance  float64
	Score         float64
	Protein       *Entity
	CanAttack     bool
	ClusterCenter bool
	SporeTo       Position
}

func (e *Entity) Scan() {
	fmt.Scan(&e.Pos.X, &e.Pos.Y, &e.Type, &e.Owner, &e.OrganID, &e.OrganDir, &e.OrganParentID, &e.OrganRootID)
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

func NewEntity() *Entity {

	e := &Entity{}
	e.Scan()
	return e
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
		// full := false
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
	if s.MyStock.B > 0 && s.MyStock.D > 0 {
		result[SporerTypeEntity] = struct{}{}
	}
	if s.MyStock.A > 0 && s.MyStock.B > 0 &&
		s.MyStock.C > 0 && s.MyStock.D > 0 {
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

type Position struct {
	X int
	Y int
}

func (p Position) Equal(pos Position) bool {
	return p.X == pos.X && p.Y == pos.Y
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

func (from Position) Distance(to Position) float64 {
	return math.Sqrt(math.Pow(float64(to.X-from.X), 2) + math.Pow(float64(to.Y-from.Y), 2))
}

func (p Position) ToCoordinates() Coordinates {
	return Coordinates{
		float64(p.X),
		float64(p.Y),
	}
}

func (p Position) ToLog() string {
	return fmt.Sprintf("(%d:%d)", p.X, p.Y)
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
				fmt.Println(e.GrowTentacle(s.GetTentacleDir(e)))
				continue
			}
		}

		if s.HasProtein(e) && organs.HasHarvester() {
			if dir := s.GetHarvesterDir(e); dir != "" {
				fmt.Println(e.GrowHarvester(dir))
				continue
			}
		}
		if len(s.mySporer) == 0 || len(s.myRoot) > len(s.mySporer) {
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
			fmt.Println(e.GrowTentacle(s.GetTentacleDir(e)))
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

	dirs := e.Pos.GetLocality()
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

func (s *State) GetSporerDir(from *Entity, to Position) string {
	if from == nil {
		return ""
	}
	degree := PointToAngle(from.Pos, to)
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

func (s *Stock) StockProduction() string {
	return fmt.Sprintf("A:%d B:%d C:%d D:%d", s.APerStep, s.BPerStep, s.CPerStep, s.DPerStep)
}

func (s *Stock) CanAttack() bool {
	return s.APerStep >= 2 && s.BPerStep >= 2 && s.CPerStep >= 2 && s.DPerStep >= 2
}

func NewStock() *Stock {
	s := &Stock{}
	s.Scan()
	return s
}

func DebugMsg(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
}

func PointToAngle(from, to Position) int {

	ax, ay, bx, by := float64(from.X), float64(from.Y), float64(to.X), float64(to.Y)

	res := math.Atan2(bx-ax, by-ay)
	degree := res * 180 / math.Pi

	DebugMsg("angle", res, res*180/math.Pi)
	return int(math.Round(degree))
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
			if !fn(n) {
				break
			}
		}
	}
}
