package messages

const COMMAND_PING = "ping"

type Ping struct{}

func ParsePing() {}

func MakePing() {}

func (v *Ping) Serialize() []byte {
	return nil
}

func (v Ping) GetCommand() string {
	return COMMAND_VERACK
}
