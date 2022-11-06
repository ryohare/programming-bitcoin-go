package messages

import "github.com/ryohare/programming-bitcoin-go/pkg/utils"

const COMMAND_GETHEADERS Command = "getheaders"

type GetHeaders struct {
	Version    uint32
	NumHashes  uint32
	StartBlock uint32
	EndBlock   uint32
}

func MakeGetHeaders(version, numHashes, startBlock, endBlock uint32) *GetHeaders {
	return &GetHeaders{
		Version:    version,
		NumHashes:  numHashes,
		StartBlock: startBlock,
		EndBlock:   endBlock,
	}
}

// Serializes the message for transmit over the network
func (g *GetHeaders) Serialize() []byte {
	// protocol is 4 bytes little endian
	result := utils.UInt32ToLittleEndianBytes(g.Version)

	// the number of hashes is a varint
	result = append(result, utils.IntToVarintBytes(int(g.NumHashes))...)

	// Start block is little endian
	result = append(result, utils.UInt32ToLittleEndianBytes(g.StartBlock)...)

	// End block is little endian
	result = append(result, utils.UInt32ToLittleEndianBytes(g.EndBlock)...)

	return result
}
