package simple

import (
	"fmt"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/block"
	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/network/messages"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
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

func TestGetAllHeaders(t *testing.T) {
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

	// previous block is used to compare the chain against the
	// currently being processed block (e.g. the hashes match)
	previous := genesisBlock

	// used for recalculating the difficulty (e.g. the bits)
	firstEpochTimestamp := previous.Timestamp

	// used to hold what the expected bits value for a block is
	expectedBits, err := block.GetLowestBitsBytes()
	if err != nil {
		t.Fatalf("failed to get the lowest bits because %s", err.Error())
	}

	// used to count the blocks we've processed starting at the genesis block (1)
	count := 1

	for i := 0; i < 20; i++ {
		// create the get headers message with the starting block
		// being the last parsed block
		prevHash, err := previous.Hash()
		if err != nil {
			t.Fatalf("failed to get the previous blocks hash because %s", err.Error())
		}

		getHeadersMessage, err := messages.MakeGetHeaders(70015, 1, prevHash, nil)
		if err != nil {
			t.Fatalf("failed to create the getheaders message because %s", err.Error())
		}
		node.Send(getHeadersMessage)

		// now wait for the headers response messages
		msg, err := node.WaitFor(messages.COMMAND_HEADERS)
		if err != nil {
			t.Fatalf("failed to wait for the headers message because %s", err.Error())
		}

		var newMsg messages.Message
		newMsg = *msg

		// ignore type casting errors for right now
		headers, _ := newMsg.(*messages.Headers)

		// now, for each header in the response (usually 2000)
		// iterate over them and validate the transactions
		for _, header := range headers.BlockHeaders {
			// check the proof of work for the block is valid
			if !header.CheckPow() {
				t.Fatalf("proof of work is not valid for block %d", count)
			}

			// now check for continuity of the bockchain
			prevHash, err := previous.Hash()
			if err != nil {
				t.Fatalf("failed to get the previous blocks hash because %s", err.Error())
			}

			if !utils.CompareByteArrays(prevHash, header.PreviousBlock) {
				t.Logf("blockchain discontinuous at block %d", count)
			}

			fmt.Printf("%d: %x\n", count, prevHash)

			// handle difficulty adjustment which is done e very
			// 2016 blocks and is stored in the bits field of the block header
			if count%2016 == 0 {
				timeDiff := previous.Timestamp - firstEpochTimestamp
				expectedBits = utils.CalculateNewBits(previous.Bits, timeDiff)
				firstEpochTimestamp = header.Timestamp
			}

			// check that the bits match the expected bits
			if !utils.CompareByteArrays(header.Bits, expectedBits) {
				t.Logf("%d: bad bits - %s", count, prevHash)
			}

			previous = header
			count += 1
		}
	}
}
