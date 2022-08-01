package bitcoin

import (
	"bytes"
	
	"math/big"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

type TransactionInput struct {
}

func (txIn TransactionInput) String() string {
	return ""	
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

func (t Transaction) String() string {
	var txInStr string
	for _,v := range(t.Inputs) {
		txInStr = fmt.Sprintf("%s\n%s\n",txInStr, v.String())
	}

	var txOutStr string
	for _,v := range(t.Outputs){
		txOutStr = fmt.Sprintf("%s\n%s\n", txOutStr, v.String())
	}

	retStr = fmt.Sprintf("Tx: %s\nVersion: %d\n txIns:\n%stxOuts:\n%slocktime: %d", t.ID(), t.Version, txInStr, txOutStr, t.Locktime )

	return retStr
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
