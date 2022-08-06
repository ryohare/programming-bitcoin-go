package tx

import (
	"testing"
)

func TestTxFetcher(t *testing.T) {
	targetTx := "0d6fe5213c0b3291f208cba8bfb59b7476dffacc4e5cb66f6eb20a080843a299"

	trans, err := TxFetcherSvc.Fetch(targetTx, true, true)

	if err != nil {
		t.Errorf("failed getting transaction because %s\n", err)
	}

	// do some checks on this so its like you know, a real test
	if len(trans.Inputs) != 1 && len(trans.Outputs) != 14 {
		t.Error("inputs and outputs parsing has failed")
	}
}
