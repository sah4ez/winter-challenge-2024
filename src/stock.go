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

func NewStock() *Stock {
	s := &Stock{}
	s.Scan()
	return s
}
