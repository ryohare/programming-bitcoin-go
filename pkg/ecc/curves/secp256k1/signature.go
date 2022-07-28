package secp256k1

import (
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

// func (s *Signature) VerifySignature(p *S256Point, z *big.Int, sig *Signature) (bool, error) {
// 	// validated
// 	G := GetGeneratorPoint()
// 	tmpSInv := new(big.Int)
// 	tmpN := new(big.Int)
// 	n2 := tmpN.Sub(GetNonce(), big.NewInt(2))
// 	sInv := tmpSInv.Exp(sig.S, n2, GetNonce())
// 	tmpSInv = new(big.Int)

// 	// validated
// 	u := tmpSInv.Mul(z, sInv)
// 	u = u.Mod(u, GetNonce())
// 	tmpSInv = new(big.Int)

// 	//verified
// 	v := tmpSInv.Mul(sig.R, sInv)
// 	v = v.Mod(v, GetNonce())

// 	// verified
// 	uG, _ := RMultiply(*G, *u)
// 	vP, _ := RMultiply(*p, *v)
// 	sum, err := point.Addition(uG.Point, vP.Point)

// 	if err != nil {
// 		return false, err
// 	}

// 	return sum.X.Num.Cmp(sig.R) == 0, nil
// }
