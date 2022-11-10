package block

import "testing"

func TestMake(t *testing.T) {
	if Make(16) == nil {
		t.Fatal("failed to make merkle tree")
	}
}
