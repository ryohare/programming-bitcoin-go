package secp256k1

import (
	"fmt"
	"math/big"

	fe "github.com/ryohare/programming-bitcoin-go/pkg/ecc/fieldelement"
	point "github.com/ryohare/programming-bitcoin-go/pkg/ecc/point"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
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
	point, err := point.RMultiply(*p.Point, *coef)

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

func GetB() *big.Int {
	b := big.NewInt(B)
	return b
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

func VerifySignature(pk PrivateKey, z *big.Int, sig *Signature) (bool, error) {
	// validated
	G := GetGeneratorPoint()
	tmpSInv := new(big.Int)
	tmpN := new(big.Int)
	n2 := tmpN.Sub(GetNonce(), big.NewInt(2))
	sInv := tmpSInv.Exp(sig.S, n2, GetNonce())
	tmpSInv = new(big.Int)

	// validated
	u := tmpSInv.Mul(z, sInv)
	u = u.Mod(u, GetNonce())
	tmpSInv = new(big.Int)

	//verified
	v := tmpSInv.Mul(sig.R, sInv)
	v = v.Mod(v, GetNonce())

	// verified
	uG, _ := RMultiply(*G, *u)
	vP, _ := RMultiply(*pk.Point, *v)
	sum, err := point.Addition(*uG.Point, *vP.Point)

	if err != nil {
		return false, err
	}

	return sum.X.Num.Cmp(sig.R) == 0, nil
}

func (s S256Point) Sec(compressed bool) []byte {
	buf := make([]byte, 0, 32)

	fmt.Println(s.Point.Y.Num.String())
	fmt.Println(s.Point.X.Num.String())
	if compressed {
		if new(big.Int).Mod(s.Point.Y.Num, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
			buf = append(buf, 0x02)
			buf = append(buf, s.Point.X.Num.Bytes()...)
		} else {
			buf = append(buf, 0x03)
			buf = append(buf, s.Point.X.Num.Bytes()...)
		}
	} else {
		buf = append(buf, 0x04)
		buf = append(buf, s.Point.X.Num.Bytes()...)
		buf = append(buf, s.Point.Y.Num.Bytes()...)
	}
	return buf
}

func Sqrt(fe1 fe.FieldElement) (*fe.FieldElement, error) {
	p1 := new(big.Int).Add(GetPrime(), big.NewInt(1))
	p1 = p1.Div(p1, big.NewInt(4))

	res, err := fe.RMultiply(&fe1, p1)

	if err != nil {
		return nil, err
	}

	fmt.Println(res.String())

	return res, nil
}

func Parse(secBin []byte) (*S256Point, error) {
	if len(secBin) < 1 {
		return nil, fmt.Errorf("secBin is too short")
	}
	if secBin[0] == 0x04 {
		x := new(big.Int).SetBytes(secBin[1:33])
		y := new(big.Int).SetBytes(secBin[33:65])
		return MakePoint(x, y), nil
	}

	isEven := secBin[0] == 0x02
	x := &fe.FieldElement{Num: new(big.Int).SetBytes(secBin[1:]), Prime: GetPrime()}

	// # right side of the equation y^2 = x^3 + 7
	alpha, err := fe.Exponentiate(x, *big.NewInt(3))

	if err != nil {
		return nil, err
	}

	alpha, err = fe.Add(alpha, &fe.FieldElement{Num: GetB(), Prime: GetPrime()})

	if err != nil {
		return nil, err
	}

	beta, err := Sqrt(*alpha)

	if err != nil {
		return nil, err
	}

	p1 := new(big.Int).Sub(GetPrime(), beta.Num)
	var evenBeta *fe.FieldElement
	var oddBeta *fe.FieldElement
	if new(big.Int).Mod(beta.Num, big.NewInt(2)).Cmp(big.NewInt(2)) == 0 {
		evenBeta = beta
		oddBeta = &fe.FieldElement{Num: p1, Prime: GetPrime()}
	} else {
		oddBeta = beta
		evenBeta = &fe.FieldElement{Num: p1, Prime: GetPrime()}
	}

	if isEven {
		return &S256Point{
			Point: &point.Point{
				A: &fe.FieldElement{
					Num:   big.NewInt(0),
					Prime: GetPrime(),
				},
				B: &fe.FieldElement{
					Num:   big.NewInt(7),
					Prime: GetPrime(),
				},
				X: &fe.FieldElement{
					Num:   x.Num,
					Prime: GetPrime(),
				},
				Y: &fe.FieldElement{
					Num:   evenBeta.Num,
					Prime: GetPrime(),
				},
			},
		}, nil
	} else {
		return &S256Point{
			Point: &point.Point{
				A: &fe.FieldElement{
					Num:   big.NewInt(0),
					Prime: GetPrime(),
				},
				B: &fe.FieldElement{
					Num:   big.NewInt(7),
					Prime: GetPrime(),
				},
				X: &fe.FieldElement{
					Num:   x.Num,
					Prime: GetPrime(),
				},
				Y: &fe.FieldElement{
					Num:   oddBeta.Num,
					Prime: GetPrime(),
				},
			},
		}, nil
	}
}

func (s S256Point) Hash160(compressed bool) []byte {
	return utils.Hash160(s.Sec(compressed))
}

func (s S256Point) Address(compressed, testnet bool) []byte {

	fmt.Println(s.Point.X.String())
	fmt.Println(s.Point.Y.String())
	h160 := s.Hash160(compressed)

	var prefix []byte
	if testnet {
		prefix = append(prefix, 0x6f)
	} else {
		prefix = append(prefix, 0x00)
	}

	payload := append(prefix, h160...)

	return utils.EncodeBase58Checksum(payload)
}
