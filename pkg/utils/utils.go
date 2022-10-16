package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
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
		prefix = append(prefix, 1)
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

func ReorderBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
	}

	return b
}

// Convert a byte array from little endian format to big endian int format
// and return as type big.int
func ConvertLittleEndianToBigInt(b []byte) *big.Int {
	littleEndian := ReorderBytes(b)
	n := new(big.Int).SetBytes(littleEndian)
	return n
}

// Convert a big.Int into little endian by
func ConvertIntToLittleEndian(i *big.Int) []byte {
	b := ReorderBytes(i.Bytes())
	return b
}

// Reads 4 bytes as a little endian variable integer and converts to a big endian integer
func LittleEndianToVarInt(reader *bytes.Reader) int {
	littleEndian := make([]byte, 4)
	reader.Read(littleEndian)
	bigEndian, _ := binary.ReadUvarint(bytes.NewReader(littleEndian))
	return int(bigEndian)
}

// Reads 4 bytes as a little endian integer and converts to a big endian integer
func LittleEndianToInt(reader *bytes.Reader) int {
	littleEndian := make([]byte, 4)
	reader.Read(littleEndian)
	bigEndian := binary.LittleEndian.Uint32(littleEndian)
	return int(bigEndian)
}

// Reads 4 bytes as a little endian integer and converts to a big endian integer
func LittleEndianToUInt64(reader *bytes.Reader) uint64 {
	littleEndian := make([]byte, 8)
	reader.Read(littleEndian)
	bigEndian := binary.LittleEndian.Uint64(littleEndian)
	return bigEndian
}

// Takes in a stream reader, reads in n bytes and reorders from
// Little Endian to Big Endian.
func LittleEndianToBigEndian(reader *bytes.Reader, length int) []byte {
	littleEndian := make([]byte, length)
	reader.Read(littleEndian)
	return ReorderBytes(littleEndian)
}

// Converts a big endian uint64 to a little endian byte array
func UInt64ToLittleEndianBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

// Converts a big endian int to a little endian byte array
func IntToLittleEndianBytes(n int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n))
	return b
}

// Takes in an int and encodes it to a var int
func IntToVarintBytes(v int) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(v))
	b := buf[:n]
	return b
}

func getIndex(s rune) int {
	for i, c := range BASE58_ALPHABET {
		if string(s) == string(c) {
			return i
		}
	}
	return -1
}

func DecodeBase58(address string) []byte {
	num := new(big.Int)
	for _, b := range address {
		i := getIndex(b)
		num = num.Mul(num, big.NewInt(58))

		if i != -1 {
			num.Add(num, big.NewInt(int64(i)))
		} else {
			fmt.Printf("rune is outside the range for base58")
		}
	}

	// combine in big endian format
	return num.Bytes()
}
