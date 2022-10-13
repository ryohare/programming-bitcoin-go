package script

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
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

func TestCombine(t *testing.T) {
	zBytes, err := hex.DecodeString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d")
	if err != nil {
		t.Fatalf("failed to decode z")
	}

	sec, err := hex.DecodeString("04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	if err != nil {
		t.Fatalf("failed to decode sec")
	}

	sig, err := hex.DecodeString("3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601")
	if err != nil {
		t.Fatalf("failed to decode sig")
	}

	scriptPubKey := Script{}
	var cmds []Command
	cmds = append(cmds, Command{Bytes: sec})
	// append OP_CHECKSIG (172)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, 172)
	cmds = append(cmds, Command{Bytes: b, OpCode: true})
	scriptPubKey.Commands = cmds

	scriptSig := &Script{
		RawScript: []byte{},
		Commands:  []Command{},
	}
	scriptSig.Commands = append(scriptSig.Commands, Command{Bytes: sig})

	combinedScript := Combine(scriptPubKey, *scriptSig)

	// combined script is now
	// 1: Sig
	// 2: PubKey
	// 3: 0xac 		// OP_CHECKSIG
	if !combinedScript.Evaluate(new(big.Int).SetBytes(zBytes), 0, 0, 0) {
		t.Fatalf("evaulate failed")
	}

}
