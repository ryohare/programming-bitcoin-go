package point

import (
	"fmt"
	"math"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
)

// constants
const INF = math.MinInt

type Point struct {
	A *fe.FieldElement
	B *fe.FieldElement
	X *fe.FieldElement
	Y *fe.FieldElement
}

func MakePoint(a, b, x, y *fe.FieldElement) (*Point, error) {

	// Ensure the point is on the curve
	// y^2 = x^3 + ax + b

	// points at infinity check
	// inifinty is a nil pointer here
	// dont do the checks at infinity, they will break shit
	if x == nil && y == nil {
		return &Point{
				A: a,
				B: b,
				X: x,
				Y: y,
			},
			nil
	}

	y2, err := fe.Exponentiate(y, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to exponentiate %v because %s", y, err.Error())
	}
	x3, err := fe.Exponentiate(x, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to exponentiate %v because %s", x, err.Error())
	}
	ax, err := fe.Multiply(a, x)
	if err != nil {
		return nil, fmt.Errorf("failed to mulitply %v because %s", y, err.Error())
	}
	rhs, err := fe.Add(x3, ax)
	if err != nil {
		return nil, fmt.Errorf("failed to add %v because %s", y, err.Error())
	}
	rhs, err = fe.Add(rhs, b)
	if err != nil {
		return nil, fmt.Errorf("failed to add %v because %s", y, err.Error())
	}

	// check the the point is on the curve
	if !fe.Equal(rhs, y2) {
		return nil, fmt.Errorf("(%d,%d) is not on the curve", x, y)
	}

	return &Point{
			A: a,
			B: b,
			X: x,
			Y: y,
		},
		nil
}

func Equal(p1, p2 *Point) bool {
	if p1.A == p2.A && p1.B == p2.B && p1.X == p2.X && p1.Y == p2.Y {
		return true
	}
	return false
}

func NotEqual(p1, p2 *Point) bool {
	return !Equal(p1, p2)
}

func Addition(p1, p2 *Point) (*Point, error) {

	// make sure both points are on the same curve
	if p1.A != p2.A || p1.B != p2.B {
		return nil, fmt.Errorf("points are not on the same curve")
	}

	// check for points at infinity, additive identity property
	if p1.X == nil {
		return p2, nil
	}
	if p2.X == nil {
		return p1, nil
	}

	// check for a straigt y line
	if p1.X == p2.X && p1.Y != p2.Y {
		return &Point{
				p1.A,
				p1.B,
				nil,
				nil,
			},
			nil
	}

	// do the addition
	if !fe.Equal(p1.X, p2.X) {

		//(p2.Y - p1.Y) / (p2.X - p1.X)
		lhs, err := fe.Subtract(p2.Y, p1.Y)
		if err != nil {
			return nil, fmt.Errorf("failed subtraction of %v and %v", p2.Y, p1.Y)
		}
		rhs, err := fe.Subtract(p2.X, p1.X)
		if err != nil {
			return nil, fmt.Errorf("failed subtraction of %v and %v", p2.X, p1.X)
		}
		s, err := fe.Divide(lhs, rhs)
		if err != nil {
			return nil, fmt.Errorf("failed division of %v and %v", lhs, rhs)
		}
		fmt.Println(s)

		// x = s^2 - p1.X - p2.X
		s2, err := fe.Exponentiate(s, 2)
		if err != nil {
			return nil, fmt.Errorf("failed to exponentiate %v by %d", s, 2)
		}
		sub1, err := fe.Subtract(s2, p1.X)
		if err != nil {
			return nil, fmt.Errorf("failed to substract %v and %v", s2, p1.X)
		}
		x, err := fe.Subtract(sub1, p2.X)
		if err != nil {
			return nil, fmt.Errorf("failed to substract %v and %v", sub1, p2.X)
		}

		// y := s*(p1.X-x) - p1.Y
		sub1, err = fe.Subtract(p1.X, x)
		if err != nil {
			return nil, fmt.Errorf("faile to substract %v and %v", p1.X, x)
		}
		sx, err := fe.Multiply(s, sub1)
		if err != nil {
			return nil, fmt.Errorf("failed to multiply %v and %v", s, sub1)
		}
		y, err := fe.Subtract(sx, p1.Y)
		if err != nil {
			return nil, fmt.Errorf("failed to subtract %v and %v", sx, p1.Y)
		}

		return &Point{
				A: p1.A,
				B: p2.B,
				X: x,
				Y: y,
			},
			nil
	}

	// one more exception
	// if Equal(p1, p2) && p1.Y == 0*p1.X {
	// 	return &Point{
	// 			A: p1.A,
	// 			B: p1.B,
	// 			X: nil,
	// 			Y: nil,
	// 		},
	// 		nil
	// }

	// // Adding against self
	// if Equal(p1, p2) {
	// 	s := (3*int64(math.Pow(float64(p1.X), 2)) + p1.A) / (2 * p1.Y)
	// 	x := int64(math.Pow(float64(s), 2)) - 2*p1.X
	// 	y := s*(p1.X-x) - p1.Y
	// 	fmt.Println(y)
	// 	return &Point{
	// 			A: p1.A,
	// 			B: p1.B,
	// 			X: x,
	// 			Y: y,
	// 		},
	// 		nil
	// }
	return nil, fmt.Errorf("failed to find addition condition which matches the two points")
}
