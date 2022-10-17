package tx

import (
	"bytes"
	"encoding/binary"
	"fmt"

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
	b := make([]byte, 8)
	binary.PutVarint(b, int64(txOut.Amount))

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
	var err error
	txOut.ScriptPubkey, err = script.Parse(reader)

	if err != nil {
		fmt.Printf("Failed to parse the script pubkey because %v\n", err.Error())
	}

	return txOut
}
