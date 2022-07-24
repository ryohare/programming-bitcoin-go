package ecc

import (
	"fmt"
	"math"
)

type Point struct {
	A int64
	B int64
	X int64
	Y int64
}

func MakePoint(a, b, x, y int64) (*Point, error) {

	// Ensure the point is on the curve
	if int64(math.Pow(float64(y), 2)) != int64(math.Pow(float64(x), 3))+a*x+b {
		return nil, fmt.Errorf("(%s,%s) is not on the curve", x, y)
	}
}

func Equal(p1, p2 Point) bool {
	if p1.A == p2.A && p1.B == p2.B && p1.X == p2.X && p1.Y == p2.Y {
		return true
	}
	return false
}
