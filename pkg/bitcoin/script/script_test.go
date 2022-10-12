package script

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func testEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestParse(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")

	reader := bytes.NewReader(scriptPubKey)

	script, err := Parse(reader)

	if err != nil {
		t.Fatalf("failed to parse the script because %v\n", err)
	}

	if script == nil {
		t.Fatalf("parsed script is nil")
	}

	cmd1, _ := hex.DecodeString("3045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01")
	cmd2, _ := hex.DecodeString("0349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")

	if !testEq(cmd1, script.Commands[0].Bytes) {
		t.Fatalf("cmd0 does not match")
	}

	if !testEq(cmd2, script.Commands[1].Bytes) {
		t.Fatalf("cmd1 does not match")
	}
}

func TestSeralize(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")

	reader := bytes.NewReader(scriptPubKey)

	script, err := Parse(reader)

	if err != nil {
		t.Fatalf("failed to parse the script because %v\n", err)
	}

	if script == nil {
		t.Fatalf("parsed script is nil")
	}

	seralizedScript := script.Serialize()

	if len(seralizedScript) != len(scriptPubKey) {
		t.Fatalf("seralized script does not match oringal script")
	}
}
