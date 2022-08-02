package bitcoin

import (
	"bytes"
	"encoding/binary"
)

type Script struct {
	RawScript []byte
}

func MakeScript() *Script {
	return &Script{RawScript: []byte{0x00}}
}

func ParseScript(reader *bytes.Reader) *Script {
	script := &Script{}

	length, _ := binary.ReadUvarint(reader)

	//TEMP
	b := make([]byte, length)
	reader.Read(b)

	return script
}
