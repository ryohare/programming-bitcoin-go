package secp256k1

import (
	"math/big"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
	point "github.com/ryohare/programming-bitcoin-go/pkg/ecc/point"
)

// Default vaules for the Secp256k1 Curve
const N = "0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
const GX = "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
const GY = "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
const P = "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"

// A/B for the Secp256k1 Curve
const A = 0
const B = 7

func GetGeneratorPoint() *point.Point {
	gx, _ := new(big.Int).SetString(GX, 16)
	gy, _ := new(big.Int).SetString(GY, 16)
	p, _ := new(big.Int).SetString(P, 16)
	a := big.NewInt(int64(A))
	b := big.NewInt(int64(B))
	point, _ := point.Make(
		&fe.FieldElement{
			Num:   a,
			Prime: p,
		},
		&fe.FieldElement{
			Num:   b,
			Prime: p,
		},
		&fe.FieldElement{
			Num:   gx,
			Prime: p,
		},
		&fe.FieldElement{
			Num:   gy,
			Prime: p,
		},
	)

	return point
}

func GetNonce() *big.Int {
	n, _ := new(big.Int).SetString(N, 16)
	return n
}
