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

	if np.X != nil || np.Y != nil {
		t.Errorf("failed to prove the generator point with the nonce makes a point at infinity")
	}
}

func TestSigVerify(t *testing.T) {
	z, _ := new(big.Int).SetString("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", 16)
	r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	px, _ := new(big.Int).SetString("04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574", 16)
	py, _ := new(big.Int).SetString("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4", 16)
	G := GetGeneratorPoint()
	p := MakePoint(px, py)

	// validated
	tmpSInv := new(big.Int)
	tmpN := new(big.Int)
	n2 := tmpN.Sub(GetNonce(), big.NewInt(2))
	sInv := tmpSInv.Exp(s, n2, GetNonce())
	tmpSInv = new(big.Int)

	// validated
	u := tmpSInv.Mul(z, sInv)
	u = u.Mod(u, GetNonce())
	fmt.Println(u.String())
	tmpSInv = new(big.Int)

	//verified
	v := tmpSInv.Mul(r, sInv)
	v = v.Mod(v, GetNonce())

	// verified
	uG, _ := RMultiply(*G, *u)
	vP, _ := RMultiply(*p, *v)
	sum, err := point.Addition(uG, vP)

	if err != nil {
		t.Errorf("failed addition because %s", err.Error())
	}

	if sum.X.Num.Cmp(r) != 0 {
		t.Error("failed to validate the signature")
	}

}
