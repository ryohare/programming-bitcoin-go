package tx

import (
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
	// reverse the address
	// txStr := fmt.Sprintf("%x ", utils.MutableReorderBytes([]byte(txID)))

	url := fmt.Sprintf("%s/tx/%s/raw", getUrl(testnet), txID)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get HTTP 200 response. Recieved %s", resp.Status)
	}

	raw, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	// This was code to ignore the segwit transactions
	// if raw[4] == 0 {
	// 	raw = append(raw[:4], raw[6:]...)
	// 	trans = ParseTransaction(raw)
	// 	reader := bytes.NewReader(raw[len(raw)-4:])
	// 	trans.Locktime = utils.LittleEndianToInt(reader)
	// } else {
	// 	trans = ParseTransaction(raw)
	// }
	tx, err := ParseTransaction(raw)

	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction because %s", err.Error())
	}

	// make sure the tx received is the hash requested
	var calculated string
	if tx.Segwit {
		calculated = tx.ID()
	} else {
		calculated = fmt.Sprintf("%x", utils.ImmutableReorderBytes(utils.Hash256(raw)))
	}

	if calculated != txID {
		return nil, fmt.Errorf("failed to retrieve the correct transaction. Received %s, requested %s", calculated, txID)
	}

	tx.Testnet = testnet
	t.cache[txID] = tx

	return tx, nil
}
