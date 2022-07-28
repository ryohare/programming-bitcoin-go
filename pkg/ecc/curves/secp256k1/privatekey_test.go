package secp256k1

import (
	"math/big"
	"testing"
)

func TestHex(t *testing.T) {
	message := "message"
	secret := new(big.Int).SetBytes([]byte(message))
	g := GetGeneratorPoint()
	p, _ := RMultiply(*g, *secret)
	pk := PrivateKey{
		Secret: "message",
		Point:  p,
	}

	pk.Hex()
}

// def test_sign(self):
// pk = PrivateKey(randint(0, N))
// z = randint(0, 2**256)
// sig = pk.sign(z)
// self.assertTrue(pk.point.verify(z, sig))
func TestSignMessage(t *testing.T) {
	pk, err := MakePrivateKey("rnd.String()")

	if err != nil {
		t.Errorf("failed to make private key because %s", err.Error())
	}

	message := new(big.Int).SetBytes([]byte("secret message"))

	if err != nil {
		t.Errorf("failed to get a random z value because %s", err.Error())
	}

	sig, err := pk.Sign(message)

	if err != nil {
		t.Errorf("failed to sign message because %s", err.Error())
	}

	status, err := VerifySignature(*pk, message, sig)

	if err != nil {
		t.Errorf("failed to verify signature because %s", err.Error())
	}

	if !status {
		t.Error("failed to verify signature")
	}
}
