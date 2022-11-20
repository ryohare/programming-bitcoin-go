package messages

type GetDataMessage struct {
	Data []byte
}

func MakeGetDataMessage() *GetDataMessage {
	return &GetDataMessage{}
}
