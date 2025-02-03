package accept_auth_method

// This package provides a message, with which the server responds when an authentication method is chosen. The authentication methods are provided by the client
import (
	"socks5_server/messages"
	"socks5_server/messages/shared"
)

type AcceptAuthMethod struct {
	method uint16
}

func (aam *AcceptAuthMethod) Method() uint16 {
	return aam.method
}

func (aam *AcceptAuthMethod) SetMethod(method uint16) error {
	if method != shared.NoAcceptableMethods && method > shared.JsonParameterBlock {
		return messages.UnknownAuthMethodError{Method: method}
	}
	if method == shared.Unassigned {
		return messages.UnknownAuthMethodError{Method: shared.Unassigned}
	}
	aam.method = method
	return nil
}

func (aam *AcceptAuthMethod) ToBytes() []byte {
	return []byte{0x05, byte(aam.method)}
}

func (aam *AcceptAuthMethod) Deserialize(buf []byte) error {
	if len(buf) < 2 {
		return messages.MalformedMessageError{}
	}
	if buf[0] != 0x05 {
		return messages.MismatchedSocksVersionError{}
	}
	if err := aam.SetMethod(uint16(buf[1])); err != nil {
		return err
	}
	return nil
}
