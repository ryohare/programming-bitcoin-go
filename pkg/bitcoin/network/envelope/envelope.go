package envelope

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

const (
	MAINNET_PORT = 8333
	TESTNET_PORT = 18333
	SIGNET_PORT  = 38333
)

var TESTNET_NETWORK_MAGIC = [...]byte{0x0b, 0x11, 0x09, 0x07}
var MAINNET_NETWORK_MAGIC = [...]byte{0xf9, 0xbe, 0xb4, 0xd9}

type Envelope struct {
	Command []byte
	Payload []byte
	Magic   []byte
}

func Make(cmd, payload []byte, testnet bool) *Envelope {
	env := &Envelope{
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

func ParseSocket(reader net.Conn, testnet bool) (*Envelope, error) {
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

	return &Envelope{
		Command: commandBytes,
		Payload: payloadBytes,
		Magic:   magicBytes,
	}, nil
}
func Parse(reader *bytes.Reader, testnet bool) (*Envelope, error) {
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

	return &Envelope{
		Command: commandBytes,
		Payload: payloadBytes,
		Magic:   magicBytes,
	}, nil
}

func (e *Envelope) Size() int {
	return len(e.Magic) + len(e.Payload) + 12
}

func (e *Envelope) Serialize() []byte {

	// first serialize the result
	result := e.Magic

	// need to backfill/padd a command to be exactly 12 bytes
	// fresh call to make should give us 0x00'd out memory

	// next push in the command
	result = append(result, e.Command...)

	// check the length of the command and add 0x00's to the
	// byte array until the command section is exactly 12 bytes
	for i := len(e.Command); i < 12; i++ {
		result = append(result, 0x00)
	}

	// next push in the size of the payload
	receiver := make([]byte, 4)
	binary.LittleEndian.PutUint32(receiver, uint32(len(e.Payload)))
	result = append(result, receiver...)

	// add in the checksum of the payload next
	result = append(result, utils.Hash256(e.Payload)[:4]...)

	// next we push in the payload itself
	result = append(result, e.Payload...)

	return result
}
