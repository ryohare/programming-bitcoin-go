package tx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/script"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

// SIGHASH byte fields
const (
	SIGHASH_ALL       uint32 = 1
	SIGNHASH_NONE     uint32 = 2
	SIGHASH_SINGLE    uint32 = 3
	SIGHASH_ANYCANPAY uint32 = 4
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
	return utils.MutableReorderBytes(utils.Hash256(serial))
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

// Get the signature hash of the transaction.
func (t Transaction) SigHash(inputIndex int, redeemScript *script.Script, sigHash uint32, testnet bool) (*big.Int, error) {
	// start with getting the version from the transaction
	// it is the first element of the serialization stored
	// in little endian formant. For memory allocation, using
	// the o'reilly book, lets allocate all the memory upfront
	// for performance reasons
	//
	// In essence, a bitcoin transaction is just 300 to 400 bytes
	// of data and has to reach any one of tens of thousands of bitcoin nodes.
	s := make([]byte, 4)
	binary.LittleEndian.PutUint32(s, uint32(t.Version))

	// next in the serialization is the number of input
	// transactions. A var int is either 4 or 8 bytes, so
	// we will allocate enough memory for the worst case
	// txInLenBytes := make([]byte, 4)
	// binary.PutUvarint(txInLenBytes, uint64(len(t.Inputs)))
	// s = append(s, txInLenBytes...)
	varint, err := utils.EncodeUVarInt(uint64(len(t.Inputs)))
	for err != nil {
		return nil, err
	}
	s = append(s, varint...)

	// Iterate over the inputs looking for the inputs requring
	// a signature (passed into the function). Also, seralize
	// each input and place into the serialization buffer
	for i, txIn := range t.Inputs {

		// see if the current index matches the index to sign. This
		// means we include the script pub key so we can push in the
		// script sig "unlock" the funds
		if i == inputIndex {
			// Create a TxIn with the correct script pub key

			signedTxIn := &TransactionInput{
				PrevTx:    txIn.PrevTx,
				PrevIndex: txIn.PrevIndex,
				ScriptSig: nil,
				Sequence:  txIn.Sequence,
			}

			// If a redeem script is provided, this is a P2SH transaction
			// and we should use the redeem script as the script sig.
			// otherwise, it is a P2PKH and we use the scriptPunkey
			// as the redeem script
			if redeemScript != nil {
				signedTxIn.ScriptSig = redeemScript
			} else {
				scriptPubKey, err := txIn.ScriptPubkey(testnet)
				if err != nil {
					return nil, fmt.Errorf("failed to parse script pubkey")
				}
				signedTxIn.ScriptSig = scriptPubKey
			}

			signedTxInBytes := signedTxIn.Serialize()
			fmt.Printf("%x\n", signedTxInBytes)
			s = append(s, signedTxInBytes...)
		} else {
			// this is an input we are not signin, and thus not spending
			// in this transaction, so we include it but we do not include
			// the script pub key
			signedTxIn := &TransactionInput{
				PrevTx:    txIn.PrevTx,
				PrevIndex: txIn.PrevIndex,
				Sequence:  txIn.Sequence,
			}
			signedTxInBytes := signedTxIn.Serialize()
			s = append(s, signedTxInBytes...)
		}
	}

	// encode the length of the txOuts into the buffer
	// Max size of a varint is 8 bytes, make the buffer
	// 8 in length for the worst case
	b, err := utils.EncodeUVarInt(uint64(len(t.Outputs)))
	if err != nil {
		return nil, err
	}
	s = append(s, b...)
	fmt.Printf("%x\n", s)
	// next we serlaized all the transactions outputs of this transaction
	for _, txOut := range t.Outputs {
		txOutBytes := txOut.Serialize()
		s = append(s, txOutBytes...)
	}

	// next field is the locktime for the transaction. This is encoded
	// as a little endian int
	locktimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktimeBytes, uint32(t.Locktime))
	s = append(s, locktimeBytes...)

	// Set the SIGHASH flag.
	sighashBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sighashBytes, sigHash)
	s = append(s, sighashBytes...)

	// To get the transaction hash, we run it through hash256
	h256 := utils.Hash256(s)

	// return a big int as we want
	return new(big.Int).SetBytes(h256), nil
}

// Verify the input can be spent by this wallet
func (t Transaction) VerifyInput(inputIndex int) (bool, error) {

	// get the input transaction referenced by the index
	txIn := t.Inputs[inputIndex]

	// pull off the script pub key
	scriptPubkey, err := txIn.ScriptPubkey(t.Testnet)

	if err != nil {
		return false, err
	}

	// check the pubkey type. If it is a P2SH, we need to create
	// the redeem script to spend the funds
	var redeemScript *script.Script
	if scriptPubkey.IsP2shScriptPubkey() {
		// Get the redeem script off the input, which is the last element
		cmd := txIn.ScriptSig.Commands[len(txIn.ScriptSig.Commands)-1]

		// Scripts always start with the length of the script, so add in the length
		// of the script so we can parse it correctly
		redeemScriptLenBytes := make([]byte, 8)
		binary.BigEndian.PutUint32(redeemScriptLenBytes, uint32(len(cmd.Bytes)))
		redeemScriptBytes := append(redeemScriptLenBytes, cmd.Bytes...)

		// Parse the script into a script
		var err error
		redeemScript, err = script.Parse(bytes.NewReader(redeemScriptBytes))
		if err != nil {
			return false, err
		}
	} else {
		redeemScript = nil
	}

	// get the script sig from the redeem script
	z, err := t.SigHash(inputIndex, redeemScript, SIGHASH_ALL, true)
	if err != nil {
		return false, err
	}

	// Combine the scripts
	combinedScript := script.Combine(*scriptPubkey, *txIn.ScriptSig)

	// valuate the transaction. If it evaluates to true, then the redeem script
	// or the pub key supplied is valid for the transaction and is allowed
	// to spend the funds encumbered with this
	return combinedScript.Evaluate(z, uint64(t.Locktime), uint64(txIn.Sequence), uint64(t.Version)), nil

}

// verify the transaction is valid
func (t Transaction) Verify(testnet bool) bool {
	if t.Fee(testnet) <= 0 {
		// cant have a negative fee
		return false
	}

	for i := range t.Inputs {
		verify, err := t.VerifyInput(i)
		if err != nil || !verify {
			return false
		}
	}
	return true
}

// Checks if this transaction is a coinbase transaction
func (t Transaction) IsCoinbase() bool {
	// first check that the number of inputs is 1
	if len(t.Inputs) != 1 {
		return false
	}

	// check that the previous transaction hash is all 0's
	prevTxHash := new(big.Int).SetBytes(t.Inputs[0].PrevTx)
	if prevTxHash.Cmp(big.NewInt(0)) != 0 {
		return false
	}

	// now check that the prev index is 0xffffffff
	prevTxInputs := big.NewInt(0xffffffff)
	if prevTxInputs.Cmp(big.NewInt(int64(t.Inputs[0].PrevIndex))) != 0 {
		return false
	}

	return true
}
