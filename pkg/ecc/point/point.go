package point

import (
	"fmt"
	"math"
	"math/big"

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

func (p Point) String() string {
	// x and y can be nil
	x := "inf"
	y := "inf"
	if p.X != nil {
		x = p.X.Num.String()
	}
	if p.Y != nil {
		y = p.Y.Num.String()
	}

	if x == "inf" {
		return "Point(infinity)"
	}

	return fmt.Sprintf(
		"Point(%s,%s)_%s_%s FieldElement(%s)",
		x,
		y,
		p.A.Num,
		p.B.Num,
		p.A.Prime,
	)
}

func Make(a, b, x, y *fe.FieldElement) (*Point, error) {

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

	y2, err := fe.Exponentiate(y, big.NewInt(2))
	if err != nil {
		return nil, fmt.Errorf("failed to exponentiate %v because %s", y, err.Error())
	}
	x3, err := fe.Exponentiate(x, big.NewInt(3))
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

	fmt.Println(y2.String())
	fmt.Println(rhs.String())

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

	// Case 0.0: self is the point at infinity, return other
	if p1.X == nil {
		return p2, nil
	}
	if p2.X == nil {
		return p1, nil
	}

	// Case 1: self.x == other.x, self.y != other.y
	// Result is point at infinity
	if p1.X == p2.X && p1.Y != p2.Y {
		return &Point{
				p1.A,
				p1.B,
				nil,
				nil,
			},
			nil
	}

	// Case 2: self.x ≠ other.x
	// Formula (x3,y3)==(x1,y1)+(x2,y2)
	// s=(y2-y1)/(x2-x1)
	// x3=s**2-x1-x2
	// y3=s*(x1-x3)-y1
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

		// x = s^2 - p1.X - p2.X
		s2, err := fe.Exponentiate(s, big.NewInt(2))
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

	// Case 4: if we are tangent to the vertical line,
	// we return the point at infinity
	// note instead of figuring out what 0 is for each type
	// we just use 0 * self.x

	// if they are the same
	if Equal(p1, p2) {

		// calculate this out here
		zeroFieldElement := &fe.FieldElement{Num: big.NewInt(0), Prime: p1.A.Prime}

		// ger the zeroith field element
		zero, err := fe.Multiply(p1.X, zeroFieldElement)
		if err != nil {
			return nil, fmt.Errorf("failed multiply because %s", err.Error())
		}

		if fe.Equal(p1.Y, zero) {
			return &Point{
					A: p1.A,
					B: p1.B,
					X: nil,
					Y: nil,
				},
				nil
		} else {
			// # Case 3: self == other
			// # Formula (x3,y3)=(x1,y1)+(x2,y2)
			// # s=(3*x1**2+a)/(2*y1)
			// # x3=s**2-2*x1
			// # y3=s*(x1-x3)-y1
			// if self == other:
			// 	s = (3 * self.x**2 + self.a) / (2 * self.y)
			// 	x = s**2 - 2 * self.x
			// 	y = s * (self.x - x) - self.y
			// 	return self.__class__(x, y, self.a, self.b)
			if Equal(p1, p2) {
				// s=(3*acc+a)/(2*y1)
				lhs := p1.X
				lhs, _ = fe.Exponentiate(lhs, big.NewInt(2))
				lhs, _ = fe.RMultiply(lhs, big.NewInt(3))
				lhs, _ = fe.Add(lhs, p1.A)
				rhs, _ := fe.RMultiply(p1.Y, big.NewInt(2))
				s, _ := fe.Divide(lhs, rhs)

				// x = s**2 - 2 * self.x
				s2, _ := fe.Exponentiate(s, big.NewInt(2))
				x2, _ := fe.RMultiply(p1.X, big.NewInt(2))
				x, _ := fe.Subtract(s2, x2)

				// y = s * (self.x - x) - self.y
				ix, _ := fe.Subtract(p1.X, x)
				sx, _ := fe.Multiply(ix, s)
				y, _ := fe.Subtract(sx, p1.Y)

				return &Point{
						A: p1.A,
						B: p2.B,
						X: x,
						Y: y,
					},
					nil
			}
		}
	}

	return nil, fmt.Errorf("failed to find addition condition which matches the two points")
}

func RMultiply(p1 *Point, coefficient *big.Int) (*Point, error) {
	product := &Point{
		p1.A,
		p1.B,
		nil,
		nil,
	}

	// bad design TODO - do the bit shift stuff
	start := big.NewInt(0)
	end := coefficient
	var one = big.NewInt(1)
	for i := new(big.Int).Set(start); i.Cmp(end) < 0; i.Add(i, one) {
		product, _ = Addition(product, p1)
	}
	return product, nil
}
