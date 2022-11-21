package tx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
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
	Version int

	// Transaction inputs which map to previousOuts from existing transactions
	Inputs []*TransactionInput

	// New Utxos being created
	Outputs []*TransactionOutput

	// Locktime afte which the transaction can be spend
	Locktime int

	// Internal Testnet flag
	Testnet bool

	// Serlaization holder - Not Used - Use tx.Serialize() to get the serlaization
	Serialization []byte

	// Internal Segwit flag
	Segwit bool

	// Withness programs
	Witness [][]byte

	// Hash of the outputs for "this" transaction
	hashOutputs []byte

	// hash of the outputs for the previous transactions (utxos)
	hashPrevOuts []byte

	// hash of the previous outputs sequences (utxos)
	hashSequence []byte
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

// Serialize a segwit transaction. Declared as private because we want the main entrypoint
// to be the serialize function for all serialization requests of this struct
func (t Transaction) serializeSegwit() []byte {
	// setup the var that will be the serialization
	var tx []byte

	// serialize the version number first
	tx = append(tx, utils.IntToLittleEndianBytes(t.Version)...)

	// add in the segwit marker and flag
	tx = append(tx, []byte{0x00, 0x01}...)

	// now we encode the number of inputs that are in the transaction
	// which gets encoded as a varint
	tx = append(tx, utils.IntToVarintBytes(len(t.Inputs))...)

	// loop over all the inputs and append their individual serialiations
	for _, txin := range t.Inputs {
		tx = append(tx, txin.Serialize()...)
	}

	// next are the outputs. Again, like the inputs, firs element is the
	// length which is encoded as var int
	tx = append(tx, utils.IntToVarintBytes(len(t.Outputs))...)

	// serialize each outut
	for _, txout := range t.Outputs {
		tx = append(tx, txout.Serialize()...)
	}

	// now the segwit magic. We need to add in the witness program for each
	// input, which map to prevOuts or Outpoints we are consuming
	for _, txin := range t.Inputs {

		// first element is the number of witness programs we have.
		// this is encoded as little endian
		tx = append(tx, utils.IntToLittleEndianBytes(len(txin.Witness))...)

		// iterate over the witness programs and add them into the serialization
		for _, witness := range t.Witness {

			// based on the parse function, we check if 0x00 was pushed into the witness array
			// in the python code, they do this by checking if the datatype is an int which
			// doesnt fly in golang because everything is a []byte here.
			// Update: Pulled out the 0x00 check and am just checking if it is of length 1,
			// indicating it is uint8
			if len(witness) == 1 /*&& witness[0] == 0x00*/ {
				tx = append(tx, utils.UInt8ToLittleEndianBytes(uint8(witness[0])))
			} else {
				// else its a varint that needs to encoded as the length of the withness program
				// followed by the program itself
				tx = append(tx, utils.IntToVarintBytes(len(witness))...)
				tx = append(tx, witness...)
			}
		}
	}

	// now just follow the normal serialization the rest of the way
	// serialize locktime as a little endian int
	tx = append(tx, utils.IntToLittleEndianBytes(t.Locktime)...)

	return tx
}

// Serialize a transaction into a byte array
func (t Transaction) Serialize() []byte {
	if t.Segwit {
		return t.serializeSegwit()
	} else {
		return t.serializeLegacy()
	}
}

// Returns the byte serialization of the transaction
func (t Transaction) serializeLegacy() []byte {
	// setup the var that will be the serialization
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

// Return the transaction Id (hash) of the transaction as a byte array
func (t Transaction) Hash() []byte {
	serial := t.serializeLegacy()
	return utils.MutableReorderBytes(utils.Hash256(serial))
}

// Returns the transaction Id as a string
func (t Transaction) ID() string {
	return fmt.Sprintf("%x", t.Hash())
}

// Parse a segwit transaction
func ParseSegwit(serialization []byte) (*Transaction, error) {
	t := &Transaction{}

	// create the reader
	reader := bytes.NewReader(serialization)

	// read the segwit version to know what we are dealing with.
	// version = 0 == normal segwit
	// version = 1 == taproot - sir not appearing in this picture
	t.Version = utils.LittleEndianToInt(reader)

	// read in the segwit marker and flag which are the next 2 bytes
	segwitMarker, _ := reader.ReadByte()
	segwitFlag, _ := reader.ReadByte()
	if segwitMarker != 0x00 && segwitFlag != 0x01 {
		return nil, fmt.Errorf("segwith markers are not correct. Received %x %x", segwitMarker, segwitFlag)
	}

	//
	// Parse the inputs to the transaction
	//
	// number of inputs is the next item in the serialization of the transaction on the wire
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
	// Parse the witness program
	//
	// each input needs a witness in order to spend it in "this" transaction
	for i, _ := range t.Inputs {
		// read in the number of witness items
		// it can be variable in the case of multisig stuff
		// like a set of signatures, or something else
		numWitnesses := utils.ReadVarIntFromBytes(reader)

		// witnesses will be an array of byte arrays
		witnessess := [][]byte{}

		// iterate over the witness programs
		for i := 0; i < int(numWitnesses); i++ {
			witnessLength := utils.ReadVarIntFromBytes(reader)

			// if the witness program is nothing, push in a null byte
			if witnessLength == 0 {
				witnessess = append(witnessess, []byte{0x00})
			}

			// otherwise, read in N bytes as the witness program
			witnessProgram, _ := ioutil.ReadAll(io.LimitReader(reader, int64(witnessLength)))

			// add in the just parsed withness program
			witnessess = append(witnessess, witnessProgram)
		}

		// witness data is coupled with the input
		t.Inputs[i].Witness = witnessess
	}

	//
	// Parse the locktime
	//
	t.Locktime = utils.LittleEndianToInt(reader)

	// set the internal segwit flag
	t.Segwit = true

	return t, nil
}

// Parse a transaction from a byte stream
func ParseTransaction(serialization []byte) (*Transaction, error) {
	t := &Transaction{}

	// make a reader to easily read in the serialization
	reader := bytes.NewReader(serialization)

	// segwith bolt-on, read 4 bytes and ignore and look for
	// the segwit marker in the 5th byte of the transaction
	for i := 0; i < 4; i++ {
		_, err := reader.ReadByte()

		if err != nil {
			return nil, err
		}
	}

	// now read for the segwith marker, if we find it, go the special
	// segwit parser, otherwise, reset the stream and continue with
	// the originally defined processing path
	segwit, _ := reader.ReadByte()
	if segwit == 0x00 {
		return ParseSegwit(serialization)
	} else {
		reader.Reset(serialization)
	}

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

	return t, nil
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

func (t Transaction) HashPrevOuts() []byte {
	// if we have not set the prevOuts for this transaction, then we need to set them
	if len(t.hashPrevOuts) == 0 {
		allPrevOuts := []byte{}
		allSequence := []byte{}

		// construct a byte array of format
		// prevHash1 + length1 ... prevHashN + lengthN
		for _, txin := range t.Inputs {
			allPrevOuts = append(
				allPrevOuts,
				utils.ImmutableReorderBytes(txin.PrevTx)...,
			)
			allPrevOuts = append(
				allPrevOuts,
				utils.IntToLittleEndianBytes(txin.PrevIndex)...,
			)
			allSequence = append(
				allSequence,
				utils.IntToLittleEndianBytes(txin.Sequence)...,
			)
		}

		// we now have two byte arrays, one with the serialization of all the txins hash + index
		// and one for all the sequences

		// Set self to the values
		t.hashPrevOuts = utils.Hash256(allPrevOuts)
		t.hashSequence = utils.Hash256(allSequence)
	}

	return t.hashPrevOuts
}

func (t Transaction) HashOutputs() []byte {
	if len(t.hashPrevOuts) == 0 {
		var allOutputs []byte
		for _, txout := range t.Outputs {
			allOutputs = append(allOutputs, txout.Serialize()...)
		}
		t.hashOutputs = allOutputs
	}

	return t.hashOutputs
}

// Returns the big.Int representation of the hash that needs to be signed for the index inputIndex
// Signing the inputs is unlocking the prevOut from the previous transaction
func (t Transaction) SigHashSegwit(inputIndex int, redeemScript, witnessScript *script.Script) (*big.Int, error) {

	// TODO check bounds on the inputIndex - Probably should be defensivily programming everywhere
	txin := t.Inputs[inputIndex]

	// This is all done per BIP143 Spec which I just dupped with the python code is doing
	var s []byte
	s = utils.IntToLittleEndianBytes(t.Version)
	s = append(s, t.hashPrevOuts...)
	s = append(s, t.hashSequence...)
	s = append(s, utils.ImmutableReorderBytes(txin.PrevTx)...)
	s = append(s, utils.IntToLittleEndianBytes(txin.PrevIndex)...)

	// handle the supplied scripts script
	// TODO - Grok this block
	var scriptCode []byte
	if witnessScript != nil {

		// witness sript was supplied, this is a p2wpkh
		scriptCode = append(scriptCode, witnessScript.Serialize()...)
	} else if redeemScript != nil {
		// if there is a redeem script and no witness script, then this is a p2sh-p2wpkh
		scriptCode = append(
			scriptCode,
			// make a p2pkh serialization
			script.MakeP2pkh(redeemScript.Commands[1].Bytes).Serialize()...,
		)
	} else {

		pubkey, err := txin.ScriptPubkey(t.Testnet)
		if err != nil {
			return nil, fmt.Errorf("failed to get ScriptPubKey because for input idx %d because %s", inputIndex, err.Error())
		}

		scriptCode = append(scriptCode,
			script.Makep2sh(pubkey.Commands[1].Bytes).Serialize()...,
		)
	}

	// get the amounts for the utxo's for the previous transaction
	s = append(s, scriptCode...)
	val, err := txin.Value(t.Testnet)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value for the outputs because %s", err.Error())
	}
	s = append(s, utils.UInt64ToLittleEndianBytes(uint64(val))...)
	s = append(s, utils.IntToLittleEndianBytes(txin.Sequence)...)
	s = append(s, t.HashOutputs()...)
	s = append(s, utils.IntToLittleEndianBytes(t.Locktime)...)
	s = append(s, utils.IntToLittleEndianBytes(int(SIGHASH_ALL))...)

	// now that we have s, which is what is to be signed for the transaction during transaction signing
	// we calculate the hash. The hash of this "serialization" is what is signed during transaction signing
	h256 := utils.Hash256(s)

	// make the big int
	return new(big.Int).SetBytes(h256), nil

}

// Verify the input can be spent by this wallet
func (t Transaction) VerifyInput(inputIndex int) (bool, error) {

	// get the input transaction referenced by the index
	txIn := t.Inputs[inputIndex]

	// pull off the prevcious output ScriptPubKey. This can be any of the
	// transaction types
	scriptPubkey, err := txIn.ScriptPubkey(t.Testnet)

	if err != nil {
		return false, fmt.Errorf("failed to get ScriptPubKey because %s", err.Error())
	}

	// check the pubkey type. If it is a P2SH, we need to create
	// the redeem script to spend the funds.
	// If it is p2wphk, we need to handle it differently as there
	// will need to parse the witness program
	var redeemScript *script.Script
	z := new(big.Int)
	var witness [][]byte
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

		// handle p2sh-pwpkh type. This is where the witness signature
		// was embedded in the redeem script.
		// The script embedded could be either w p2wpkh or p2wsh type script
		if redeemScript.IsP2wpkhScriptPubkey() {
			// calculate Z the special way
			z, err = t.SigHashSegwit(inputIndex, redeemScript, nil)

			if err != nil {
				return false, fmt.Errorf("failed to calculate z because %s", err.Error())
			}

			// since we are segwit aware, assign the witness data field
			witness = txIn.Witness
		} else {
			witness = nil
		}

	} else if scriptPubkey.IsP2wpkhScriptPubkey() {
		// handle segwit type Outpoint (prevHash.ScriptPubKey)
		// most of the processing path is the same as p2sh, we just
		// need to handle the witness data correctly in this path
		z, err = t.SigHashSegwit(inputIndex, nil, nil)
		if err != nil {
			return false, fmt.Errorf("failed to calculate z because %s", err.Error())
		}
		witness = t.Witness
	} else {
		z, err = t.SigHash(inputIndex, nil, SIGHASH_ALL, true)
		if err != nil {
			return false, fmt.Errorf("failed to calculate z because %s", err.Error())
		}
		witness = nil
		redeemScript = nil
	}

	// Combine the scripts
	combinedScript := script.Combine(*scriptPubkey, *txIn.ScriptSig)

	// valuate the transaction. If it evaluates to true, then the redeem script
	// or the pub key supplied is valid for the transaction and is allowed
	// to spend the funds encumbered with this
	return combinedScript.Evaluate(z, uint64(t.Locktime), uint64(txIn.Sequence), uint64(t.Version), witness), nil
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
