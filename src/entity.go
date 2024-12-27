package main

import "fmt"

type Entity struct {
	Pos           Position
	Type          string
	Owner         int
	OrganID       int
	OrganDir      string
	OrganParentID int
	OrganRootID   int

	NextDistance float64
	Score        float64
	Protein      *Entity
	CanAttack    bool
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
	return fmt.Sprintf("%s %d %d %d", SporeCmd, e.OrganID, e.Protein.Pos.X, e.Protein.Pos.Y)
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
	return fmt.Sprintf("(%d:%d)%s:%d:%d:%.2f", e.Pos.X, e.Pos.Y, e.Type, e.OrganID, e.Owner, e.NextDistance)
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
