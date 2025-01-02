package main

import (
	"container/heap"
	"strings"
)

func (s *State) PathFind(start, dest *Entity) *Path {

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

	if !start.Walkable() || !dest.Walkable() {
		return nil
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
