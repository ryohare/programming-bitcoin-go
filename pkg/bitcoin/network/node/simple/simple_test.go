package simple

import (
	"fmt"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/block"
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

func TestGetHeaders(t *testing.T) {
	// make a node
	node, err := MakeNode(
		true,
		"testnet.programmingbitcoin.com",
		18333,
	)
	if err != nil {
		t.Fatalf("failed to create a simple node because %s", err.Error())
	}

	// handshake with the peer
	if !node.Handshake() {
		t.Fatalf("failed to handshake with the specified peer")
	}

	// start processing from the genesis block
	genesisBlock, err := block.GetTestnetGenesisBlock()
	if err != nil {
		t.Fatalf("failed to get the genesis block because %s", err.Error())
	}

	genesisBlockHash, err := genesisBlock.Hash()
	if err != nil {
		t.Fatalf("failed to get the genesis block hash because %s", err.Error())
	}

	getHeadersMessage, err := messages.MakeGetHeaders(70015, 1, genesisBlockHash, nil)
	if err != nil {
		t.Fatalf("failed to create the getheaders message because %s", err.Error())
	}

	// send the get headers message
	node.Send(getHeadersMessage)

}
