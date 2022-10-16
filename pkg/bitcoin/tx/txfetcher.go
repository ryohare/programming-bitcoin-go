package tx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ryohare/programming-bitcoin-go/pkg/utils"
)

var TxFetcherSvc *TxFetcher

type TxFetcher struct {
	cache map[string]*Transaction
}

func init() {
	TxFetcherSvc = &TxFetcher{
		cache: make(map[string]*Transaction),
	}
}

func getUrl(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api"
	} else {
		return "https://blockstream.info/mainnet/api"
	}
}

func (t TxFetcher) Fetch(txID string, testnet, fresh bool) (*Transaction, error) {
	if !fresh {
		if val, ok := t.cache[txID]; ok {
			fmt.Println(val)
			return val, nil
		}
	}
	// need to reporder the bytes for string print

	txStr := fmt.Sprintf("%x", utils.ReorderBytes([]byte(txID)))

	url := fmt.Sprintf("%s/tx/%s/raw", getUrl(testnet), txStr)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	// check --- i dont know what right now
	var trans *Transaction
	if raw[4] == 0 {
		raw = append(raw[:4], raw[6:]...)
		trans = ParseTransaction(raw)
		reader := bytes.NewReader(raw[len(raw)-4:])
		trans.Locktime = utils.LittleEndianToInt(reader)
	} else {
		trans = ParseTransaction(raw)
	}

	trans.Testnet = testnet
	t.cache[txID] = trans

	return trans, nil
}
