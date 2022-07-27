package fieldelement

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAdd(t *testing.T) {
	fe1 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}
	fe2 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}

	f, _ := Add(fe1, fe2)

	if f.Num.Cmp(big.NewInt(1)) != 0 && f.Prime.Cmp(big.NewInt(13)) != 0 {
		t.Errorf("failed adding %v and %v", fe1.String(), fe2.String())
	}
}

func TestSub(t *testing.T) {
	fe1 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}
	fe2 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}

	feb, _ := Subtract(fe1, fe2)

	fmt.Println(feb)
}

func TestMultiply(t *testing.T) {
	fe1 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}
	fe2 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}

	fem, _ := Multiply(fe1, fe2)

	fmt.Println(fem)

	if fem.Num.Cmp(big.NewInt(10)) != 0 && fem.Prime.Cmp(big.NewInt(13)) != 0 {
		t.Error("Failed multiply")
	}
}

func TestDivide(t *testing.T) {
	fe1 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}
	fe2 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}

	res, _ := Divide(fe1, fe2)

	if res.Num.Cmp(big.NewInt(1)) != 0 && res.Prime.Cmp(big.NewInt(13)) != 0 {
		t.Errorf("failed division")
	}

}

func TestExponentiation(t *testing.T) {
	fe1 := &FieldElement{
		Num:   big.NewInt(7),
		Prime: big.NewInt(13),
	}

	fee, _ := Exponentiate(fe1, *big.NewInt(3))
	fmt.Println(fee)

	if fee.Num.Cmp(big.NewInt(5)) != 0 && fee.Prime.Cmp(big.NewInt(13)) != 0 {
		t.Error("failed exponentiation")
	}

	// if fee.Num != 5 && fee.Prime != 13 {
	// 	t.Error("Failed exponentiation")
	// }
}
