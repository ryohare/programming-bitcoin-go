package utils

import (
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
