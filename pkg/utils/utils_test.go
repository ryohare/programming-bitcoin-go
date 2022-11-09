package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestHash256(t *testing.T) {
	e := new(big.Int).SetBytes(Hash256([]byte("my secret")))
	z := new(big.Int).SetBytes(Hash256([]byte("my message")))

	eC, _ := new(big.Int).SetString("62971298242950415662486979275162298594154135681004836692467839909933090737920", 10)
	zC, _ := new(big.Int).SetString("992574323290069558693408995600997375871533518660852402323633869568647941752", 10)

	if e.Cmp(eC) != 0 {
		t.Error("the secret does not match")
	}
	if z.Cmp(zC) != 0 {
		t.Error("the message does not match")
	}

}

func TestEncodeBase58(t *testing.T) {
	// 7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d
	// eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c
	// c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6

	a1, _ := new(big.Int).SetString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 16)
	b1 := EncodeBase58(a1.Bytes())

	if string(b1) != "9MA8fRQrT4u8Zj8ZRd6MAiiyaxb2Y1CMpvVkHQu5hVM6" {
		t.Error("base58 encoding failed")
	}
	a2, _ := new(big.Int).SetString("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", 16)
	b2 := EncodeBase58(a2.Bytes())

	if string(b2) != "4fE3H2E6XMp4SsxtwinF7w9a34ooUrwWe4WsW1458Pd" {
		t.Error("base58 encoding failed")
	}
	a3, _ := new(big.Int).SetString("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", 16)
	b3 := EncodeBase58(a3.Bytes())

	if string(b3) != "EQJsjkd6JaGwxrjEhfeqPenqHwrBmPQZjJGNSCHBkcF7" {
		t.Error("base58 encoding failed")
	}
}

func TestCovertLittleEndian(t *testing.T) {
	//0xdeadbeef
	n, _ := new(big.Int).SetString("deadbeef", 16)
	a := ConvertLittleEndianToBigInt(n.Bytes())
	fmt.Println(a)
}

func TestConvertIntToLittleEndian(t *testing.T) {
	val := big.NewInt(4022250974)
	a := ConvertIntToLittleEndian(val)
	fmt.Println(a)
}

func TestDecodeBase58(t *testing.T) {
	pyAnswer, _ := hex.DecodeString("d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f")
	addr := "mzx5YhAH9kNHtcN481u6WkjeHjYtVeKVh2"
	goAnswer, err := DecodeBase58(addr)
	if err != nil {
		t.Fatalf("failed to decode address")
	}

	for i := range goAnswer {
		if goAnswer[i] != pyAnswer[i] {
			t.Fatalf("failed validation on byte %d", i)
		}
	}
}

func TestP2shCalculation(t *testing.T) {
	h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")

	// main net test
	result := EncodeBase58Checksum(append([]byte{0x05}, h160...))

	expectedResult := "3CLoMMyuoDQTPRD3XYZtCvgvkadrAdvdXh"

	resultStr := string(result)

	if expectedResult != resultStr {
		t.Fatalf("failed to verify the address")
	}

}

func TestRedeemScriptHashing(t *testing.T) {

	binaryRedeemScript, _ := hex.DecodeString("5221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152ae")
	binaryHashResult, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	h160 := Hash160(binaryRedeemScript)

	if len(h160) != len(binaryHashResult) {
		t.Fatalf("lengths are not the same")
	}

	for i := range h160 {
		if h160[i] != binaryHashResult[i] {
			t.Fatalf("failed to verify result at byte position %d", i)
		}
	}

}

func decodeHexString(str string, t *testing.T) []byte {
	hexBytes, err := hex.DecodeString(str)
	if err != nil {
		t.Fatalf("failed to decode string because %s", err.Error())
	}
	return hexBytes
}

func TestMerkleParentLevel(t *testing.T) {
	txHashes := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
		"2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908",
		"b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0",
	}

	txHashesBytes := make([][]byte, 0, len(txHashes))
	for _, txHash := range txHashes {
		txHashesBytes = append(txHashesBytes, decodeHexString(txHash, t))
	}

	// to test, we need to boil it down until there is only 1 hash left
	for {
		if len(txHashesBytes) == 1 {
			break
		}
		var err error
		txHashesBytes, err = MerkleParentLevel(txHashesBytes)
		if err != nil {
			t.Fatalf("failed to make merkle parent level because %s", err.Error())
		}
	}

	target, _ := hex.DecodeString("acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6")

	if !CompareByteArrays(target, txHashesBytes[0]) {
		t.Fatalf("byte arrays do not match %s vs %s", target, txHashesBytes[0])
	}
}

func TestMerkleRoot(t *testing.T) {
	txHashes := []string{
		"42f6f52f17620653dcc909e58bb352e0bd4bd1381e2955d19c00959a22122b2e",
		"94c3af34b9667bf787e1c6a0a009201589755d01d02fe2877cc69b929d2418d4",
		"959428d7c48113cb9149d0566bde3d46e98cf028053c522b8fa8f735241aa953",
		"a9f27b99d5d108dede755710d4a1ffa2c74af70b4ca71726fa57d68454e609a2",
		"62af110031e29de1efcad103b3ad4bec7bdcf6cb9c9f4afdd586981795516577",
		"766900590ece194667e9da2984018057512887110bf54fe0aa800157aec796ba",
		"e8270fb475763bc8d855cfe45ed98060988c1bdcad2ffc8364f783c98999a208",
	}

	txHashesBytes := make([][]byte, 0, len(txHashes))
	for _, txHash := range txHashes {
		txHashesBytes = append(txHashesBytes, decodeHexString(txHash, t))
	}

	// all the hashes are stored littlen endian, so we need to reorder the byte arary
	for _, txHash := range txHashesBytes {
		MutableReorderBytes(txHash)
	}

	// Get the merkle root
	root, err := MerkleRoot(txHashesBytes)
	if err != nil {
		t.Fatalf("failed to get merkle root because %s", err.Error())
	}

	// need to reorder the root into little endian for evaulation
	MutableReorderBytes(root)

	target, _ := hex.DecodeString("654d6181e18e4ac4368383fdc5eead11bf138f9b7ac1e15334e4411b3c4797d9")

	if !CompareByteArrays(target, root) {
		t.Fatalf("byte arrays are not equal %s vs %s", target, root)
	}
}
