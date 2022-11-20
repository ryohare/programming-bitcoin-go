package block

import (
	"math/big"

	"github.com/spaolacci/murmur3"
)

const BIP37_CONSTANT = 0xfba4c795

type BloomFilter struct {
	Size          int
	BitField      []byte
	FunctionCount int
	Tweak         int
}

func MakeBloomFilter(size, functionCount, tweak int) *BloomFilter {
	ret := &BloomFilter{
		Size:          size,
		FunctionCount: functionCount,
		Tweak:         tweak,
	}

	b := make([]byte, size*8)
	ret.BitField = b
	return ret
}

func (b *BloomFilter) Add(item []byte) {
	for i := 0; i < b.FunctionCount; i++ {
		// BIP0037 spec seed is i*BIP37_CONSTANT + self.tweak
		seed := uint32(i*BIP37_CONSTANT + b.Tweak)

		// get the murmur3 hash with the calculated seed
		sum := murmur3.New32WithSeed(seed).Sum32()

		// set the bit at the hash mod the bitfield size (self.size*8)
		// bit = h % (self.size * 8)
		// need to use a big int because % is not defined for uint32
		tmp := new(big.Int).SetInt64(int64(b.Size * 8))
		tmp.Mod(big.NewInt(int64(sum)), tmp)
		bit := uint32(tmp.Uint64())
		// bit := int(sum) % (b.Size * 8)

		// self.bit_field[bit] = 1
		b.BitField[bit] = 1
	}
}
