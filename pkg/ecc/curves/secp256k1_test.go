package secp256k1

import (
	"fmt"
	"math/big"
	"testing"

	point "github.com/ryohare/programming-bitcoin-go/pkg/ecc/point"
)

func TestGetGeneratorPoint(t *testing.T) {

	p := GetPrime()
	gy := GetGy()
	gy = gy.Exp(gy, big.NewInt(2), p)

	gx := GetGx()
	gx = gx.Exp(gx, big.NewInt(3), nil)
	gx = gx.Add(gx, big.NewInt(7))
	gx = gx.Mod(gx, p)

	if gx.Cmp(gy) != 0 {
		t.Error("gx, gy or p is off")
	}

	pi := GetGeneratorPoint()
	n := GetNonce()

	np, err := point.RMultiply(pi, *n)

	if err != nil {
		t.Errorf("failed to validate generator point because %s", err.Error())
	}

	fmt.Println(np)
}
