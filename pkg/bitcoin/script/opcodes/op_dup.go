package opcodes

type OpDup struct {
	RawBytes []byte
}

func (o *OpDup) OpCode() int {
	// TODO or do we return byte here 0x076
	return 118
}

// Duplicate the top item on the stack
func (o *OpDup) Execute(s Stack) bool {
	// OpDup will duplidate the top option on the stack

	if len(s) < 1 {
		// nothing to duplicate, this is an error
		return false
	}

	toDup := (s)[len(s)]
	s = append(s, toDup)

	return true
}

// Returns the raw command information
func (o *OpDup) GetBytes() []byte {
	return o.RawBytes
}
