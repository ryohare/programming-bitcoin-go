package messages

import (
	"fmt"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

const COMMAND_GETHEADERS Command = "getheaders"

type GetHeaders struct {
	Version    uint32
	NumHashes  uint32
	StartBlock []byte
	EndBlock   []byte
}

func (h GetHeaders) GetCommand() Command {
	return COMMAND_GETHEADERS
}

func MakeGetHeaders(version, numHashes uint32, startBlock, endBlock []byte) (*GetHeaders, error) {
	// make sure we have a start block
	if startBlock == nil {
		return nil, fmt.Errorf("must specify a start block")
	}

	// next, if the endBlock is null, allocate an empty 32 byte array
	_endBlock := endBlock
	if endBlock == nil {
		_endBlock = make([]byte, 32)
	}

	return &GetHeaders{
		Version:    version,
		NumHashes:  numHashes,
		StartBlock: startBlock,
		EndBlock:   _endBlock,
	}, nil
}

// Serializes the message for transmit over the network
func (g *GetHeaders) Serialize() []byte {
	// protocol is 4 bytes little endian
	result := utils.UInt32ToLittleEndianBytes(g.Version)

	// the number of hashes is a varint
	result = append(result, utils.IntToVarintBytes(int(g.NumHashes))...)

	// Start block is little endian
	result = append(result, utils.ImmutableReorderBytes(g.StartBlock)...)

	// End block is little endian
	result = append(result, utils.ImmutableReorderBytes(g.EndBlock)...)

	return result
}
