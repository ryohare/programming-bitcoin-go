package envelope

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSmoke(t *testing.T) {
	Make(nil, nil, true)
}

func TestParse(t *testing.T) {
	msg, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")

	env, err := Parse(bytes.NewReader(msg), false)

	if err != nil {
		t.Fatalf("failed to parse the envelope %s", err.Error())
	}

	fmt.Printf("%v\n", env)
}
