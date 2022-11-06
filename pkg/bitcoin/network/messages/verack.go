package messages

const COMMAND_VERACK Command = "verack"

type VersionAck struct{}

func ParseVerAck() {}

func MakeVerAck() {}

func (v *VersionAck) Serialize() []byte {
	return nil
}

func (v VersionAck) GetCommand() Command {
	return COMMAND_VERACK
}
