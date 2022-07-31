package secp256k1

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAddress(t *testing.T) {
	//  5002 (use uncompressed SEC on testnet)
	//  2020^5 (use compressed SEC on testnet)
	//  0x12345deadbeef (use compressed SEC on mainnet)

	b1 := big.NewInt(5002)
	priv, _ := MakePrivateKeyFromBigInt(b1)
	a1 := priv.Point.Address(false, true)

	fmt.Println(a1)
}
