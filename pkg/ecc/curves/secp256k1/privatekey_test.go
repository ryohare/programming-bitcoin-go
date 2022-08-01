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

func TestGetDeterministsicK(t *testing.T) {
	pk, _ := MakePrivateKey("secret")
	// "secret" == 126879297332596

	pk.GetDeterministsicK(big.NewInt(1000))
}

func TestWif(t *testing.T) {
	n1 := big.NewInt(5003)
	priv, _ := MakePrivateKeyFromBigInt(n1)
	w := priv.Wif(true, true)

	if string(w) != "cMahea7zqjxrtgAbB7LSGbcQUr1uX1ojuat9jZodMN8rFTv2sfUK" {
		t.Error("wif format does not match expected version (1)")
	}

	// 2021^5 (uncompressed, testnet)
	b2 := new(big.Int).Exp(big.NewInt(2021), big.NewInt(5), nil)
	priv, _ = MakePrivateKeyFromBigInt(b2)
	w = priv.Wif(false, true)

	if string(w) != "91avARGdfge8E4tZfYLoxeJ5sGBdNJQH4kvjpWAxgzczjbCwxic" {
		t.Error("wif format does not match expected version (2)")
	}

	// 0x54321deadbeef (compressed, mainnet)
	b3, _ := new(big.Int).SetString("54321deadbeef", 16)
	priv, _ = MakePrivateKeyFromBigInt(b3)
	w = priv.Wif(true, false)

	if string(w) != "KwDiBf89QgGbjEhKnhXJuH7LrciVrZi3qYjgiuQJv1h8Ytr2S53a" {
		t.Error("wif format does not match expected version (3)")
	}
}
