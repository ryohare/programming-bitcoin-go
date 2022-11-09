package block

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type MerkleBlock struct {
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

	// Total number of transactions in the block
	Total int

	// Variable length, tx_hashes
	TxHashes [][]byte

	// Number of hashes included in the block
	NumberOfHashes uint64

	// the flags field gives information about where the hashes go within the Merkle tree
	Flags []byte
}

// Parses a block header from a bytestream
func ParseMerkleBlock(reader *bytes.Reader) (*MerkleBlock, error) {
	mb := &MerkleBlock{}

	// the merkleblock is mostly the same format as the block header
	// so we will just copy pasta that code here for now

	// parse the version first off the raw block which is a little endian int
	// convert the bytes to an integer, interpreted as little endian
	mb.Version = utils.LittleEndianToInt(reader)

	// previous block is 32 bytes, on the block chain as little endian, so
	// they will need to be reversed here
	prevBlockBytesLittleEndian, err := ioutil.ReadAll(io.LimitReader(reader, 32))
	if err != nil {
		return nil, err
	}
	mb.PreviousBlock = utils.ImmutableReorderBytes(prevBlockBytesLittleEndian)

	// next in the stream is the merkle root which like the previous block is
	// 32 bytes stored on chain as little endian format
	merkleRootBytesLittleEndian, err := ioutil.ReadAll(io.LimitReader(reader, 32))
	if err != nil {
		return nil, err
	}
	mb.MerkleRoot = utils.ImmutableReorderBytes(merkleRootBytesLittleEndian)

	// next off the stream is the time stamp. This is 4 bytes little endian stored
	// on chain.
	mb.Timestamp = utils.LittleEndianToInt(reader)

	// next is the bits section which is just 4 bytes
	bitsBytes, err := ioutil.ReadAll(io.LimitReader(reader, 4))
	if err != nil {
		return nil, err
	}
	mb.Bits = bitsBytes

	// next is the Nonce which is 4 bytes long
	nonceBytes, err := ioutil.ReadAll(io.LimitReader(reader, 4))
	if err != nil {
		return nil, err
	}
	mb.Nonce = nonceBytes

	// parge the total number of transactions field with is 4 bytes LE
	total := utils.LittleEndianToInt(reader)
	mb.Total = total

	// Get the number of hashes included which is a varint
	numOfHashes := utils.ReadVarIntFromBytes(reader)
	mb.NumberOfHashes = numOfHashes

	// now we parse each transaction hash
	hashes := make([][]byte, 0, numOfHashes)

	// iterate over the transactions
	for i := 0; i < int(numOfHashes); i++ {
		tmp, err := ioutil.ReadAll(io.LimitReader(reader, 32))
		if err != nil {
			return nil, err
		}

		// we need to reorder the bytes to be in LE
		utils.MutableReorderBytes(tmp)
		hashes = append(hashes, tmp)
	}
	mb.TxHashes = hashes

	// lastly, we read in the flags section of the merkle block
	// which starts with a length which is a varint
	flagsLen := utils.ReadVarIntFromBytes(reader)

	// finally read in the flags byte array
	// the flags field gives information about where the hashes
	// go within the Merkle tree
	flags, err := ioutil.ReadAll(io.LimitReader(reader, int64(flagsLen)))
	if err != nil {
		return nil, err
	}
	mb.Flags = flags

	return mb, nil

}

func verifyMerkleRoot(hashes [][]byte, merkleRoot []byte) bool {
	// step 1, reorder all the hashes passed in
	reordered := make([][]byte, len(hashes))
	copy(reordered, hashes)
	for _, v := range reordered {
		utils.MutableReorderBytes(v)
	}
	root, err := utils.MerkleRoot(reordered)
	if err != nil {
		return false
	}

	// reorder the root for returning
	utils.MutableReorderBytes(root)

	return utils.CompareByteArrays(root, merkleRoot)
}

func (m *MerkleBlock) VerifyMerkleRoot() bool {
	return verifyMerkleRoot(m.TxHashes, m.MerkleRoot)
}
