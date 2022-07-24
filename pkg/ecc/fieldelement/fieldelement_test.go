package fieldelement

import (
	"testing"
)

func TestAdd(t *testing.T) {
	fe1 := &FieldElement{
		Num:   7,
		Prime: 13,
	}
	fe2 := &FieldElement{
		Num:   7,
		Prime: 13,
	}

	fea, _ := Add(fe1, fe2)

	if fea.Num != 1 && fea.Prime != 13 {
		t.Errorf("Failed addition")
	}
}

func TestSub(t *testing.T) {
	fe1 := &FieldElement{
		Num:   7,
		Prime: 13,
	}
	fe2 := &FieldElement{
		Num:   7,
		Prime: 13,
	}

	feb, _ := Subtract(fe1, fe2)

	if feb.Num != 0 && feb.Prime != 13 {
		t.Error("Failed subtraction")
	}
}

func TestMultiply(t *testing.T) {
	fe1 := &FieldElement{
		Num:   7,
		Prime: 13,
	}
	fe2 := &FieldElement{
		Num:   7,
		Prime: 13,
	}

	fem, _ := Multiply(fe1, fe2)

	if fem.Num != 10 && fem.Prime != 13 {
		t.Error("Failed multiply")
	}
}
func TestDivide(t *testing.T) {
	fe1 := &FieldElement{
		Num:   7,
		Prime: 13,
	}
	fe2 := &FieldElement{
		Num:   7,
		Prime: 13,
	}

	Divide(fe1, fe2)

}

func TestExponentiation(t *testing.T) {
	fe1 := &FieldElement{
		Num:   7,
		Prime: 13,
	}

	fee, _ := Exponentiate(fe1, int64(3))

	if fee.Num != 5 && fee.Prime != 13 {
		t.Error("Failed exponentiation")
	}
}
