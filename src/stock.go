package main

import "fmt"

type Stock struct {
	A int
	B int
	C int
	D int
}

func (s *Stock) Scan() {
	fmt.Scan(&s.A, &s.B, &s.C, &s.D)
}

func NewStock() *Stock {
	s := &Stock{}
	s.Scan()
	return s
}
