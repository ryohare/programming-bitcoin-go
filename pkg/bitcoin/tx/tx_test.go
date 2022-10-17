package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/script"
	"github.com/ryohare/programming-bitcoin-go/pkg/ecc/curves/secp256k1"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

const testTx = `010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600`

func TestParseTransaction(t *testing.T) {
	tx, err := hex.DecodeString(testTx)

	if err != nil {
		t.Errorf("failed to parse testTx because %s", err.Error())
	}

	ParseTransaction(tx)
}

func TestSerializeTransaction(t *testing.T) {
	tx, err := hex.DecodeString(testTx)

	if err != nil {
		t.Errorf("failed to parse testTx because %s", err.Error())
	}

	serial := ParseTransaction(tx)

	// re-serialize
	serial.Serialize()
}

func TestSigHash(t *testing.T) {

}

func TestTxE2e(t *testing.T) {
	// this address is no longer live on the blockchain, so we cannot submit the transcation
	// but this test will validate that the transaction is constructed correctly
	prevTx, _ := hex.DecodeString("0d6fe5213c0b3291f208cba8bfb59b7476dffacc4e5cb66f6eb20a080843a299")
	prevIndex := 13
	txIn := MakeTransactionInput(prevTx, prevIndex, nil, 0xffffffff)

	// set the change ammount
	changeAmount := int(0.33 * 100000000)

	// change address in bytes
	changeAddress, err := utils.DecodeBase58("mzx5YhAH9kNHtcN481u6WkjeHjYtVeKVh2")
	if err != nil {
		t.Fatalf("failed to decode address")
	}

	// create the script for sending to the change address
	changeScript := script.MakeP2pkh(changeAddress)

	// create the output transaction now
	changeOutputTx := &TransactionOutput{
		Amount:       uint64(changeAmount),
		ScriptPubkey: changeScript,
	}

	// spend to address
	spendToAddress, err := utils.DecodeBase58("mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf")
	if err != nil {
		t.Fatalf("failed to decode address")
	}

	// create the script for the spend to address
	spendToScript := script.MakeP2pkh(spendToAddress)

	// create the output object
	spendToTx := &TransactionOutput{
		Amount:       uint64(0.1 * 100000000),
		ScriptPubkey: spendToScript,
	}

	// created the inputs and the outputs, can now make the transaction object
	tx := Transaction{
		Version:       1,
		Inputs:        []*TransactionInput{txIn},
		Outputs:       []*TransactionOutput{changeOutputTx, spendToTx},
		Locktime:      0,
		Testnet:       false,
		Serialization: []byte{},
	}

	pyAnswer := "010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d00000000ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000"
	goAnswer := fmt.Sprintf("%x", tx.Serialize())

	fmt.Println(pyAnswer)
	fmt.Println(goAnswer)

	// dump tx for fun
	fmt.Printf("%v\n", tx.Serialize())

	// next we need to sign the input
	z, err := tx.SigHash(0, nil, SIGHASH_ALL, true)
	if err != nil {
		t.Fatalf("failed to sign input because %s\n", err.Error())
	}

	// private key we need to sign...
	privateKey, err := secp256k1.MakePrivateKeyFromBigInt(big.NewInt(8675309))
	if err != nil {
		t.Fatalf("failed to make private key because %s\n", err.Error())
	}

	// get the der formated signature for sighash
	_sig, err := privateKey.Sign(z)
	if err != nil {
		t.Fatalf("failed to sign the transaction because %s\n", err)
	}
	der := _sig.Der()

	// append the SIGHASH flag to the transaction
	sighash := byte(SIGHASH_ALL)
	sig := append(der, sighash)

	// get the public key in SEC compressed format
	secPubKey := privateKey.Point.Sec(true)

	// create the script sig for the input
	// only requires the secPubKey
	scriptSigCmds := []script.Command{
		script.Command{
			Bytes: sig,
		},
		script.Command{
			Bytes: secPubKey,
		},
	}

	// set the script sig for the input added to the tx
	tx.Inputs[0].ScriptSig = &script.Script{Commands: scriptSigCmds}

	// now get the byte stream of the transaction
	binaryTx := tx.Serialize()

	pyAnswer = "010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d0000006c4930460221008ed46aa2cf12d6d81065bfabe903670165b538f65ee9a3385e6327d80c66d3b50221003124f804410527497329ec4715e18558082d489b218677bd029e7fa306a72236012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000"
	goAnswer = fmt.Sprintf("%x", binaryTx)
	fmt.Println(pyAnswer)
	fmt.Println(goAnswer)
}
