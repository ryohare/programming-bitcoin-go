package opcodes

type OpHash256 struct {
	RawBytes []byte
}

func (o *OpHash256) OpCode() int {
	// or do we return byte here 0x
	return 170
}

// Execute the op on the supplied stack
func (o *OpHash256) Execute(s Stack) bool {
	// OpDup will duplidate the top option on the stack

	if len(s) < 1 {
		// nothing to duplicate, this is an error
		return false
	}

	toHash := (s)[len(s)]
	s = append(s, toHash)

	return true
}

// Returns the raw command information
func (o *OpHash256) GetBytes() []byte {
	return o.RawBytes
}
