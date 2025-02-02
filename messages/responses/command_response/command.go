package command_response

import (
	"socks5_server/messages"
	"socks5_server/messages/shared"
)

// Response types as defined in RFC1928
const (
	Success                       = 0
	SocksServerFailure            = 1
	ConnectionNotAllowedByRuleSet = 2
	NetworkUnreachable            = 3
	HostUnreachable               = 4
	ConnectionRefused             = 5
	TtlExpired                    = 6
	CommandNotSupported           = 7
	AddressTypeNotSupported       = 8
)

const statusIdPos = 1
const atypPos = 3
const dstAddrStartPos = atypPos + 1

// Represents a response to command request for proxying data as defined in RC1928
type CommandResponse struct {
	Status   uint16
	BND_ADDR shared.DstAddr
	BND_PORT uint16
}

// Converts CommandResponse into wire-transferable byte array
func (cmd *CommandResponse) ToBytes() ([]byte, error) {
	req := make([]byte, 0)
	req = append(req, []byte{0x05, byte(cmd.Status), 0x00, byte(cmd.BND_ADDR.Type)}...)
	dstAddrBytes, err := cmd.BND_ADDR.ToBytes()
	if err != nil {
		return []byte{}, err
	}
	req = append(req, dstAddrBytes...)
	req = append(req, byte(cmd.BND_PORT>>8), byte(cmd.BND_PORT))
	return req, nil
}

// Constructs command from a wire-transferable data
func (cmd *CommandResponse) Deserialize(req []byte) error {
	if len(req) < 10 { //ver+cmd+rsv+atyp+ipv4+port is 10 bytes at least
		return messages.MalformedMessageError{}
	}
	if req[0] != messages.PROTOCOL_VERSION {
		return messages.MismatchedSocksVersionError{}
	}
	if req[2] != 0x00 {
		return messages.InvalidReservedFieldError{Value: uint16(req[2])}
	}
	if err := cmd.deserializeStatus(req); err != nil {
		return err
	}
	if err := cmd.deserializeBindAddrAndPort(req); err != nil {
		return err
	}
	return nil
}

func (cmd *CommandResponse) deserializeStatus(req []byte) error {
	candidate := uint16(req[statusIdPos])
	if candidate > 9 {
		return &InvalidStatusError{Status: candidate}
	}
	cmd.Status = candidate
	return nil
}

func (cmd *CommandResponse) deserializeAtyp(req []byte) (uint16, error) {
	candidate := uint16(req[atypPos])
	if candidate != shared.ATYP_IPV4 && candidate != shared.ATYP_FQDN && candidate != shared.ATYP_IPV6 {
		return 0, &InvalidAtypError{Atyp: candidate}
	}
	return candidate, nil
}

func (cmd *CommandResponse) deserializeBindAddrAndPort(req []byte) error {
	atyp, err := cmd.deserializeAtyp(req)
	if err != nil {
		return err
	}
	destAddr := shared.DstAddr{}
	nextByte, err := destAddr.Deserialize(req[dstAddrStartPos:], atyp)
	if err != nil {
		return err
	}
	if len(req) < dstAddrStartPos+nextByte+2 {
		return messages.MalformedMessageError{}
	}
	cmd.BND_ADDR = destAddr
	cmd.BND_PORT = uint16(req[dstAddrStartPos+nextByte])<<8 | uint16(req[dstAddrStartPos+nextByte+1])
	return nil
}
