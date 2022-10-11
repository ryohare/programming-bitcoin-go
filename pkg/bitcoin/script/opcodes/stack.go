package opcodes

type StackElementType int

const ( // iota is reset to 0
	data    = iota // c0 == 0
	command        // c1 == 1
)

// Stack is purely abstract.
type StackElement interface {
	OpCode() int
	GetBytes() []byte
	Execute() bool
	Type() StackElementType
}

type Stack []StackElement

func (s Stack) Push(v StackElement) Stack {
	return append(s, v)
}

func (s Stack) pop() StackElement {
	l := len(s)
	return s[l-1]
}
