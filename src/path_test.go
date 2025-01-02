package main

import "testing"

func TestPathScore(t *testing.T) {
	testCases := []struct {
		name    string
		s       State
		from    Position
		to      Position
		exp     float64
		expFind bool
		fillFn  func(*State)
	}{
		{
			name: "one step right",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(0, 1),
			exp:     2.0,
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
			exp:     4.0,
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
			exp:     3.0,
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
			exp:     5.0,
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
			exp:     7.0,
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
			exp:     19.0,
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
			exp:     28.0,
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
			exp:     6.0,
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
			exp:     5.0,
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
	}

	for _, tc := range testCases {
		t.Run(tc.name+":"+tc.from.ToLog()+"->"+tc.to.ToLog(), func(t *testing.T) {
			tc.s.initMatrix()
			if tc.fillFn != nil {
				tc.fillFn(&tc.s)
			}
			act, find := tc.s.PathScore(tc.from, tc.to)
			if act != tc.exp {
				t.Error("unexpected score", act, tc.exp)
			}
			if find != tc.expFind {
				t.Error("unexpected find result", find, tc.expFind)
			}
		})
	}
}
