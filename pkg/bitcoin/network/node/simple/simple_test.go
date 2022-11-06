package simple

import (
	"fmt"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/network/messages"
)

func TestSmoke(t *testing.T) {
	_, err := MakeNode(
		true,
		"testnet.programmingbitcoin.com",
		18333,
	)
	if err != nil {
		t.Fatalf("failed to create a simple node because %s", err.Error())
	}
}

func TestManualHandshake(t *testing.T) {
	node, err := MakeNode(
		true,
		"testnet.programmingbitcoin.com",
		18333,
	)
	if err != nil {
		t.Fatalf("failed to create a simple node because %s", err.Error())
	}

	// create a version message to start the handshake process
	version := messages.MakeVersion(true)

	// send the version message
	node.Send(version)

	// now we should expect back a verack and a corresponding
	// version message from the peer.
	env, err := node.Read()
	if err != nil {
		t.Fatalf("failed to read the network envelope because %s", err.Error())
	}
	fmt.Printf("%s\n", string(env.Command))
	env, err = node.Read()
	if err != nil {
		t.Fatalf("failed to read the network envelope because %s", err.Error())
	}
	fmt.Printf("%s\n", string(env.Command))
}

func TestHandshake(t *testing.T) {
	node, err := MakeNode(
		true,
		"testnet.programmingbitcoin.com",
		18333,
	)
	if err != nil {
		t.Fatalf("failed to create a simple node because %s", err.Error())
	}

	if !node.Handshake() {
		t.Fatalf("failed to handshake with the specified peer")
	}
}
