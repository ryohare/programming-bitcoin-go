package secp256k1

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type PrivateKey struct {
	Secret string
	Point  *S256Point
}

func MakePrivateKey(secret string) (*PrivateKey, error) {
	pk := &PrivateKey{}
	s := new(big.Int).SetBytes([]byte(secret))
	pk.Secret = secret
	var err error
	g := GetGeneratorPoint()
	pk.Point, err = RMultiply(*g, *s)

	if err != nil {
		return nil, err
	}

	return pk, nil
}

func (p PrivateKey) Hex() string {
	return fmt.Sprintf("0x%.64x", []byte(p.Secret))
}

func (pk PrivateKey) Sign(z *big.Int) (*Signature, error) {

	// k - 32 bytes = 256 bit K
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	k := new(big.Int).SetBytes(b)

	// // r
	// kG, err := RMultiply(*GetGeneratorPoint(), *k)
	// if err != nil {
	// 	return nil, err
	// }
	// r := kG.X.Num

	// // kInt
	// n2 := new(big.Int).Sub(GetNonce(), big.NewInt(2))
	// kInv := new(big.Int).Exp(k, n2, GetNonce())

	// // s
	// secret := new(big.Int).SetBytes([]byte(p.Secret))
	// rs := new(big.Int).Mul(r, secret)
	// rs = rs.Add(rs, z)
	// rs = rs.Mul(kInv, rs)
	// s := rs.Mod(rs, GetNonce())

	// // if s > N/2:
	// if s.Cmp(new(big.Int).Div(GetNonce(), big.NewInt(2))) > 0 {
	// 	s = new(big.Int).Sub(GetNonce(), s)
	// }

	// return &Signature{
	// 		R: r,
	// 		S: s,
	// 	},
	// 	nil

	e := new(big.Int).SetBytes([]byte(pk.Secret))
	G := GetGeneratorPoint()
	N := GetNonce()

	rPoint, err := RMultiply(*G, *k)
	if err != nil {
		return nil, err
	}
	r := rPoint.Point.X.Num

	n2 := new(big.Int).Sub(N, big.NewInt(2))

	tmp := big.NewInt(0)
	kInv := tmp.Exp(k, n2, N)
	s := new(big.Int).Mul(r, e)
	s = new(big.Int).Add(s, z)
	s = s.Mul(s, kInv)
	s = s.Mod(s, N)

	return &Signature{
			R: r,
			S: s,
		},
		nil
}
