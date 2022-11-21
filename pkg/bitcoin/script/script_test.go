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
	if !combinedScript.Evaluate(new(big.Int).SetBytes(zBytes), 0, 0, 0, nil) {
		t.Fatalf("evaulate failed")
	}
}

func TestScriptOpCodes(t *testing.T) {
	// X
	// OP_DUP
	// OP_DUP
	// OP_MUL
	// OP_ADD
	// OP_6
	// OP_EQUAL
	// solving the question comes down to 6 = x^2+x
	scriptPubkeyBytes := []byte{0x76, 0x76, 0x95, 0x93, 0x56, 0x87}
	scriptPubKey := &Script{
		RawScript: scriptPubkeyBytes,
	}
	for _, b := range scriptPubkeyBytes {
		scriptPubKey.Commands = append(scriptPubKey.Commands, Command{Bytes: []byte{b}})
	}

	// solving x^2+x, x=2 >>>> op_2
	scriptSigBytes := []byte{0x52}
	scriptSig := &Script{
		RawScript: scriptPubkeyBytes,
	}
	scriptSig.Commands = append(scriptSig.Commands, Command{Bytes: scriptSigBytes})
	combinedScript := Combine(*scriptPubKey, *scriptSig)

	if !combinedScript.Evaluate(big.NewInt(0), 0, 0, 0, nil) {
		t.Fatalf("failed to evaulate the script")
	}
}

func TestCollision(t *testing.T) {
	scriptPubkeyBytes := []byte{0x6e, 0x87, 0x91, 0x69, 0xa7, 0x7c, 0xa7, 0x87}
	c1Bytes, _ := hex.DecodeString("55044462d312e330a25e2e3cfd30a0a0a312030206f626a0a3c3c2f57696474682032203020522f4865696768742033203020522f547970652034203020522f537562747970652035203020522f46696c7465722036203020522f436f6c6f7253706163652037203020522f4c656e6774682038203020522f42697473506572436f6d706f6e656e7420383e3e0a73747265616d0affd8fffe00245348412d3120697320646561642121212121852fec092339759c39b1a1c63c4c97e1fffe017f46dc93a6b67e013b029aaa1db2560b45ca67d688c7f84b8c4c791fe02b3df614f86db1690901c56b45c1530afedfb76038e972722fe7ad728f0e4904e046c230570fe9d41398abe12ef5bc942be33542a4802d98b5d70f2a332ec37fac3514e74ddc0f2cc1a874cd0c78305a21566461309789606bd0bf3f98cda8044629a1")
	c2Bytes, _ := hex.DecodeString("255044462d312e330a25e2e3cfd30a0a0a312030206f626a0a3c3c2f57696474682032203020522f4865696768742033203020522f547970652034203020522f537562747970652035203020522f46696c7465722036203020522f436f6c6f7253706163652037203020522f4c656e6774682038203020522f42697473506572436f6d706f6e656e7420383e3e0a73747265616d0affd8fffe00245348412d3120697320646561642121212121852fec092339759c39b1a1c63c4c97e1fffe017346dc9166b67e118f029ab621b2560ff9ca67cca8c7f85ba84c79030c2b3de218f86db3a90901d5df45c14f26fedfb3dc38e96ac22fe7bd728f0e45bce046d23c570feb141398bb552ef5a0a82be331fea48037b8b5d71f0e332edf93ac3500eb4ddc0decc1a864790c782c76215660dd309791d06bd0af3f98cda4bc4629b1")

	scriptSig := &Script{}
	scriptSig.Commands = append(scriptSig.Commands, Command{Bytes: c1Bytes})
	scriptSig.Commands = append(scriptSig.Commands, Command{Bytes: c2Bytes})

	scriptPubKey := &Script{}
	scriptPubKey.Commands = append(scriptPubKey.Commands, Command{Bytes: scriptPubkeyBytes})

	combinedScript := Combine(*scriptPubKey, *scriptSig)

	if !combinedScript.Evaluate(big.NewInt(0), 0, 0, 0, nil) {
		t.Fatalf("failed to evaluate script")
	}
}

func TestP2pkh(t *testing.T) {
	MakeP2pkh([]byte{})
}

func TestMakeP2wpkhSmoke(t *testing.T) {
	MakeP2wpkh([]byte{})
}

func TestP2wpkh(t *testing.T) {

}
