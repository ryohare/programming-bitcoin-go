package bitcoin

import "bytes"

type Script struct {
	RawScript []byte
}

func MakeScript() *Script {
	return &Script{RawScript: []byte{0x00}}
}

func Parse(reader *bytes.Reader) {

}
