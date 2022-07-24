package point

import (
	"testing"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
)

func TestPoint(t *testing.T) {
	var prime int64 = 223
	a := &fe.FieldElement{Num: 0, Prime: prime}
	b := &fe.FieldElement{Num: 7, Prime: prime}
	x1 := &fe.FieldElement{Num: 192, Prime: prime}
	y1 := &fe.FieldElement{Num: 105, Prime: prime}
	x2 := &fe.FieldElement{Num: 17, Prime: prime}
	y2 := &fe.FieldElement{Num: 56, Prime: prime}
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
	//For the curve __y__^2^ = __x__^3^ + 5__x__ + 7, what is (2,5) + (–1,–1)?
	//y2 = x3 + ax + b
	var prime int64 = 223
	a := &fe.FieldElement{Num: 0, Prime: prime}
	b := &fe.FieldElement{Num: 7, Prime: prime}
	x1 := &fe.FieldElement{Num: 192, Prime: prime}
	y1 := &fe.FieldElement{Num: 105, Prime: prime}
	x2 := &fe.FieldElement{Num: 17, Prime: prime}
	y2 := &fe.FieldElement{Num: 56, Prime: prime}

	res, err := Addition(p1, p2)

	if err != nil {
		t.Errorf("failed addition because %s", err.Error())
	}

	if res.X != 3 && res.Y != -7 {
		t.Errorf("failed at (2,e) + (-1,-1) (%v)", res)
	}

	// for the curve __y__^2^ = __x__^3^ + 5__x__ + 7, what is (–1,–1) + (–1,–1)?
	res, err = Addition(p2, p2)

	if err != nil {
		t.Errorf("failed addition because %s", err.Error())
	}

	if res.X != 18 && res.Y != 77 {
		t.Errorf("failed at (-1,-1) + (-1,-1) (%v)", res)
	}
}
