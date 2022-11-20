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
const TwoWeeks int = 60 * 60 * 24 * 14

func GetMaxTarget() *big.Int {
	ff := big.NewInt(0xffff)
	two := big.NewInt(256)
	exp := big.NewInt(0x1d - 3)
	two = new(big.Int).Exp(two, exp, nil)
	ff = new(big.Int).Mul(ff, two)
	return ff
}

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

func MutableReorderBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
	}

	return b
}

func ImmutableReorderBytes(b []byte) []byte {
	bb := make([]byte, len(b))
	for i := 0; i < len(b)/2; i++ {
		bb[i], bb[len(b)-i-1] = b[len(b)-i-1], b[i]
	}

	return bb
}

// Convert a byte array from little endian format to big endian int format
// and return as type big.int
func ConvertLittleEndianToBigInt(b []byte) *big.Int {
	littleEndian := ImmutableReorderBytes(b)
	n := new(big.Int).SetBytes(littleEndian)
	return n
}

// Convert a big.Int into little endian by
func ConvertIntToLittleEndian(i *big.Int) []byte {
	b := ImmutableReorderBytes(i.Bytes())
	return b
}

// Reads 4 bytes as a little endian variable integer and converts to a big endian integer
func LittleEndianToVarInt(reader *bytes.Reader) int {
	littleEndian := make([]byte, 4)
	reader.Read(littleEndian)
	bigEndian, _ := binary.ReadVarint(bytes.NewReader(littleEndian))
	return int(bigEndian)
}

func ReadVarIntFromBytes(reader *bytes.Reader) uint64 {

	// read the first byte to see what the rest of the story is
	i, _ := reader.ReadByte()

	if i == 0xfd {
		// 0xfd means the next two bytes are the number
		// so read 2 bytes and convert from little endian
		return uint64(LittleEndianToShort(reader))

	} else if i == 0xfe {
		// 4 bytes little endian
		return uint64(LittleEndianToUInt32(reader))

	} else if i == 0xff {
		// 8 bytes little endian
		return uint64(LittleEndianToUInt64(reader))
	} else {
		// i is just an integer, return as is
		return uint64(i)
	}
}

// Reads 4 bytes as a little endian integer and converts to a big endian integer
func LittleEndianToShort(reader *bytes.Reader) int {
	littleEndian := make([]byte, 2)
	reader.Read(littleEndian)
	bigEndian := binary.LittleEndian.Uint16(littleEndian)
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
func LittleEndianToUInt32(reader *bytes.Reader) uint32 {
	littleEndian := make([]byte, 4)
	reader.Read(littleEndian)
	bigEndian := binary.LittleEndian.Uint32(littleEndian)
	return bigEndian
}

// Reads 8 bytes as a little endian integer and converts to a big endian integer
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
	return ImmutableReorderBytes(littleEndian)
}

func ShortToBigEndianBytes(n int16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(n))
	return b
}

func ShortToLittleEndianBytes(n int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(n))
	return b
}

// Converts a big endian int to a little endian byte array
func IntToLittleEndianBytes(n int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n))
	return b
}

func UInt8ToLittleEndianBytes(n uint8) byte {
	// there is no libary in binary that does this, so we need try something else

	// TODO - Check the endians of the machine to see if we even
	// need to do this
	// _n := bits.Reverse8(n)

	// okay, this is really kludgy but, I cant think of a better way
	// put it as a uint16, depending on the machine, the byte will
	// either be MSB or LSB. So, check which one is non 0x00 and
	// return that one
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, uint16(n))

	if b[len(b)-1] != 0x00 {
		return b[len(b)-1]
	} else {
		return b[0]
	}

}

func UInt16ToLittleEndianBytes(n uint16) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, uint16(n))
	return b
}

func UInt32ToLittleEndianBytes(n uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(n))
	return b
}

func UInt64ToLittleEndianBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(n))
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

func DecodeBase58(address string) ([]byte, error) {
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

	// convert the byte array to big endian (its little endian right now)
	b := num.Bytes()

	// check is the last 4 bytes of the address, shave them off
	checksum := b[len(b)-4:]

	// validate the checksum of the address
	toVerify := b[:len(b)-4]

	h256 := Hash256(toVerify)

	for i := range checksum {
		if checksum[i] != h256[i] {
			return nil, fmt.Errorf("failed to verify checksum")
		}
	}

	// the first byte is the network prefix and the last 4 are the checksum
	// the middle 20 are the actual 20 byte address, hash160
	return b[1 : len(b)-4], nil
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func EncodeUVarInt(i uint64) ([]byte, error) {

	//bi := big.NewInt(0x10000000000000000)
	bi := big.NewInt(0x1000000000000000)
	st := big.NewInt(16)
	ib := big.NewInt(int64(i))

	if i < 0xfd {
		return []byte{byte(i)}, nil
	} else if i < 0x10000 {
		return append([]byte{0xfd}, ShortToLittleEndianBytes(int16(i))...), nil
	} else if i < 0x100000000 {
		return append([]byte{0xfe}, IntToLittleEndianBytes(int(i))...), nil
	} else if bi.Mul(bi, st).Cmp(ib) == -1 {
		return append([]byte{0xff}, UInt64ToLittleEndianBytes(i)...), nil
	} else {
		return nil, fmt.Errorf("integer is too large %d", i)
	}
}

func H160ToP2pkhAddress(h160 []byte, testnet bool) []byte {
	prefix := byte(0x00)
	if testnet {
		prefix = 0x6f
	}
	return EncodeBase58(append([]byte{prefix}, h160...))
}

func H160ToP2shAddress(h160 []byte, testnet bool) []byte {
	prefix := byte(0x05)
	if testnet {
		prefix = 0xc4
	}
	return EncodeBase58(append([]byte{prefix}, h160...))
}

func BitsToTarget(bits []byte) *big.Int {
	exponent := bits[len(bits)-1]
	coeffecient := LittleEndianToInt(bytes.NewReader(bits[:len(bits)-1]))

	_target := new(big.Int).Exp(big.NewInt(256), big.NewInt(int64(exponent)-3), nil)
	_target = new(big.Int).Mul(_target, big.NewInt(int64(coeffecient)))
	return _target
}

// func BitsToTarget(bits []byte) uint32 {
// 	// get trhe last element which is the exponent
// 	exponent := bits[len(bits)-1]

// 	// get the coefficient which is stored in the first 3 bytes of the
// 	coeffecient := LittleEndianToUInt32(bytes.NewReader(bits[:len(bits)-2]))

// 	// return the result which is caculated as follows:
// 	return coeffecient*256 ^ (uint32(exponent) - 3)
// }

// Turns a target integer into a bits byte array
func TargetsToBits(target uint32) []byte {

	// step 1, is to convert the integer into a big endian bytes array
	rawBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(rawBytes, target)

	// kill leading 0's in the number, because the python code does
	newRawBytes := make([]byte, 1)
	for _, b := range rawBytes {
		if b != 0x00 {
			newRawBytes = append(newRawBytes, b)
		} else {
			break
		}
	}

	// bits format is a way to express large numbers succicntly
	// Supports both negative and positive numbers.
	// If the first bit in the coeffecient is a 1, the bits field
	// is interpreted as a negative number. Target is always positive
	exponent := make([]byte, 4)
	coefficient := []byte{0x00}
	if newRawBytes[0] > 0x7f {
		binary.LittleEndian.PutUint32(exponent, uint32(len(newRawBytes)+1))
		coefficient = append(coefficient, newRawBytes[:2]...)
	} else {
		// exponent is how long the number is base256
		binary.LittleEndian.PutUint32(exponent, uint32(len(newRawBytes)))

		// coefficient is the first 3 digits of the base 256 number
		coefficient = newRawBytes[:3]
	}

	// coefficient is the little endian with the exponent going last
	newBits := append(ImmutableReorderBytes(coefficient), exponent...)

	return newBits
}

func CompareByteArrays(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Calculates the new bits given a 2016-block time differential and the previous bits
func CalculateNewBits(previousBits []byte, timeDifferential int) []byte {
	newTimeDiff := timeDifferential

	// if the time differential; is greater than 8 weeks, set to 8 weeks
	if timeDifferential > TwoWeeks*4 {
		newTimeDiff = TwoWeeks * 4
	}

	// if the time differential is less than half a week, set top half a week
	if timeDifferential < TwoWeeks/4 {
		newTimeDiff = TwoWeeks / 4
	}

	// the new target is the previous target * time differential / two weeks
	// newTarget := BitsToTarget(previousBits) * newTimeDiff / TwoWeeks
	newTarget := new(big.Int).Mul(BitsToTarget(previousBits), big.NewInt(int64(newTimeDiff)))
	newTarget = newTarget.Div(newTarget, big.NewInt(int64(TwoWeeks)))

	// if the new target is bigger than the MAX_TARGET, to to MAX_TARGET
	// convert to the new target to bits
	// return TargetsToBits(uint32(newTarget))
	if newTarget.Cmp(GetMaxTarget()) == 1 {
		newTarget = GetMaxTarget()
	}

	fmt.Println(newTarget)

	return TargetToBits(newTarget)
}

// Convert a big int (32 byte) target into the corresponding 4 byte bits array
func TargetToBits(target *big.Int) []byte {
	// buff := []byte{0x00}

	// get raw bytes in big endian format
	// rawBytes := append(target.Bytes(), buff...)
	rawBytes := target.Bytes()

	// strip leading 0's because they get in the way
	// a sample target is:
	// 0000000000000000007615000000000000000000000000000000000000000000
	newRawBytes := make([]byte, 1)
	for k, v := range rawBytes {
		if v != 0x00 {
			newRawBytes = rawBytes[k:len(rawBytes)]
			break
		}
	}

	// the bits format is a way to express large numbers with constrained space
	// supporting both positive and negative numbers. If the first bit in the
	// coefficient is a 1, the bits field is interpreted as a negative number
	// however the target number itself is always positive

	// holds a single byte which is the length of the hash. Range on the exponent
	// is 0 < x < 32, so it will always be in the byte[3], the length we need.
	exponent := make([]byte, 4)

	// only 4 bytes long, but want to initialize it with a 0x00 in the front
	// to start so we can use append make the byte array
	coeffecient := make([]byte, 1)

	// first check is if the number is negative or not
	if newRawBytes[0] > 0x7f {
		// number is positive

		// exponent is the last byte of the bits array
		binary.BigEndian.PutUint32(exponent, uint32(len(newRawBytes)+1))

		// coefficient is a leading 0x00, plus the first 2 bytes of the hash target
		// coefficient was already preloaded with a 0x00
		coeffecient = append(newRawBytes[:2], 0x00)
	} else {
		// number is negative

		// exponent is the last byte of the bits array
		binary.BigEndian.PutUint32(exponent, uint32(len(newRawBytes)))

		// coefficient is now the first 3 bytes of the new target
		coeffecient = newRawBytes[:3]
		for i, j := 0, len(coeffecient)-1; i < j; i, j = i+1, j-1 {
			coeffecient[i], coeffecient[j] = coeffecient[j], coeffecient[i]
		}
	}

	// need to reorder the cofficient to be little endian, (3 bytes) then the last
	// byte is the exponent
	newBits := append(coeffecient, exponent[len(exponent)-1])

	// done
	return newBits
}

// Take the binary hashes and calculates the hash256
func MerkleParent(h1, h2 []byte) []byte {
	// concatenate the hashes together
	h := append(h1, h2...)

	// do a hash256 of the concatenated byte array
	// and return
	return Hash256(h)
}

// Takes in an array of byte arrays (tx hashes) and
// returns a list that is 1/2 the length
func MerkleParentLevel(hashes [][]byte) ([][]byte, error) {
	// make sure we have more than 1 element in the array
	if len(hashes) == 1 {
		return nil, fmt.Errorf("length of hashes is 1")
	}

	// make sure the length of the hashes is even. If the
	// array is odd length, add the last element to the list
	if len(hashes)%2 == 1 {
		// odd
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	// parent level will be an array of byte arrays
	parentLevel := make([][]byte, 0, len(hashes)/2)

	// calculate the merkle parent for each hash
	// in the ordered list, which will result in an
	// array that is 1/2 the size of the passed in one.
	for i := 0; i < len(hashes); i += 2 {
		parent := MerkleParent(hashes[i], hashes[i+1])
		parentLevel = append(parentLevel, parent)
	}

	return parentLevel, nil
}

// Takes a list of binary hashes and returns the merkle root
func MerkleRoot(hashes [][]byte) ([]byte, error) {
	currentLevel := make([][]byte, len(hashes))
	copy(currentLevel, hashes)

	// keep hashing until we have a length of 1
	for {
		if len(currentLevel) == 1 {
			break
		}
		var err error
		currentLevel, err = MerkleParentLevel(currentLevel)
		if err != nil {
			return nil, err
		}
	}

	// only element left in the array is the merkle root
	return currentLevel[0], nil
}

func IsNull(b []byte) bool {
	for _, v := range b {
		if v != 0x00 {
			return false
		}
	}

	return true
}

func BitFieldToBytes(bits []byte) ([]byte, error) {
	if len(bits)%8 != 0 {
		return nil, fmt.Errorf("bitfield is not divisible by 8")
	}
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(len(bits)/8))

	// iterate over the bits
	for i, bit := range bits {
		byteIndex, bitIndex := divmod(i, 8)

		// if the bit is set to 1 in the bit mask
		if bit == 1 {
			result[byteIndex] |= 1 << bitIndex
		}
	}

	return result, nil
}

func BytesToBitField(b []byte) []byte {
	flagBits := []byte{}

	for _, _b := range b {
		for i := 0; i < 8; i++ {
			flagBits = append(flagBits, _b&1)
			_b >>= 1
		}
	}

	return flagBits
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
