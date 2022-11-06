package messages

type Message interface {
	Serialize() []byte
	GetCommand() string
}
