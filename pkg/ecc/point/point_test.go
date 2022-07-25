package point

import (
	"math/big"
	"testing"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
)

func TestPoint(t *testing.T) {
	var prime int64 = 223
	a := &fe.FieldElement{Num: big.NewInt(0), Prime: big.NewInt(prime)}
	b := &fe.FieldElement{Num: big.NewInt(7), Prime: big.NewInt(prime)}
	x1 := &fe.FieldElement{Num: big.NewInt(192), Prime: big.NewInt(prime)}
	y1 := &fe.FieldElement{Num: big.NewInt(105), Prime: big.NewInt(prime)}
	x2 := &fe.FieldElement{Num: big.NewInt(17), Prime: big.NewInt(prime)}
	y2 := &fe.FieldElement{Num: big.NewInt(56), Prime: big.NewInt(prime)}
	_, err := MakePoint(a, b, x1, y1)

	if err != nil {
		t.Errorf(
			"Failed to create point p1 because %s", err.Error(),
		)
	}

	_, err = MakePoint(a, b, x2, y2)

	if err != nil {
		t.Errorf(
			"Failed to create point p2 because %s", err.Error(),
		)
	}
}

func TestAddition(t *testing.T) {
	/*
		(170,142) + (60,139)= Point(220,181)_0_7 FieldElement(223)
		(47,71) + (17,56)	= Point(215,68)_0_7 FieldElement(223)
		(143,98) + (76,66)  = Point(47,71)_0_7 FieldElement(223)
	*/
	var prime int64 = 223
	a := &fe.FieldElement{Num: big.NewInt(0), Prime: big.NewInt(prime)}
	b := &fe.FieldElement{Num: big.NewInt(7), Prime: big.NewInt(prime)}
	x1 := &fe.FieldElement{Num: big.NewInt(192), Prime: big.NewInt(prime)}
	y1 := &fe.FieldElement{Num: big.NewInt(105), Prime: big.NewInt(prime)}
	x2 := &fe.FieldElement{Num: big.NewInt(17), Prime: big.NewInt(prime)}
	y2 := &fe.FieldElement{Num: big.NewInt(56), Prime: big.NewInt(prime)}

	// (170,142) + (60,139)= Point(220,181)_0_7 FieldElement(223)
	x3 := &fe.FieldElement{Num: big.NewInt(170), Prime: big.NewInt(prime)}
	y3 := &fe.FieldElement{Num: big.NewInt(142), Prime: big.NewInt(prime)}
	x4 := &fe.FieldElement{Num: big.NewInt(60), Prime: big.NewInt(prime)}
	y4 := &fe.FieldElement{Num: big.NewInt(139), Prime: big.NewInt(prime)}

	// (47,71) + (17,56)	= Point(215,68)_0_7 FieldElement(223)
	x5 := &fe.FieldElement{Num: big.NewInt(47), Prime: big.NewInt(prime)}
	y5 := &fe.FieldElement{Num: big.NewInt(71), Prime: big.NewInt(prime)}
	x6 := &fe.FieldElement{Num: big.NewInt(17), Prime: big.NewInt(prime)}
	y6 := &fe.FieldElement{Num: big.NewInt(56), Prime: big.NewInt(prime)}

	x7 := &fe.FieldElement{Num: big.NewInt(143), Prime: big.NewInt(prime)}
	y7 := &fe.FieldElement{Num: big.NewInt(98), Prime: big.NewInt(prime)}
	x8 := &fe.FieldElement{Num: big.NewInt(76), Prime: big.NewInt(prime)}
	y8 := &fe.FieldElement{Num: big.NewInt(66), Prime: big.NewInt(prime)}

	p1 := &Point{
		A: a,
		B: b,
		X: x1,
		Y: y1,
	}

	p2 := &Point{
		A: a,
		B: b,
		X: x2,
		Y: y2,
	}

	p3 := &Point{
		A: a,
		B: b,
		X: x3,
		Y: y3,
	}
	p4 := &Point{
		A: a,
		B: b,
		X: x4,
		Y: y4,
	}

	p5 := &Point{
		A: a,
		B: b,
		X: x5,
		Y: y5,
	}
	p6 := &Point{
		A: a,
		B: b,
		X: x6,
		Y: y6,
	}

	p7 := &Point{
		A: a,
		B: b,
		X: x7,
		Y: y7,
	}
	p8 := &Point{
		A: a,
		B: b,
		X: x8,
		Y: y8,
	}

	res, err := Addition(p1, p2)

	if err != nil {
		t.Errorf("failed to add points %v and %v because %s", p1, p2, err.Error())
	}

	if res.X.Num.Cmp(big.NewInt(170)) != 0 && res.Y.Num.Cmp(big.NewInt(142)) != 0 {
		t.Errorf("failed to add points %v and %v", p1, p2)
	}

	res, err = Addition(p3, p4)

	if err != nil {
		t.Errorf("failed to add points %v and %v because %s", p1, p2, err.Error())
	}

	if res.X.Num.Cmp(big.NewInt(220)) != 0 && res.Y.Num.Cmp(big.NewInt(181)) != 0 {
		t.Errorf("failed to add points %v and %v", p1, p2)
	}

	res, err = Addition(p5, p6)

	if err != nil {
		t.Errorf("failed to add points %v and %v because %s", p1, p2, err.Error())
	}

	if res.X.Num.Cmp(big.NewInt(215)) != 0 && res.Y.Num.Cmp(big.NewInt(68)) != 0 {
		t.Errorf("failed to add points %v and %v", p1, p2)
	}

	res, err = Addition(p7, p8)

	if err != nil {
		t.Errorf("failed to add points %v and %v because %s", p1, p2, err.Error())
	}

	if res.X.Num.Cmp(big.NewInt(47)) != 0 && res.Y.Num.Cmp(big.NewInt(71)) != 0 {
		t.Errorf("failed to add points %v and %v", p1, p2)
	}
}

func TestRMultiplication(t *testing.T) {
	var prime int64 = 223
	a := &fe.FieldElement{Num: big.NewInt(0), Prime: big.NewInt(prime)}
	b := &fe.FieldElement{Num: big.NewInt(7), Prime: big.NewInt(prime)}
	x1 := &fe.FieldElement{Num: big.NewInt(47), Prime: big.NewInt(prime)}
	y1 := &fe.FieldElement{Num: big.NewInt(71), Prime: big.NewInt(prime)}
	p1 := &Point{
		X: x1,
		Y: y1,
		A: a,
		B: b,
	}

	res, _ := RMultiply(p1, big.NewInt(1))

	if res.X.Num.Cmp(big.NewInt(47)) != 0 && res.Y.Num.Cmp(big.NewInt(71)) != 0 {
		t.Errorf("failed for point %s", res.String())
	}
	res, _ = RMultiply(p1, big.NewInt(2))
	if res.X.Num.Cmp(big.NewInt(36)) != 0 && res.Y.Num.Cmp(big.NewInt(111)) != 0 {
		t.Errorf("failed for point %s", res.String())
	}
	res, _ = RMultiply(p1, big.NewInt(8))
	if res.X.Num.Cmp(big.NewInt(116)) != 0 && res.Y.Num.Cmp(big.NewInt(55)) != 0 {
		t.Errorf("failed for point %s", res.String())
	}
	res, _ = RMultiply(p1, big.NewInt(15))
	if res.X.Num.Cmp(big.NewInt(139)) != 0 && res.Y.Num.Cmp(big.NewInt(86)) != 0 {
		t.Errorf("failed for point %s", res.String())
	}

}
