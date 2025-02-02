package messages

import (
	"fmt"
)

type MismatchedSocksVersionError struct{}

func (e MismatchedSocksVersionError) Error() string {
	return "Mismatched socks version"
}

type UnknownAuthMethodError struct {
	Method uint16
}
type MalformedMessageError struct {
}

func (e MalformedMessageError) Error() string {
	return "Message is malformed"
}
func (e UnknownAuthMethodError) Error() string {
	return fmt.Sprintf("Unknown auth method: %d", e.Method)
}

type InvalidReservedFieldError struct {
	Value uint16
}

func (e InvalidReservedFieldError) Error() string {
	return fmt.Sprintf("Invalid reserved field: %d", e.Value)
}

type InvalidAtypError struct {
	Atyp uint16
}

func (e *InvalidAtypError) Error() string {
	return fmt.Sprintf("Invalid Atyp: %d", e.Atyp)
}
