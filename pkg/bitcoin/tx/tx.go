package tx

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

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

	// varint for the length of the outputs
	tx = append(tx, utils.IntToVarintBytes(len(t.Outputs))...)

	// serialize each of the outputs now
	for _, v := range t.Outputs {
		tx = append(tx, v.Serialize()...)
	}

	// add the locktime which needs to be serialized as a little endian int
	tx = append(tx, utils.IntToLittleEndianBytes(t.Locktime)...)

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
		ip := ParseTransactionInput(reader)
		t.Inputs = append(t.Inputs, ip)
	}

	//
	// Parse the outputs
	//
	// first is the varint for the length fof the inputs
	numOfOutputs, _ := binary.ReadUvarint(reader)

	// iterate over the outputs and append them to the outputs list
	for i := 0; i < int(numOfOutputs); i++ {
		op := ParseTransactionOutput(reader)
		t.Outputs = append(t.Outputs, op)
	}

	//
	// Parse the locktime
	//
	t.Locktime = utils.LittleEndianToInt(reader)

	return t
}

// Calculates the fee which should be used for a transaction
func (t Transaction) Fee(testnet bool) uint64 {
	var inputSum uint64
	var outputSum uint64

	for _, txIn := range t.Inputs {
		val, _ := txIn.Value(testnet)
		inputSum += uint64(val)
	}
	for _, txOut := range t.Outputs {
		outputSum += txOut.Amount
	}

	return inputSum - outputSum
}
