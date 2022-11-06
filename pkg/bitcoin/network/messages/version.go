package messages

import (
	"encoding/hex"
	"net"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

const COMMAND_VERSION Command = "version"

func (v Version) GetCommand() Command {
	return COMMAND_VERSION
}

type Version struct {
	Version uint32

	// Capabilities of the node
	Services uint64

	// Unix Timestamp in Little Endian
	Timestamp uint64

	// 8 byte little endian number
	ReceiverServices int

	// IPv4, IPv6, OnionCat. This impl only supports IPv4
	ReceiverIp net.IP

	// Big Endian
	ReceiverPort uint16

	// 8 byte little endian number
	SenderServices int

	// IPv4, IPv6, OnionCat. This impl only supports IPv4
	SenderIp net.IP

	// Big Endian
	SenderPort uint16

	// Rnadom number used to detect self connections
	Nonce []byte

	// Identification of the software being run
	UserAgent []byte

	// Block height or latest block the node knows of
	LatestBlock uint32

	// Used for Bloomfilters in SVP nodes
	Relay bool
}

func MakeVersion(testnet bool) *Version {
	port := 8333
	if testnet {
		port = 18333
	}

	return &Version{
		Version:          70015,
		Services:         0,
		Timestamp:        0,
		ReceiverServices: 0,
		ReceiverIp:       net.ParseIP("127.0.0.1"),
		ReceiverPort:     uint16(port),
		SenderServices:   0,
		SenderIp:         net.ParseIP("127.0.0.1"),
		SenderPort:       uint16(port),
		Nonce:            []byte{},
		UserAgent:        []byte("/programmingbitcoin:0.1"),
		LatestBlock:      0,
		Relay:            false,
	}
}

func IPv4Serialization(ip net.IP) []byte {
	prefixBytes, _ := hex.DecodeString("00000000000000000000ffff")
	var ipBytes []byte = ip
	prefixBytes = append(prefixBytes, ipBytes...)
	return prefixBytes
}

func (v *Version) Serialize() []byte {
	// version is 4 bytes little endian
	result := utils.UInt32ToLittleEndianBytes(v.Version)

	// services is an 8 byte little endian integer
	result = append(result, utils.UInt64ToLittleEndianBytes(v.Services)...)

	// time stamp is 8 bytles little endian
	result = append(result, utils.UInt64ToLittleEndianBytes(v.Timestamp)...)

	// receiver services is 8 bytes little endian
	result = append(result, utils.UInt64ToLittleEndianBytes(uint64(v.ReceiverServices))...)

	// IPv4 is 10 00 bytes and 2 FF then the receiver IP
	result = append(result, IPv4Serialization(v.ReceiverIp)...)

	// Receiver port is 2 bytles big endian endian
	result = append(result, utils.ShortToBigEndianBytes(int16(v.ReceiverPort))...)

	// sender services is 8 bytes little endian
	result = append(result, utils.UInt64ToLittleEndianBytes(uint64(v.SenderServices))...)

	// IPv4 is 10 00 bytes and 2 FF then the receiver IP
	result = append(result, IPv4Serialization(v.SenderIp)...)

	// Receiver port is 2 bytles big endian endian
	result = append(result, utils.ShortToBigEndianBytes(int16(v.SenderPort))...)

	// sender services is the nonce value
	result = append(result, v.Nonce...)

	// user agent is a variable length string, starting with a varint signifying the length to read
	length, _ := utils.EncodeUVarInt(uint64(len(v.UserAgent)))
	result = append(result, length...)
	result = append(result, v.UserAgent...)

	// latest block is 4 bytes little endian
	result = append(result, utils.UInt32ToLittleEndianBytes(uint32(v.LatestBlock))...)

	// relay is 00 for false, and 01 for true
	if v.Relay {
		result = append(result, 0x01)
	} else {
		result = append(result, 0x00)
	}

	return result
}
