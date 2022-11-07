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
func TestH160ToP2shAddress(t *testing.T) {

}

// func TestManualBits(t *testing.T) {
// 	lastBlockBytes, _ := hex.DecodeString("00000020fdf740b0e49cf75bb3d5168fb3586f7613dcc5cd89675b0100000000000000002e37b144c0baced07eb7e7b64da916cd3121f2427005551aeb0ec6a6402ac7d7f0e4235954d801187f5da9f5")
// 	firstBlockBytes, _ := hex.DecodeString("000000201ecd89664fd205a37566e694269ed76e425803003628ab010000000000000000bfcade29d080d9aae8fd461254b041805ae442749f2a40100440fc0e3d5868e55019345954d80118a1721b2e")
// 	targetBytes, _ := hex.DecodeString("0000000000000000007615000000000000000000000000000000000000000000")

// 	t.Logf("goal value is %x", targetBytes)

// 	// parse the byte streams into blocks
// 	lastBlock, err := block.ParseHeader(bytes.NewReader(lastBlockBytes))
// 	if err != nil {
// 		t.Fatalf("failed to parse last block because %s", err.Error())
// 	}
// 	firstBlock, err := block.ParseHeader(bytes.NewReader(firstBlockBytes))
// 	if err != nil {
// 		t.Fatalf("failed to parse firdst block because %s", err.Error())
// 	}

// 	// calculate the time differential
// 	timeDifferential := lastBlock.Timestamp - firstBlock.Timestamp

// 	// adjust the dificultuy
// 	if timeDifferential > TwoWeeks*4 {
// 		timeDifferential = TwoWeeks * 4
// 	}
// 	if timeDifferential < TwoWeeks/4 {
// 		timeDifferential = TwoWeeks / 4
// 	}

// 	// newTarget := lastBlock.Target() * timeDifferential

// }
