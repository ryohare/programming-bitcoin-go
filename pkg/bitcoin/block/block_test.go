package block

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/script"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

func TestParsingGensisBlock(t *testing.T) {
	gensisBytes, _ := hex.DecodeString("4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73")

	s, err := script.Parse(bytes.NewReader(gensisBytes))

	if err != nil {
		t.Fatalf("failed to parse the gensis block script sig because %s", err.Error())
	}

	fmt.Println(string(s.Commands[2].Bytes))

	if string(s.Commands[2].Bytes) != "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks" {
		t.Fatalf("failed to parse gensis block text")
	}
}

func TestParsingBlockHeight(t *testing.T) {
	blockBytes, _ := hex.DecodeString("5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00")
	s, err := script.Parse(bytes.NewReader(blockBytes))
	if err != nil {
		t.Fatalf("failed to parse script sig because %s", err.Error())
	}

	// get the block height from the commands
	targetHeighBytes := utils.IntToLittleEndianBytes(465879)
	for i := range s.Commands[0].Bytes {
		if targetHeighBytes[i] != s.Commands[0].Bytes[i] {
			t.Fatalf("failed to validate the block height for the parsed byte stream at offset %d", i)
		}
	}

}

func TestSerializeBlockHeader(t *testing.T) {
	bh := &BlockHeader{
		Version:       1,
		PreviousBlock: make([]byte, 32),
		MerkleRoot:    make([]byte, 32),
		Timestamp:     1234,
		Bits:          make([]byte, 4),
		Nonce:         make([]byte, 4),
	}

	serializedBlockHeader, err := bh.SerializeHeader()
	if err != nil {
		t.Fatalf("failed to serialize the block header because %s", err.Error())
	}

	// undo it back into the struct
	newBh, err := ParseHeader(bytes.NewReader(serializedBlockHeader))
	if err != nil {
		t.Fatalf("failed to parse the perviously serialized block header because %s", err.Error())
	}

	if newBh.Version != bh.Version || newBh.Timestamp != bh.Timestamp {
		t.Fatalf("serialization failure occured")
	}
}

func TestParseBlockHeader(t *testing.T) {
	blockBytes, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	b, err := ParseHeader(bytes.NewReader(blockBytes))
	if err != nil {
		t.Fatalf("failed to parse the block")
	}

	fmt.Printf("%x\n", b.Version)
	fmt.Printf("%x\n", b.Version>>29)

	// check if BIP9 is enabled
	if b.Version>>29 == 1 {
		// indicates  that BIP9 is enabled, which is the signaling protocol
		fmt.Println("BIP9 enabled")
	} else {
		t.Fatalf("failed to verify that BIP9 is enabled")
	}

	// check for bip91 which is another signaling protocol
	if b.Version>>4&1 == 1 {
		t.Fatal("failed to verify BIP91 as enabled")
	} else {
		fmt.Println("BIP91 enabled")
	}

	// check if BIP141 is enabled (segwit)
	if b.Version>>1&1 == 1 {
		fmt.Println("Segwit enabled")
	} else {
		t.Fatal("segwit is not enabled")
	}
}

func TestViewPow(t *testing.T) {
	blockBytes, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")

	// hash the value to see the block hash
	h256 := utils.Hash256(blockBytes)

	// order the bytes
	h256 = utils.MutableReorderBytes(h256)

	check := "0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523"

	// print out the hash
	fmt.Printf("%64x\n", h256)

	if check != fmt.Sprintf("%64x", h256) {
		t.Fatalf("failed to verify the hash for the block")
	}
}

func TestCalculatePowTarget(t *testing.T) {
	bits, _ := hex.DecodeString("e93c0118")
	exponent := bits[len(bits)-1]
	coeffecient := utils.LittleEndianToInt(bytes.NewReader(bits[:len(bits)-1]))

	target := new(big.Int).Exp(big.NewInt(256), big.NewInt(int64(exponent)-3), nil)
	fmt.Println(target)
	target = target.Mul(target, big.NewInt(int64(coeffecient)))

	// reporting
	fmt.Println(exponent)
	fmt.Println(coeffecient)
	fmt.Println(target)
	fmt.Printf("%64x\n", target.Bytes())

	// proof example
	proofBytes, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	proofBytesHash := utils.Hash256(proofBytes)

	// use good o'le big ints...
	proof := new(big.Int).SetBytes(proofBytesHash)

	// validate the the target is higher than the proof
	// print(proof < target)
	fmt.Println(proof.Cmp(target))

	// now construct the difficulty
	// (0xffff * 256^)
	_tmp := new(big.Int).Exp(big.NewInt(256), big.NewInt(0x1d-3), nil)
	difficulty := _tmp.Mul(big.NewInt(0xffff), _tmp)
	difficulty = difficulty.Div(difficulty, target)
	fmt.Println(difficulty)

}

func TestBlockDifficultyTarget(t *testing.T) {
	lastBlockBytes, _ := hex.DecodeString("00000020fdf740b0e49cf75bb3d5168fb3586f7613dcc5cd89675b0100000000000000002e37b144c0baced07eb7e7b64da916cd3121f2427005551aeb0ec6a6402ac7d7f0e4235954d801187f5da9f5")
	firstBlockBytes, _ := hex.DecodeString("000000201ecd89664fd205a37566e694269ed76e425803003628ab010000000000000000bfcade29d080d9aae8fd461254b041805ae442749f2a40100440fc0e3d5868e55019345954d80118a1721b2e")

	lastBlock, err := ParseHeader(bytes.NewReader(lastBlockBytes))
	if err != nil {
		t.Fatalf("failed to parse last block")
	}
	firstBlock, err := ParseHeader(bytes.NewReader(firstBlockBytes))
	if err != nil {
		t.Fatalf("failed to parse last block")
	}

	timeDifferential := lastBlock.Timestamp - firstBlock.Timestamp

	// make sure that it took more that 8 weeks to find the last 2015 blocks
	// otherwise we need to decrease
	if timeDifferential > TwoWeeks*4 {
		timeDifferential = TwoWeeks * 4
	}

	if timeDifferential < TwoWeeks/4 {
		timeDifferential = TwoWeeks / 4
	}

	oldTarget := lastBlock.Target()
	newTarget := oldTarget.Mul(oldTarget, big.NewInt(int64(timeDifferential)))
	fmt.Println(newTarget)

}
