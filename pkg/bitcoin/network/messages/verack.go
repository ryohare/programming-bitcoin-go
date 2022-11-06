package messages

const COMMAND_VERACK Command = "verack"

type VersionAck struct{}

func ParseVerAck(payload []byte) *VersionAck {
	return &VersionAck{}
}

func MakeVerAck() {}

func (v *VersionAck) Serialize() []byte {
	return nil
}

func (v VersionAck) GetCommand() Command {
	return COMMAND_VERACK
}
