package secp256k1

import (
	"math/big"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
	point "github.com/ryohare/programming-bitcoin-go/pkg/ecc/point"
)

// Default vaules for the Secp256k1 Curve
// all base 16 except for the P
const N = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
const GX = "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
const GY = "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
const P = "115792089237316195423570985008687907853269984665640564039457584007908834671663" //base10

// A/B for the Secp256k1 Curve
const A = 0
const B = 7

type S256Point struct {
	Point *point.Point
}

func RMultiply(p S256Point, coefficient big.Int) (*S256Point, error) {
	_c := new(big.Int).Set(&coefficient)
	coef := _c.Mod(&coefficient, GetNonce())
	point, err := point.RMultiply(p.Point, *coef)

	if err != nil {
		return nil, err
	}

	return &S256Point{
		Point: point,
	}, nil
}

func MakePoint(x, y *big.Int) *S256Point {
	p, _ := new(big.Int).SetString(P, 10)
	point := &point.Point{
		A: &fe.FieldElement{
			Num:   big.NewInt(A),
			Prime: p,
		},
		B: &fe.FieldElement{
			Num:   big.NewInt(B),
			Prime: p,
		},
		X: &fe.FieldElement{
			Num:   x,
			Prime: p,
		},
		Y: &fe.FieldElement{
			Num:   y,
			Prime: p,
		},
	}

	return &S256Point{Point: point}
}

func GetGx() *big.Int {
	gx, _ := new(big.Int).SetString(GX, 16)
	return gx
}

func GetGy() *big.Int {
	gy, _ := new(big.Int).SetString(GY, 16)
	return gy
}

func GetPrime() *big.Int {
	p, _ := new(big.Int).SetString(P, 10)
	return p
}

func GetGeneratorPoint() *S256Point {
	gx, _ := new(big.Int).SetString(GX, 16)
	gy, _ := new(big.Int).SetString(GY, 16)
	p, _ := new(big.Int).SetString(P, 10)
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

	return &S256Point{Point: point}
}

func GetNonce() *big.Int {
	n, _ := new(big.Int).SetString(N, 16)
	return n
}
