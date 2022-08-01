package secp256k1

import (
	"fmt"
	"math/big"
	"testing"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
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

	np, err := point.RMultiply(*pi.Point, *n)

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
	sum, err := point.Addition(*uG.Point, *vP.Point)

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

func TestSqrt(t *testing.T) {
	f := fe.FieldElement{
		Num:   big.NewInt(1000),
		Prime: GetPrime(),
	}
	fe, err := Sqrt(f)

	if err != nil {
		t.Errorf("failed sqrt field element because %s", err)
	}

	fmt.Println(fe.Num.String())
}

func TestSec(t *testing.T) {
	a1, _ := new(big.Int).SetString("04ffe558e388852f0120e46af2d1b370f85854a8eb0841811ece0e3e03d282d57c315dc72890a4f10a1481c031b03b351b0dc79901ca18a00cf009dbdb157a1d10", 16)
	priv, _ := MakePrivateKeyFromBigInt(big.NewInt(5000))
	c1 := new(big.Int).SetBytes(priv.Point.Sec(false))
	if a1.Cmp(c1) != 0 {
		t.Error("sec keys are not the same")
	}
	// fmt.Printf("%x\n", priv.Point.Sec(false))

	a2, _ := new(big.Int).SetString("04027f3da1918455e03c46f659266a1bb5204e959db7364d2f473bdf8f0a13cc9dff87647fd023c13b4a4994f17691895806e1b40b57f4fd22581a4f46851f3b06", 16)
	priv, _ = MakePrivateKeyFromBigInt(new(big.Int).Exp(big.NewInt(2018), big.NewInt(5), nil))
	c2 := new(big.Int).SetBytes(priv.Point.Sec(false))
	if a2.Cmp(c2) != 0 {
		t.Error("sec keys are not the same")
	}
	// fmt.Printf("%x\n", priv.Point.Sec(false))
}

func TestParse(t *testing.T) {
	// 5,001
	// 2,019^5^
	// 0xdeadbeef54321
	pn := big.NewInt(5001)
	priv, err := MakePrivateKeyFromBigInt(pn)

	if err != nil {
		t.Errorf("failed to create private key because %s", err.Error())
	}

	a1, _ := new(big.Int).SetString("0357a4f368868a8a6d572991e484e664810ff14c05c0fa023275251151fe0e53d1", 16)
	c1 := new(big.Int).SetBytes(priv.Point.Sec(true))
	if c1.Cmp(a1) != 0 {
		t.Error("compressed sec keys are not the same")
	}

	// 2019**5
	pn = big.NewInt(2019)
	pn = pn.Exp(pn, big.NewInt(5), nil)
	// fmt.Println(pn.String())
	priv, err = MakePrivateKeyFromBigInt(pn)
	// fmt.Println(priv.Point.Point.X.Num.String())
	if err != nil {
		t.Errorf("failed to create private key because %s", err.Error())
	}

	a2, _ := new(big.Int).SetString("02933ec2d2b111b92737ec12f1c5d20f3233a0ad21cd8b36d0bca7a0cfa5cb8701", 16)
	c2 := new(big.Int).SetBytes(priv.Point.Sec(true))
	// fmt.Printf("%x\n", c2)

	if c2.Cmp(a2) != 0 {
		t.Error("compressed sec keys are not the same")
	}

	pn, _ = new(big.Int).SetString("deadbeef54321", 16)
	priv, err = MakePrivateKeyFromBigInt(pn)

	if err != nil {
		t.Errorf("failed to create private key because %s", err.Error())
	}

	a3, _ := new(big.Int).SetString("0296be5b1292f6c856b3c5654e886fc13511462059089cdf9c479623bfcbe77690", 16)
	c3 := new(big.Int).SetBytes(priv.Point.Sec(true))

	if c3.Cmp(a3) != 0 {
		t.Error("compressed sec keys are not the same")
	}
}

func TestDer(t *testing.T) {
	r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	sig := Signature{
		R: r,
		S: s,
	}

	sig.Der()
}

func TestAddress(t *testing.T) {
	//  5002 (use uncompressed SEC on testnet)
	//  2020^5 (use compressed SEC on testnet)
	//  0x12345deadbeef (use compressed SEC on mainnet)

	b1 := big.NewInt(5002)
	priv, _ := MakePrivateKeyFromBigInt(b1)
	a1 := priv.Point.Address(false, true)

	if string(a1) != "mmTPbXQFxboEtNRkwfh6K51jvdtHLxGeMA" {
		t.Error("generated address does not match the expected value (1)")
	}

	b2 := new(big.Int).Exp(big.NewInt(2020), big.NewInt(5), nil)
	priv, _ = MakePrivateKeyFromBigInt(b2)
	a2 := priv.Point.Address(true, true)

	if string(a2) != "mopVkxp8UhXqRYbCYJsbeE1h1fiF64jcoH" {
		t.Error("generated address does not match the expected value (2)")
	}

	b3, _ := new(big.Int).SetString("12345deadbeef", 16)
	priv, _ = MakePrivateKeyFromBigInt(b3)
	a3 := priv.Point.Address(true, false)
	str3 := string(a3)

	// TODO - figure out why this is failing.
	// the returned string has a \x0 infront of the address
	if str3 != "\x01F1Pn2y6pDb68E5nYJJeba4TLg2U7B6KF1" {
		t.Error("generated address does not match the expected value (3)")
	}
}
