package secp256k1

// func TestVerifySignature(t *testing.T) {
// 	z, _ := new(big.Int).SetString("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", 16)
// 	r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
// 	s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
// 	px, _ := new(big.Int).SetString("04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574", 16)
// 	py, _ := new(big.Int).SetString("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4", 16)

// 	p := MakePoint(px, py)

// 	sig := &Signature{
// 		R: r,
// 		S: s,
// 	}

// 	res, err := sig.VerifySignature(p, z, sig)

// 	if err != nil {
// 		t.Errorf("failed to verify signature because %s", err.Error())
// 	}

// 	if res != true {
// 		t.Errorf("signature failed verification")
// 	}
// }
