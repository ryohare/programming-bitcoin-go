package txfetcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ryohare/programming-bitcoin-go/pkg/bitcoin/transaction/tx"
	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

var Service *TxFetcher

type TxFetcher struct {
	cache map[string]*tx.Transaction
}

func init() {
	Service = &TxFetcher{}
}

func getUrl(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api/"
	} else {
		return "https://blockstream.info/mainnet/api/"
	}
}

func (t TxFetcher) Fetch(txID string, testnet, fresh bool) (*tx.Transaction, error) {

	if !fresh {
		if val, ok := t.cache[txID]; ok {
			fmt.Println(val)
			return val, nil
		}
	}

	url := fmt.Sprintf("%s/tx/%s.hex", getUrl(testnet), txID)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	// check --- i dont know what right now
	var trans *tx.Transaction
	if raw[4] == 0 {
		raw = append(raw[:4], raw[6:]...)
		trans = tx.ParseTransaction(raw)
		reader := bytes.NewReader(raw[len(raw)-4:])
		trans.Locktime = utils.LittleEndianToInt(reader)
	} else {
		trans = tx.ParseTransaction(raw)
	}

	trans.Testnet = testnet
	t.cache[txID] = trans

	return trans, nil
}