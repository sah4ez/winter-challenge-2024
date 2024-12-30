package main

//import (
//	"testing"
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
//			"45 degree",
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
