package secp256k1

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

func ParseSignature(sigBin []byte) (*Signature, error) {
	reader := bytes.NewReader(sigBin)

	// check the prefix first looking for DER format 0x30
	compound, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if compound != 0x30 {
		return nil, fmt.Errorf("bad signature, invalid format")
	}

	length, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if int(length)+2 != len(sigBin) {
		return nil, fmt.Errorf("bad signature length")
	}

	// look for the first marker indicating the rvalue is next
	marker, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if marker != 0x02 {
		return nil, fmt.Errorf("bad signature, missing r marker")
	}

	// read in the r value which is a big endian and convert to a big int
	// the first byte is the length of r followed by n bytes of r in big endian
	rlength, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	lr := io.LimitReader(reader, int64(rlength))
	rBin, err := ioutil.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	r := new(big.Int)
	r.SetBytes(rBin)

	// look for the next marker
	marker, err = reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if marker != 0x02 {
		return nil, fmt.Errorf("bad signature, missing s marker")
	}

	// read in the s value which is big endian and convert to a big int
	// the first byte is the length of s followed by n bytes of s in big endian
	slength, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	lr = io.LimitReader(reader, int64(slength))
	sBin, err := ioutil.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	s := new(big.Int)
	s.SetBytes(sBin)

	// double check the lenghts of the data
	if len(sigBin) != 6+int(rlength)+int(slength) {
		return nil, fmt.Errorf("singature is too long")
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}
