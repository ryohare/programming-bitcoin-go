package point

import (
	"testing"
)

func TestPoint(t *testing.T) {
	_, err := MakePoint(5, 7, -1, -1)

	if err != nil {
		t.Errorf("failed to create point p1 because %s", err.Error())
	}

	_, err = MakePoint(5, 7, -1, -2)

	if err == nil {
		t.Error("failed to validate point p2 is NOT on curve.")
	}
}

func TestAddition(t *testing.T) {
	//For the curve __y__^2^ = __x__^3^ + 5__x__ + 7, what is (2,5) + (–1,–1)?
	//y2 = x3 + ax + b

	p1 := &Point{
		A: 5,
		B: 7,
		X: 2,
		Y: 5,
	}
	p2 := &Point{
		A: 5,
		B: 7,
		X: -1,
		Y: -1,
	}

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
