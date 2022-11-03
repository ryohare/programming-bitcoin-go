package envelope

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

const (
	MAINNET_PORT = 8333
	TESTNET_PORT = 18333
	SIGNET_PORT  = 38333
)

var TESTNET_NETWORK_MAGIC = [...]byte{0x0b, 0x11, 0x09, 0x07}
var MAINNET_NETWORK_MAGIC = [...]byte{0xf9, 0xbe, 0xb4, 0xd9}

type Evenvelope struct {
	Command []byte
	Payload []byte
	Magic   []byte
}

func Make(cmd, payload []byte, testnet bool) *Evenvelope {
	env := &Evenvelope{
		Command: cmd,
		Payload: payload,
	}

	if testnet {
		env.Magic = TESTNET_NETWORK_MAGIC[:]
	} else {
		env.Magic = MAINNET_NETWORK_MAGIC[:]
	}

	return env
}

func Parse(reader *bytes.Reader, testnet bool) (*Evenvelope, error) {
	// Read the first 4 bytes of the stream as the stream magic
	magicBytes, err := ioutil.ReadAll(io.LimitReader(reader, 4))
	if err != nil {
		return nil, err
	}

	if len(magicBytes) == 0 {
		return nil, fmt.Errorf("connection reset by peer via bitcoin network")
	}

	var expectedNetwork []byte
	if testnet {
		expectedNetwork = TESTNET_NETWORK_MAGIC[:]
	} else {
		expectedNetwork = MAINNET_NETWORK_MAGIC[:]
	}

	if !utils.CompareByteArrays(magicBytes, expectedNetwork) {
		return nil, fmt.Errorf("magic is not correct, read %x received %x", magicBytes, expectedNetwork)
	}

	// read in the command section which is the next 12 bytes of the stream
	commandBytes, err := ioutil.ReadAll(io.LimitReader(reader, 12))
	if err != nil {
		return nil, err
	}

	// Read the payload section now, which is variable size lengths
	payloadLength, err := ioutil.ReadAll(io.LimitReader(reader, 4))
	if err != nil {
		return nil, err
	}
	// the bytes are ordered little endian, so read it as a little endian int
	length := binary.LittleEndian.Uint32(payloadLength)

	// next 4 bytes are the checksum
	checksum, err := ioutil.ReadAll(io.LimitReader(reader, 4))
	if err != nil {
		return nil, err
	}

	// Now, read the payload given we have a length to read
	payloadBytes, err := ioutil.ReadAll(io.LimitReader(reader, int64(length)))
	if err != nil {
		return nil, err
	}

	// check the checksum, which is the first 4 bytes of the hash256 of the payload
	calculatedChecksumBytes := utils.Hash256(payloadBytes)[:4]

	if !utils.CompareByteArrays(calculatedChecksumBytes, checksum) {
		return nil, fmt.Errorf("checksomes do not match, %x vs %x", checksum, calculatedChecksumBytes)
	}

	return &Evenvelope{
		Command: commandBytes,
		Payload: payloadBytes,
		Magic:   magicBytes,
	}, nil
}
