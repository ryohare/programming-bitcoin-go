package secp256k1

import (
	"fmt"
	"math/big"
	"testing"

	point "github.com/ryohare/programming-bitcoin-go/pkg/ecc/point"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
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

	np, err := point.RMultiply(pi.Point, *n)

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
	tmpSInv = new(big.Int)

	//verified
	v := tmpSInv.Mul(r, sInv)
	v = v.Mod(v, GetNonce())

	// verified
	uG, _ := RMultiply(*G, *u)
	vP, _ := RMultiply(*p, *v)
	sum, err := point.Addition(uG.Point, vP.Point)

	if err != nil {
		t.Errorf("failed addition because %s", err.Error())
	}

	if sum.X.Num.Cmp(r) != 0 {
		t.Error("failed to validate the signature")
	}
}

func TestSigCreate(t *testing.T) {
	e := new(big.Int).SetBytes(utils.Hash256([]byte("my secret")))
	z := new(big.Int).SetBytes(utils.Hash256([]byte("my message")))
	k := big.NewInt(int64(1234567890))
	G := GetGeneratorPoint()
	N := GetNonce()

	rPoint, err := RMultiply(*G, *k)
	if err != nil {
		t.Errorf("failed RMultiply because %s", err.Error())
	}
	r := rPoint.Point.X.Num

	n2 := new(big.Int).Sub(N, big.NewInt(2))

	tmp := big.NewInt(0)
	kInv := tmp.Exp(k, n2, N)
	s := new(big.Int).Mul(r, e)
	s = new(big.Int).Add(s, z)
	s = s.Mul(s, kInv)
	s = s.Mod(s, N)

	expectedValue, _ := new(big.Int).SetString("84619427107180774700812105800546110854811249640081541635353684743141289004217", 10)
	if s.Cmp(expectedValue) != 0 {
		t.Error("s value does not match the expected value")
	}
}

func TestGlobalSigVerify(t *testing.T) {
	z, _ := new(big.Int).SetString("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", 16)
	r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	px, _ := new(big.Int).SetString("04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574", 16)
	py, _ := new(big.Int).SetString("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4", 16)

	p := MakePoint(px, py)

	pk := PrivateKey{
		Secret: "secret",
		Point:  p,
	}

	sig := &Signature{
		R: r,
		S: s,
	}

	res, err := VerifySignature(pk, z, sig)

	if err != nil {
		t.Errorf("failed to verify signature because %s", err.Error())
	}

	if res != true {
		t.Errorf("signature failed verification")
	}

}

func TestSec(t *testing.T) {
	a1, _ := new(big.Int).SetString("04ffe558e388852f0120e46af2d1b370f85854a8eb0841811ece0e3e03d282d57c315dc72890a4f10a1481c031b03b351b0dc79901ca18a00cf009dbdb157a1d10", 16)
	priv, _ := MakePrivateKeyFromBigInt(big.NewInt(5000))
	c1 := new(big.Int).SetBytes(priv.Point.Sec())
	if a1.Cmp(c1) != 0 {
		t.Error("sec keys are not the same")
	}
	fmt.Printf("%x\n", priv.Point.Sec())

	a2, _ := new(big.Int).SetString("042f01e5e15cca351daff3843fb70f3c2f0a1bdd05e5af888a67784ef3e10a2a015c4da8a741539949293d082a132d13b4c2e213d6ba5b7617b5da2cb76cbde904", 16)
	priv, _ = MakePrivateKeyFromBigInt(new(big.Int).Exp(big.NewInt(2018), big.NewInt(5), big.NewInt(10)))
	c2 := new(big.Int).SetBytes(priv.Point.Sec())
	if a2.Cmp(c2) != 0 {
		t.Error("sec keys are not the same")
	}
	fmt.Printf("%x\n", priv.Point.Sec())
}
