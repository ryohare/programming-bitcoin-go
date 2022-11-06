package messages

type Command string

type Message interface {
	Serialize() []byte
	GetCommand() Command
}
