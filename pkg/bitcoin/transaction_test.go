package bitcoin

import (
	"math/big"
	"testing"
)

func TestParseTransaction(t *testing.T) {

	// version 1 in little endian
	b1, _ := new(big.Int).SetString("01000000", 16)

	ParseTransaction(b1.Bytes())
}
