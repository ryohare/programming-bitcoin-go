package messages

import "testing"

func TestMakeGetHeaders(t *testing.T) {
	MakeGetHeaders(
		70015,
		1,
		make([]byte, 32),
		nil,
	)
}
