package responses

import "socks/pkg/socks5/shared"

func SelectionMessage(authMethod int) ([]byte, error) {
	return []byte{shared.PROTOCOL_VERSION, byte(authMethod)}, nil
}
