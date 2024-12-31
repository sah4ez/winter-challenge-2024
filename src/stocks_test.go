package main

import (
	"strings"
	"testing"
)

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
