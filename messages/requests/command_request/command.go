package command_request

// Implements a message send by the client to the proxy server requesting a service/command.
import (
	"socks5_server/messages"
	"socks5_server/messages/shared"
)

const commandIdPos = 1
const atypPos = 3
const dstAddrStartPos = atypPos + 1

// Provides name->int mapping for the SOCKS5 commands, as defined in RFC1928
const (
	CONNECT       = 1
	BIND          = 2
	UDP_ASSOCIATE = 3
)

// CommandRequest Represents a request for proxying data, it contains the command type and destination address
type CommandRequest struct {
	CMD      uint16
	DST_ADDR shared.DstAddr
	DST_PORT uint16
}

// ToBytes Converts CommandRequest into wire-transferable byte array
func (cmd *CommandRequest) ToBytes() ([]byte, error) {
	dstAddrBytes, err := cmd.DST_ADDR.ToBytes()
	if err != nil {
		return []byte{}, err
	}

	req := make([]byte, 0)
	req = append(req, []byte{0x05, byte(cmd.CMD), 0x00, byte(cmd.DST_ADDR.Type)}...)
	req = append(req, dstAddrBytes...)
	req = append(req, byte(cmd.DST_PORT>>8), byte(cmd.DST_PORT))
	return req, nil
}

// Deserialize Constructs command from a wire-transferable data
func (cmd *CommandRequest) Deserialize(req []byte) error {
	if len(req) < 10 { //ver+cmd+rsv+atyp+ipv4+port is 10 bytes at least
		return messages.MalformedMessageError{}
	}
	if req[0] != messages.PROTOCOL_VERSION {
		return messages.MismatchedSocksVersionError{}
	}
	if req[2] != 0x00 {
		return messages.InvalidReservedFieldError{Value: uint16(req[2])}
	}
	if err := cmd.deserializeCmd(req); err != nil {
		return err
	}
	if err := cmd.deserializeDstAddrAndPort(req); err != nil {
		return err
	}
	return nil
}

func (cmd *CommandRequest) deserializeCmd(req []byte) error {
	candidate := uint16(req[commandIdPos])
	if candidate < 1 || candidate > 3 {
		return &InvalidCommandError{CommandType: candidate}
	}
	cmd.CMD = candidate
	return nil
}

func (cmd *CommandRequest) deserializeAtyp(req []byte) (uint16, error) {
	candidate := uint16(req[atypPos])
	if candidate != shared.ATYP_IPV4 && candidate != shared.ATYP_FQDN && candidate != shared.ATYP_IPV6 {
		return 0, &messages.InvalidAtypError{Atyp: candidate}
	}
	return candidate, nil
}

// Deserializes the address type, the address its self and the port
func (cmd *CommandRequest) deserializeDstAddrAndPort(req []byte) error {
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
	cmd.DST_ADDR = destAddr
	cmd.DST_PORT = uint16(req[dstAddrStartPos+nextByte])<<8 | uint16(req[dstAddrStartPos+nextByte+1])
	return nil
}
