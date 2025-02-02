package messages

const PROTOCOL_VERSION byte = 0x05

type MessageType int

type Socks5Message interface {
	ToByte() []byte
	Deserialize([]byte) error
}
