package available_auth_methods

// Implements the first message in the life-cycle of a socks5 connection - the auth methods which the clients supports
// In addition to that, it provides constant related to valid values for the METHODS field
// and safe way to work with them(adding is the only functionality which made sense).
import (
	"socks5_server/messages"
	"socks5_server/messages/shared"
)

const (
	messageVersionIndex          = 0
	messageAvailMethodsIndex     = 1
	messageAuthMethodsStartIndex = 2
)

// AvailableAuthMethods Represents a message, containing all the available auth methods provided by CLIENT. The message field is private, because the library ensures
// consistency and validity of the methods by a setter.
type AvailableAuthMethods struct {
	methods []uint16
}

// Methods is Getter for a field in the AvailableAuthMethods struct.
func (m *AvailableAuthMethods) Methods() []uint16 {
	return m.methods
}

// AddMethod Checks if given method is valid and if so, it appends it to any other methods already presented in the instance
func (m *AvailableAuthMethods) AddMethod(method uint16) error {
	if method > shared.JsonParameterBlock || method == shared.Unassigned {
		return messages.UnknownAuthMethodError{Method: method}
	}
	m.methods = append(m.methods, method)
	return nil
}

// AddMultipleMethods is wrapper around AddMethod to support adding multiple methods with single call
func (m *AvailableAuthMethods) AddMultipleMethods(methods []uint16) error {
	for _, method := range methods {
		if err := m.AddMethod(method); err != nil {
			return err
		}
	}
	return nil
}

// ToBytes Converts the structure into wire-transferable data
func (m *AvailableAuthMethods) ToBytes() []byte {
	typesBytes := make([]byte, 0)
	for _, method := range m.methods {
		typesBytes = append(typesBytes, byte(method))
	}
	headersBytes := []byte{messages.PROTOCOL_VERSION, byte(len(m.methods))}
	return append(headersBytes, typesBytes...)
}

// Deserialize Constructs AvailableAuthMethods from bytes transferred over the wire
func (m *AvailableAuthMethods) Deserialize(buf []byte) error {
	if len(buf) < 3 {
		return messages.MalformedMessageError{}
	}
	if buf[messageVersionIndex] != messages.PROTOCOL_VERSION {
		return messages.MismatchedSocksVersionError{}
	}

	authMethodsCount := uint16(buf[messageAvailMethodsIndex])

	lastAuthMethodIndex := int(messageAuthMethodsStartIndex + authMethodsCount)
	if lastAuthMethodIndex >= len(buf) {
		return messages.MalformedMessageError{}
	}

	for i := messageAuthMethodsStartIndex; i < lastAuthMethodIndex; i++ {
		currentAuthMethod := uint16(buf[i])
		if err := m.AddMethod(currentAuthMethod); err != nil {
			return err
		}
	}
	return nil
}
