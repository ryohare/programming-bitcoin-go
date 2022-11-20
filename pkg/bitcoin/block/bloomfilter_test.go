package block

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	b1 := []byte("Hellow World")
	b2 := []byte("Goodbye!")
	target, _ := hex.DecodeString("4000600a080000010940")

	filter := MakeBloomFilter(10, 5, 99)

	filter.Add(b1)
	filter.Add(b2)

	// filterBytes := utils.BytesToBitField()

	fmt.Printf("%x\n%x\n", target, filter.BitField)
}
