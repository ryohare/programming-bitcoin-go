package bitcoin

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type TransactionInput struct {
	// 32 byte hash256 of previous previous transaction's contents
	// little endian
	PrevTx []byte

	// 4 bytes little endian
	PrevIndex int

	ScriptSig *Script

	// 4 bytes little endian
	Sequence int
}

func (txIn TransactionInput) String() string {
	return fmt.Sprintf("%s:%d", txIn.PrevTx, txIn.PrevIndex)
}

func (txIn TransactionInput) Hex() string {
	return ""
}

func ParseInput(reader *bytes.Reader) *TransactionInput {
	txIn := &TransactionInput{}

	// read in prev_tx first
	txIn.PrevTx = utils.LittleEndianToBigEndian(reader, 32)

	// prev_index is next
	txIn.PrevIndex = utils.LittleEndianToInt(reader)

	// script_sig is next
	txIn.ScriptSig = ParseScript(reader)

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

func MakeTransactionInput(prevTx []byte, prevIndex int, scriptSig *Script, sequence uint64) *TransactionInput {
	if sequence == 0 {
		sequence = 0xffffffff
	}

	txIn := &TransactionInput{}

	txIn.PrevTx = prevTx
	txIn.PrevIndex = prevIndex

	// no script sig was passed, use basic script
	if scriptSig == nil {
		txIn.ScriptSig = MakeScript()
	} else {
		txIn.ScriptSig = scriptSig
	}

	return txIn
}

type TransactionOutput struct {
	Amount       uint64
	ScriptPubkey *Script
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
func ParseOutput(reader *bytes.Reader) *TransactionOutput {
	txOut := &TransactionOutput{}

	// read in the amount first
	txOut.Amount = utils.LittleEndianToUInt64(reader)

	// read in the script pub key next
	txOut.ScriptPubkey = ParseScript(reader)

	return txOut
}

type Transaction struct {
	// 4 bytes little endian
	Version       int
	Inputs        []*TransactionInput
	Outputs       []*TransactionOutput
	Locktime      int
	Testnet       bool
	Serialization []byte
}

func (t Transaction) String() string {
	var txInStr string
	for _, v := range t.Inputs {
		txInStr = fmt.Sprintf("%s\n%s\n", txInStr, v.String())
	}

	var txOutStr string
	for _, v := range t.Outputs {
		txOutStr = fmt.Sprintf("%s\n%s\n", txOutStr, v.String())
	}

	retStr := fmt.Sprintf("Tx: %s\nVersion: %d\n txIns:\n%stxOuts:\n%slocktime: %d", t.ID(), t.Version, txInStr, txOutStr, t.Locktime)

	return retStr
}

// Returns the byte serialization of the transaction
func (t Transaction) Serialize() []byte {
	var tx []byte

	// serialize the version number first
	tx = append(tx, utils.IntToLittleEndianBytes(t.Version)...)

	// varint for the length of the inputs
	tx = append(tx, utils.IntToVarintBytes(len(t.Inputs))...)

	// serialize each of the inputs now
	for _, v := range t.Inputs {
		tx = append(tx, v.Serialize()...)
	}

	return tx
}

func (t Transaction) Hash() []byte {
	serial := t.Serialize()
	return utils.ReorderBytes(utils.Hash256(serial))
}

func (t Transaction) ID() string {
	return string(t.Hash())
}

func ParseTransaction(serialization []byte) *Transaction {
	t := &Transaction{}

	// make a reader to easily read in the serialization
	reader := bytes.NewReader(serialization)

	//
	// parse the version
	//
	t.Version = utils.LittleEndianToInt(reader)

	//
	// Parse the inputs
	//
	// first is the varint for the length of the inputs
	numOfInputs, _ := binary.ReadUvarint(reader)

	// iterate over the inputs and append them to the inputs list
	for i := 0; i < int(numOfInputs); i++ {
		ip := ParseInput(reader)
		t.Inputs = append(t.Inputs, ip)
	}

	//
	// Parse the outputs
	//
	// first is the varint for the length fof the inputs
	numOfOutputs, _ := binary.ReadUvarint(reader)

	// iterate over the outputs and append them to the outputs list
	for i := 0; i < int(numOfOutputs); i++ {
		op := ParseOutput(reader)
		t.Outputs = append(t.Outputs, op)
	}

	//
	// Parse the locktime
	//
	t.Locktime = utils.LittleEndianToInt(reader)

	return t
}
