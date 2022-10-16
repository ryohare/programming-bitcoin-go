package tx

import (
	"encoding/hex"
	"testing"
)

func TestTxFetcher(t *testing.T) {
	b, _ := hex.DecodeString("99a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d")
	trans, err := TxFetcherSvc.Fetch(string(b), true, true)

	if err != nil {
		t.Errorf("failed getting transaction because %s\n", err)
	}

	// do some checks on this so its like you know, a real test
	if len(trans.Inputs) != 1 && len(trans.Outputs) != 14 {
		t.Error("inputs and outputs parsing has failed")
	}
}
