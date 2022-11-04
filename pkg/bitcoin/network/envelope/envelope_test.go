package envelope

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
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

func TestSerialize(t *testing.T) {
	msg, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")

	env, err := Parse(bytes.NewReader(msg), false)

	if err != nil {
		t.Fatalf("failed to parse the envelope %s", err.Error())
	}

	msgBytes := env.Serialize()

	fmt.Printf("%x\n", msg)
	fmt.Printf("%x\n", msgBytes)

	if !utils.CompareByteArrays(msg, msgBytes) {
		t.Fatal("failed to verify the bytes")
	}
}
