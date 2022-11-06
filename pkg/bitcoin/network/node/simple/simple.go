package simple

import (
	"fmt"
	"net"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/network/envelope"
	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/network/messages"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type Node struct {
	Testnet bool
	Host    string
	Port    uint16
	Socket  net.Conn
}

func MakeNode(testnet bool, host string, Port uint16) (*Node, error) {
	// determine which port to used based on testnet/mainnet
	port := uint16(8333)
	if testnet {
		port = 18333
	}

	// assuming we have a host name and it is not an IP address
	// this is the assumption we are making, we need to do a resolution
	// of the IP address
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	if len(ips) < 1 {
		return nil, fmt.Errorf("failed to resolve %s to IP address", host)
	}

	// attempt to open a socket to the remote peer
	connStr := fmt.Sprintf("%s:%d", ips[0], port)

	conn, err := net.Dial("tcp", connStr)
	if err != nil {
		return nil, err
	}

	// we not have an open TCP connection to the remote peer. We should
	// store this in the node and return the node object

	return &Node{
		Port:    port,
		Socket:  conn,
		Testnet: testnet,
	}, nil
}

// sends a network message to the the remote peer
func (n *Node) Send(msg messages.Message) error {
	// create a network envelope for the message
	env := envelope.Make([]byte(msg.GetCommand()), msg.Serialize(), n.Testnet)

	// send the envelope to the remote peer
	_, err := n.Socket.Write(env.Serialize())

	if err != nil {
		return err
	}
	return nil
}

// Returns a network envelope read from the remote peer
func (n *Node) Read() (*envelope.Envelope, error) {
	env, err := envelope.ParseSocket(n.Socket, n.Testnet)

	if err != nil {
		return nil, err
	}

	return env, nil
}

// Synchronous blocking call waiting for a particular network message
func (n *Node) WaitFor(command string) (*messages.Message, error) {

	cmd := ""
	payload := []byte{}
	for {
		env, err := n.Read()
		if err != nil {
			return nil, err
		}

		// check the envelope's command
		if utils.CompareByteArrays(env.Command, []byte(command)) {
			cmd = string(env.Command)
			payload = env.Payload
			break
		} else {
			// check for overhead messages and service them
			cmd := string(env.Command)

			if cmd == new(messages.Version).GetCommand() {
				// received a version command message, send back
				// a version ack message
				verack := &messages.VersionAck{}
				n.Send(verack)
			} else if cmd == new(messages.Ping).GetCommand() {
				// send a pong response. The payload of the message
				// envelope is the nonce value needed for the construction
				// of the correct pong message
				pong := messages.MakePong(env.Payload)
				n.Send(pong)
			}
		}
	}

	// command is the command we are waiting for, return it as the correct message type
	// parse of course back to the user.
	switch cmd {
	case messages.COMMAND_VERSION:
		var msg messages.Message

		// this should be a parse call here, but there is no impl for parsing of a
		// version message at this time from the payload, ....
		fmt.Println(payload)
		versionMsg := messages.MakeVersion(n.Testnet)
		msg = versionMsg
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown command matched %s", cmd)
	}
}
