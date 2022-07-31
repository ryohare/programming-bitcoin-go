package utils

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const BASE58_ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Double sha256 hash
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

func EncodeBase58(s []byte) []byte {
	count := 0

	for _, c := range s {
		if c == 0 {
			count += 1
		} else {
			break
		}
	}

	num := new(big.Int).SetBytes(s)

	var prefix []byte
	var result []byte
	for i := 0; i < count; i++ {
		prefix = append(prefix, 0x01)
	}

	for {
		if num.Cmp(big.NewInt(0)) <= 0 {
			break
		}
		var mod *big.Int
		num, mod = new(big.Int).DivMod(num, big.NewInt(58), big.NewInt(58))
		b := byte(BASE58_ALPHABET[mod.Int64()])
		result = append([]byte{b}, result...)
	}

	if prefix != nil {
		return append(prefix, result...)
	} else {
		return result
	}
}

// sha256 followed by ripemd160
func Hash160(s []byte) []byte {
	sha := sha256.New()
	sha.Write(s)

	ripe := ripemd160.New()
	ripe.Write(sha.Sum(nil))
	return ripe.Sum(nil)
}

func EncodeBase58Checksum(b []byte) []byte {
	slice := Hash256(b)[:4]
	buf := append(b, slice...)
	return EncodeBase58(buf)
}
