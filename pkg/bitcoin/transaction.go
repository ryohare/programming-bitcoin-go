package bitcoin

import (
	"bytes"
	
	"math/big"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type TransactionInput struct {
}

type TransactionOutput struct {
}

type Transaction struct {
	Version       int
	Inputs        []TransactionInput
	Outputs       []TransactionOutput
	Locktime      int
	Testnet       bool
	Serialization []byte
}

func (t Transaction) Serialize() []byte {
	return []byte{0x00}
}

func (t Transaction) String() string {
	return ""
}

func (t Transaction) Hash() []byte {
	serial := t.Serialize()
	return utils.ReorderBytes(utils.Hash256(serial))
}

func (t Transaction) ID() string {
	return string(t.Hash())
}
     
func ParseTransaction(serialization []byte) *Transaction{
	t := &Transaction{}

	reader := bytes.NewReader(serialization)

	version := make([]byte, 4)
	reader.Read(version)

	// version is stored little endian
	beVersion := new(big.Int).SetBytes(utils.ReorderBytes(version))
	t.Version = int(beVersion.Int64())

	return t
}
