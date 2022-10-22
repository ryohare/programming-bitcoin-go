package block

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"math/big"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

const TwoWeeks = 60 * 60 * 24 * 14

type BlockHeader struct {
	// 4 bytes - Used to encode capabilities about the miner for the mined block
	Version int

	// 32 bytes - Hash of the previous block functioning as a point
	PreviousBlock []byte

	// 32 bytes - Proof of inclusion hash to verify a specific transaction hash is included
	MerkleRoot []byte

	// 4 bytes - Timestamp when the block was mined
	Timestamp int

	// 4 Bytes - Used for validating Proof of Work (POW)
	Bits []byte

	// 4 Bytes - Randomness introduced to solve for the block difficulty by the miner
	Nonce []byte
}

// type Block struct {
// }

// // // Return the block has a byte arrasy
// // func (b *Block) Bytes() []byte {

// // }

// // // Parses a block from a supplied byte stream
// func Parse(reader *bytes.Reader) (*Block, error) {
// 	// parse the block header
// 	bh := BlockHeader{}
// }

// Parses a block header from a bytestream
func ParseHeader(reader *bytes.Reader) (*BlockHeader, error) {
	bh := &BlockHeader{}

	// parse the version first off the raw block which is a little endian int
	// convert the bytes to an integer, interpreted as little endian
	bh.Version = utils.LittleEndianToInt(reader)

	// previous block is 32 bytes, on the block chain as little endian, so
	// they will need to be reversed here
	prevBlockBytesLittleEndian, err := ioutil.ReadAll(io.LimitReader(reader, 32))
	if err != nil {
		return nil, err
	}
	bh.PreviousBlock = utils.ImmutableReorderBytes(prevBlockBytesLittleEndian)

	// next in the stream is the merkle root which like the previous block is
	// 32 bytes stored on chain as little endian format
	merkleRootBytesLittleEndian, err := ioutil.ReadAll(io.LimitReader(reader, 32))
	if err != nil {
		return nil, err
	}
	bh.MerkleRoot = utils.ImmutableReorderBytes(merkleRootBytesLittleEndian)

	// next off the stream is the time stamp. This is 4 bytes little endian stored
	// on chain.
	bh.Timestamp = utils.LittleEndianToInt(reader)

	// next is the bits section which is just 4 bytes
	bitsBytes, err := ioutil.ReadAll(io.LimitReader(reader, 4))
	if err != nil {
		return nil, err
	}
	bh.Bits = bitsBytes

	// next is the Nonce which is 4 bytes long
	nonceBytes, err := ioutil.ReadAll(io.LimitReader(reader, 4))
	if err != nil {
		return nil, err
	}
	bh.Nonce = nonceBytes

	return bh, nil
}

func (b *BlockHeader) SerializeHeader() ([]byte, error) {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(b.Version))

	// next in the byte stream does the pervious block
	result = append(result, utils.ImmutableReorderBytes(b.PreviousBlock)...)

	// Next is the merkle root
	result = append(result, utils.ImmutableReorderBytes(b.MerkleRoot)...)

	// Next is the timestamp
	result = append(result, utils.IntToLittleEndianBytes(b.Timestamp)...)

	// Next is teh bits section
	result = append(result, b.Bits...)

	// Finally is the nonce
	result = append(result, b.Nonce...)

	return result, nil
}

// Return a hash of the block header
func (b *BlockHeader) Hash() ([]byte, error) {
	s, err := b.SerializeHeader()
	if err != nil {
		return nil, err
	}
	sha256 := utils.Hash256(s)

	// reorder the bytes for return
	return utils.MutableReorderBytes(sha256), nil
}

func (b *BlockHeader) Bip9() bool {
	return b.Version>>29 == 1
}

func (b *BlockHeader) Bip91() bool {
	return b.Version>>4&1 == 1
}

func (b *BlockHeader) Bip141() bool {
	return b.Version>>1&1 == 1
}

func (b *BlockHeader) Target() *big.Int {
	return utils.BitsToTarget(b.Bits)
}
