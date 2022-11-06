package messages

const COMMAND_VERACK = "verack"

type VersionAck struct{}

func ParseVerAck() {}

func MakeVerAck() {}

func (v *VersionAck) Serialize() []byte {
	return nil
}

func (v VersionAck) GetCommand() string {
	return COMMAND_VERACK
}
