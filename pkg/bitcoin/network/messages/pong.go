package messages

import (
	"io"
	"io/ioutil"
	"net"
)

const COMMAND_PONG = "pong"

type Pong struct {
	Nonce []byte
}

func ParsePong(socket net.Conn) (*Pong, error) {
	// pong, is just reading the 8 bytes sent to the receiver
	// and throwing it back at them
	nonce, err := ioutil.ReadAll(io.LimitReader(socket, 8))
	if err != nil {
		return nil, err
	}

	return MakePong(nonce), nil
}

func MakePong(nonce []byte) *Pong {
	return &Pong{
		Nonce: nonce,
	}
}

func (v *Pong) Serialize() []byte {
	return v.Nonce
}

func (v Pong) GetCommand() string {
	return COMMAND_PONG
}
