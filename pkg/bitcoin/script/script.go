package script

import (
	"bytes"
	"encoding/binary"
)

type Script struct {
	RawScript []byte
}

func Make() *Script {
	return &Script{RawScript: []byte{0x00}}
}

func Parse(reader *bytes.Reader) *Script {
	script := &Script{}

	length, _ := binary.ReadUvarint(reader)

	//TEMP
	b := make([]byte, length)
	reader.Read(b)

	script.RawScript = b

	return script
}

func (s Script) Serialize() []byte {
	return s.RawScript
}
