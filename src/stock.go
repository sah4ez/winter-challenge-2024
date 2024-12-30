package main

import "fmt"

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

func (s *Stock) NeedCollectProtein(protein string) bool {
	if protein == AProteinTypeEntity {
		return s.APerStep == 0 && s.A == 0
	}
	if protein == BProteinTypeEntity {
		return s.BPerStep == 0 && s.B == 0
	}
	if protein == CProteinTypeEntity {
		return s.CPerStep == 0 && s.C == 0
	}
	if protein == DProteinTypeEntity {
		return s.DPerStep == 0 && s.D == 0
	}
	return false
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
