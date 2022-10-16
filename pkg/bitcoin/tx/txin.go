package tx

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

func (txIn TransactionInput) FetchTx(testnet bool) (*Transaction, error) {
	// reorder the bytes to little endian
	b := utils.ReorderBytes(txIn.PrevTx)
	return TxFetcherSvc.Fetch(string(b), testnet, false)
}

// Get the output value by looking up the tx hash. Returns the amount in satoshi.
func (txIn TransactionInput) Value(testnet bool) (int, error) {
	tx, err := txIn.FetchTx(testnet)

	if err != nil {
		return -1, err
	}

	return int(tx.Outputs[txIn.PrevIndex].Amount), nil

}

func ParseTransactionInput(reader *bytes.Reader) *TransactionInput {
	txIn := &TransactionInput{}

	// read in prev_tx first
	txIn.PrevTx = utils.LittleEndianToBigEndian(reader, 32)

	// prev_index is next
	txIn.PrevIndex = utils.LittleEndianToInt(reader)

	// script_sig is next
	var err error
	txIn.ScriptSig, err = script.Parse(reader)

	if err != nil {
		fmt.Printf("Failed parse script sig because %v\n", err.Error())
	}

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

// Default values that need to be passed in are:
// scriptSig = nil
// sequence = 0xffffffff
func MakeTransactionInput(prevTx []byte, prevIndex int, scriptSig *script.Script, sequence uint64) *TransactionInput {
	if sequence == 0 {
		sequence = 0xffffffff
	}

	txIn := &TransactionInput{}

	txIn.PrevTx = prevTx
	txIn.PrevIndex = prevIndex

	// no script sig was passed, use basic script
	if scriptSig == nil {
		txIn.ScriptSig = script.MakeScript()
	} else {
		txIn.ScriptSig = scriptSig
	}

	return txIn
}

// Get the ScriptPubKey by looking up the tx hash. Returns a Script object.
func (txIn TransactionInput) ScriptPubkey(testnet bool) (*script.Script, error) {
	tx, err := txIn.FetchTx(testnet)

	if err != nil {
		return nil, err
	}

	return tx.Outputs[txIn.PrevIndex].ScriptPubkey, nil
}
