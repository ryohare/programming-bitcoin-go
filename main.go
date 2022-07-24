package main

import (
	"fmt"

	ecc "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
)

func main() {
	fe1 := &ecc.FieldElement{
		Num:   7,
		Prime: 13,
	}
	fe2 := &ecc.FieldElement{
		Num:   7,
		Prime: 13,
	}

	fea, _ := ecc.Add(fe1, fe2)

	fmt.Println(fea)

	feb, _ := ecc.Subtract(fe1, fe2)

	fmt.Println(feb)

	fem, _ := ecc.Multiply(fe1, fe2)

	fmt.Println(fem)

	fee, _ := ecc.Exponentiate(fe1, int64(3))

	fmt.Println(fee)

}
