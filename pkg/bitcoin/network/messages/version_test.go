package messages

import (
	"net"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/network/envelope"
)

func TestSerializeSmoke(t *testing.T) {
	v := MakeVersion(false)
	envelope.Make([]byte(COMMAND_VERSION), v.Serialize(), true)
}

func TestSerialize(t *testing.T) {

	// peers we are using
	host := "testnet.programmingbitcoin.com:18333"
	// port := uint16(18333)

	// connect to peer
	conn, err := net.Dial("tcp", host)
	if err != nil {
		t.Fatalf("failed to connect to %s because %s", host, err.Error())
	}
	defer conn.Close()

	// make the version message
	version := MakeVersion(true)

	// make the network envelope and load up the version message
	env := envelope.Make([]byte(COMMAND_VERSION), version.Serialize(), true)

	// write the message to the wire
	conn.Write(env.Serialize())

	// taking this out because verifying a network interaction is not
	// in the unit test spirit
	// expect 2 responses, a verack and a version from the peer
	// resp := make([]byte, 1014)
	// reader :=
	// _, err = conn.Read(resp)
	// if err != nil {
	// 	t.Fatalf("failed to read the version message becase %s", err.Error())
	// }
	// fmt.Printf("%x\n", resp)

	// _, err = conn.Read(resp)
	// if err != nil {
	// 	t.Fatalf("failed to read the verack message becase %s", err.Error())
	// }
	// fmt.Printf("%x\n", resp)
}
