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
	}{
		{
			name: "one step right",
			s: State{
				w: 10,
				h: 10,
			},
			from:    NewPos(0, 0),
			to:      NewPos(0, 1),
			exp:     1.0,
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
			exp:     3.0,
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
			exp:     2.0,
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
			exp:     4.0,
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
			exp:     6.0,
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
			exp:     18.0,
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
			exp:     18.0,
			expFind: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name+":"+tc.from.ToLog()+"->"+tc.to.ToLog(), func(t *testing.T) {
			tc.s.initMatrix()
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
