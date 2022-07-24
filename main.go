package main

import (
	"fmt"

	ecc "github.com/ryohare/programming-bitcoin-go/pkg/ecc"
)

func main() {
	fe := &ecc.FieldElement{
		Num:   7,
		Prime: 13,
	}

	fmt.Println(fe)
}
