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

func (s Signature) Der() []byte {
	rbin := s.R.Bytes()
	sbin := s.S.Bytes()

	//lsrip extra bins
	for k, v := range rbin {
		if v != 0x00 {
			rbin = rbin[k:]
			break
		}
	}
	for k, v := range sbin {
		if v != 0x00 {
			sbin = sbin[k:]
			break
		}
	}

	fmt.Printf("\n%x\n", rbin)

	// strip out the leading binaries
	// dont think this is necessary in golang since we get the bytes array on demand

	// if rbin has a high bit, add a \x00
	if rbin[0]&0x80 != 0 {
		rbin = append([]byte{0x00}, rbin...)
	}

	var resbin []byte
	pbytes := append([]byte{0x02}, byte(len(rbin)))
	resbin = append(pbytes, rbin...)
	result := resbin

	// if sbin has a high bit, add a \x00
	if sbin[0]%0x80 != 0 {
		sbin = append([]byte{0x00}, sbin...)
	}

	var sesbin []byte
	pbytes = append([]byte{0x02}, byte(len(sbin)))
	sesbin = append(pbytes, sbin...)
	result = append(result, sesbin...)

	result = append([]byte{0x30, byte(len(result))}, result...)
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
