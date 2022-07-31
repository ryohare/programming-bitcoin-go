package secp256k1

import (
	"fmt"
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

// def der(self):
// rbin = self.r.to_bytes(32, byteorder='big') # remove all null bytes at the beginning rbin = rbin.lstrip(b'\x00')
// # if rbin has a high bit, add a \x00
// if rbin[0] & 0x80:
// rbin = b'\x00' + rbin
// result = bytes([2, len(rbin)]) + rbin
// sbin = self.s.to_bytes(32, byteorder='big') # remove all null bytes at the beginning sbin = sbin.lstrip(b'\x00')
// # if sbin has a high bit, add a \x00
// if sbin[0] & 0x80:
// sbin = b'\x00' + sbin
// result += bytes([2, len(sbin)]) + sbin return bytes([0x30, len(result)]) + result
func (s Signature) Der() []byte {
	rbin := s.R.Bytes()
	sbin := s.S.Bytes()

	// strip out the leading binaries
	// dont think this is necessary in golang since we get the bytes array on deman

	// if rbin has a high bit, add a \x00
	if rbin[0]&0x80 != 0 {
		rbin = append([]byte{0x00}, rbin...)
	}

	var resbin []byte
	resbin = append([]byte{0x02}, resbin...)
	resbin = append(resbin, byte(len(rbin)))
	resbin = append(resbin, rbin...)
	result := resbin

	// if sbin has a high bit, add a \x00
	if sbin[0]&0x80 != 0 {
		sbin = append([]byte{0x00}, sbin...)
	}

	var sesbin []byte
	sesbin = append([]byte{0x02}, sesbin...)
	sesbin = append(sesbin, byte(len(sbin)))
	sesbin = append(sesbin, sbin...)
	fmt.Printf("%x\n", sesbin)
	fmt.Printf("%x\n", result)
	// result = new(big.Int).Add(new(big.Int).SetBytes(sesbin), new(big.Int).SetBytes(result)).Bytes()
	result = append(result, sesbin...)

	fmt.Printf("%x\n", result)

	result = append([]byte{0x30, byte(len(result))}, result...)
	fmt.Printf("%x\n", result)
	// sesbin = append([]byte{0x30}, result...)
	// fmt.Printf("%x\n", sesbin)
	return result
}
