package txin

import (
	"bytes"
	"fmt"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/script"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type Fetcher interface {
	Fetch()
}

type TransactionInput struct {
	// 32 byte hash256 of previous previous transaction's contents
	// little endian
	PrevTx []byte

	// 4 bytes little endian
	PrevIndex int
	ScriptSig *script.Script

	// 4 bytes little endian
	Sequence int
}

func (txIn TransactionInput) String() string {
	return fmt.Sprintf("%s:%d", txIn.PrevTx, txIn.PrevIndex)
}

func (txIn TransactionInput) Hex() string {
	return ""
}

// func (txIn TransactionInput) FetchTx(testnet bool) {
// 	txfetcher.Service.Fetch(txIn.PrevTx, testnet)
// }

func Parse(reader *bytes.Reader) *TransactionInput {
	txIn := &TransactionInput{}

	// read in prev_tx first
	txIn.PrevTx = utils.LittleEndianToBigEndian(reader, 32)

	// prev_index is next
	txIn.PrevIndex = utils.LittleEndianToInt(reader)

	// script_sig is next
	txIn.ScriptSig = script.Parse(reader)

	// sequence is next
	txIn.Sequence = utils.LittleEndianToInt(reader)

	return txIn
}

// Returns the byte serialization of the transaction input
func (txIn TransactionInput) Serialize() []byte {
	// get the reversed byte order of the output hash for the input
	var b []byte

	b = append(b, utils.ReorderBytes(txIn.PrevTx)...)

	// previous index converted from big endian to little endian
	b = append(b, utils.IntToLittleEndianBytes(txIn.PrevIndex)...)

	// script sig
	b = append(b, txIn.ScriptSig.Serialize()...)

	// sequence little endian
	b = append(b, utils.IntToLittleEndianBytes(txIn.Sequence)...)

	return b
}

func MakeTransactionInput(prevTx []byte, prevIndex int, scriptSig *script.Script, sequence uint64) *TransactionInput {
	if sequence == 0 {
		sequence = 0xffffffff
	}

	txIn := &TransactionInput{}

	txIn.PrevTx = prevTx
	txIn.PrevIndex = prevIndex

	// no script sig was passed, use basic script
	if scriptSig == nil {
		txIn.ScriptSig = script.Make()
	} else {
		txIn.ScriptSig = scriptSig
	}

	return txIn
}
