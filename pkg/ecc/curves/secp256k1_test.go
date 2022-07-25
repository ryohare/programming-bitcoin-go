package secp256k1

import (
	"fmt"
	"testing"

	point "github.com/ryohare/programming-bitcoin-go/pkg/ecc/point"
)

func TestGetGeneratorPoint(t *testing.T) {
	p := GetGeneratorPoint()
	n := GetNonce()

	np, err := point.RMultiply(p, n)

	if err != nil {
		t.Errorf("failed to validate generator point because %s", err.Error())
	}

	fmt.Println(np)
}
