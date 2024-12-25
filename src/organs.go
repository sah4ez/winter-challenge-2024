package main

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
