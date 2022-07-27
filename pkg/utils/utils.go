package utils

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// def hash256(s):
//     '''two rounds of sha256'''
//     return hashlib.sha256(hashlib.sha256(s).digest()).digest()
// func Hash256(s big.Int) *big.Int {
// 	h := sha256.New()
// 	hh := sha256.New()
// 	h.Write(s.Bytes())
// 	hh.Write(h.Sum(nil))
// 	digest := hh.Sum(nil)

// 	return new(big.Int).SetBytes(digest)
// }

func Hash256(s []byte) []byte {
	h := sha256.New()
	hh := sha256.New()

	h.Write(s)
	hh.Write(h.Sum(nil))
	digest := hh.Sum(nil)

	return digest
}

func ToHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %x or upper case
}

func ToHexRat(n *big.Rat) string {
	return fmt.Sprintf("%x", n) // or %x or upper case
}
