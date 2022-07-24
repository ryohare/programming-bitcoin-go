package point

import (
	"fmt"
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
	/*(170,142) + (60,139)
	(47,71) + (17,56)
	(143,98) + (76,66)A*/
	var prime int64 = 223
	a := &fe.FieldElement{Num: 0, Prime: prime}
	b := &fe.FieldElement{Num: 7, Prime: prime}
	x1 := &fe.FieldElement{Num: 192, Prime: prime}
	y1 := &fe.FieldElement{Num: 105, Prime: prime}
	x2 := &fe.FieldElement{Num: 17, Prime: prime}
	y2 := &fe.FieldElement{Num: 56, Prime: prime}

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
	p3, err := Addition(p1, p2)

	if err != nil {
		t.Errorf("failed to add points %v and %v because %s", p1, p2, err.Error())
	}

	fmt.Println(p3)
}
