package main

import (
	"fmt"
	"math"
	"os"
)

func DebugMsg(msg ...any) {
	fmt.Fprintln(os.Stderr, msg...)
}

func PointToAngle(from, to Position) int {

	ax, ay, bx, by := float64(from.X), float64(from.Y), float64(to.X), float64(to.Y)

	res := math.Atan2(bx-ax, by-ay)
	degree := res * 180 / math.Pi

	// DebugMsg("angle", res, res*180/math.Pi)
	return int(math.Round(degree))
}

func AngleToDir(degree int) string {

	if -45 <= degree && degree <= 45 {
		return DirS
	}
	if 45 <= degree && degree <= 135 {
		return DirE
	}
	if -135 <= degree && degree <= -45 {
		return DirW
	}

	return DirN
}
