package messages

const COMMAND_PING Command = "ping"

type Ping struct{}

func ParsePing() {}

func MakePing() {}

func (v *Ping) Serialize() []byte {
	return nil
}

func (v Ping) GetCommand() Command {
	return COMMAND_PING
}
