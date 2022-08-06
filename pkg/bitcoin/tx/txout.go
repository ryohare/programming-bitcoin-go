package tx

import (
	"bytes"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/script"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type TransactionOutput struct {
	Amount       uint64
	ScriptPubkey *script.Script
}

func (txOut TransactionOutput) String() string {
	return ""
}

// Returns the byte serialization of the transaction output
func (txOut TransactionOutput) Serialize() []byte {
	var b []byte

	b = append(b, utils.UInt64ToLittleEndianBytes(txOut.Amount)...)
	b = append(b, txOut.ScriptPubkey.Serialize()...)

	return b
}

//Takes a byte stream and parses the tx_output at the start.
// Returns a TxOut object.
func ParseTransactionOutput(reader *bytes.Reader) *TransactionOutput {
	txOut := &TransactionOutput{}

	// read in the amount first
	txOut.Amount = utils.LittleEndianToUInt64(reader)

	// read in the script pub key next
	txOut.ScriptPubkey = script.Parse(reader)

	return txOut
}
