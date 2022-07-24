package main

import (
	"fmt"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
)

func main() {

	//
	// Field element unit tests
	//

	fe1 := &fe.FieldElement{
		Num:   7,
		Prime: 13,
	}
	fe2 := &fe.FieldElement{
		Num:   7,
		Prime: 13,
	}

	fea, _ := fe.Add(fe1, fe2)

	fmt.Println(fea)

	feb, _ := fe.Subtract(fe1, fe2)

	fmt.Println(feb)

	fem, _ := fe.Multiply(fe1, fe2)

	fmt.Println(fem)

	fee, _ := fe.Exponentiate(fe1, int64(3))

	fmt.Println(fee)

}
