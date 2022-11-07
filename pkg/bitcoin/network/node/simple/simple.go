package simple

import (
	"bytes"
	"fmt"
	"net"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/network/envelope"
	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/network/messages"
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

// Perform a handshake function with a specific node
func (n *Node) Handshake() bool {
	// start of the handshake with a version message
	version := messages.MakeVersion(n.Testnet)

	// send the version message
	n.Send(version)

	// after a version message, we expect back a version message from
	// peer as well as a verack message
	_, err := n.WaitFor(messages.COMMAND_VERACK)
	return err == nil
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
func (n *Node) WaitFor(command messages.Command) (*messages.Message, error) {

	var cmd messages.Command
	payload := []byte{}
	for {
		env, err := n.Read()
		if err != nil {
			return nil, err
		}

		// check the envelope's command
		// command is a netascii string, so convert everything to
		// a string and then command and compare
		_cmd := bytes.Trim(env.Command, "\x00")
		if messages.Command(string(_cmd)) == command {
			cmd = messages.Command(string(bytes.Trim(env.Command, "\x00")))
			payload = env.Payload
			break
		} else {
			// check for overhead messages and service them
			cmd := messages.Command(bytes.Trim(env.Command, "\x00"))

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
	case messages.COMMAND_VERACK:
		verack := messages.ParseVerAck(payload)
		var msg messages.Message
		msg = verack
		return &msg, nil
	case messages.COMMAND_HEADERS:
		// parse the headers into a headers object
		headers, err := messages.ParseHeaders(bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}

		// now cast it into its base interface
		var msg messages.Message
		msg = headers
		return &msg, nil

	default:
		return nil, fmt.Errorf("unknown command matched %s", cmd)
	}
}
