package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/block"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type Headers struct {
	BlockHeaders []*block.BlockHeader
}

const COMMAND_HEADERS Command = "headers"

func ParseHeaders(reader *bytes.Reader) (*Headers, error) {
	// first part of the stream is the number of the block headers
	// which is stored as a type varint
	// numBlocks, err := binary.ReadVarint(io.ByteReader(reader))
	numBlocks := int(utils.ReadVarIntFromBytes(reader))
	// if err != nil {
	// 	return nil, err
	// }

	// all the block headers we've parsed
	var bhs []*block.BlockHeader

	// range over the message response and parse each header information
	// from the getblocks response. After each block is a varint which holds
	// the number of transactions in the block
	for i := 0; i < numBlocks; i++ {
		b, err := block.ParseHeader(reader)
		bhs = append(bhs, b)
		if err != nil {
			return nil, err
		}

		// read the varint to get the number of transactions in the block
		numTxs, err := binary.ReadVarint(io.ByteReader(reader))
		if err != nil {
			return nil, err
		}

		// as per the notes, if we dont get 0 here, something is wrong
		if numTxs != 0 {
			return nil, fmt.Errorf("number of txs is not 0")
		}
	}

	return &Headers{BlockHeaders: bhs}, nil
}

func (h Headers) Serialize() []byte {
	return nil
}

func (h Headers) GetCommand() Command {
	return COMMAND_HEADERS
}
