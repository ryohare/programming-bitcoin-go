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
