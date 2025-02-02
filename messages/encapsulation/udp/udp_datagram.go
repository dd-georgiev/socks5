package udp

// Represents a UDP datagram encapsulation, used by the client and the server when communicating to remote server.
// This encapsulation is used by both the client and the server
import (
	"socks5_server/messages"
	"socks5_server/messages/shared"
)

var rsvVal = []byte{0x00, 0x00}

const rsvStart = 0
const rsvEnd = 1
const fragPos = 2
const atypPos = 3
const dstAddrStartPos = 4

type UDPDatagram struct {
	Frag         uint16
	DST_ADDR     shared.DstAddr
	DST_PORT     uint16
	DATA         []byte
	dataStartPos int
}

// Deserialize converts bytes received over the wire to UDPDatagram
func (dgram *UDPDatagram) Deserialize(req []byte) error {
	if len(req) < 10 {
		return messages.MalformedMessageError{}
	}
	if req[rsvStart] != 0 || req[rsvEnd] != 0 {
		invalidRsvVal := uint16(req[rsvStart])<<8 | uint16(req[rsvEnd])
		return messages.InvalidReservedFieldError{Value: invalidRsvVal}
	}

	dgram.Frag = uint16(req[fragPos])
	if err := dgram.deserializeDstAddrAndPort(req); err != nil {
		return err
	}

	dgram.DATA = req[dgram.dataStartPos:]
	return nil
}

// ToBytes converts the UDPDatagram to wire-transferable byte array
func (dgram *UDPDatagram) ToBytes() ([]byte, error) {
	addrBytes, err := dgram.DST_ADDR.ToBytes()
	if err != nil {
		return []byte{}, err
	}
	res := make([]byte, 0)
	res = append(res, rsvVal...)
	res = append(res, byte(dgram.Frag))
	res = append(res, byte(dgram.DST_ADDR.Type))
	res = append(res, addrBytes...)
	res = append(res, byte(dgram.DST_PORT>>8), byte(dgram.DST_PORT))
	res = append(res, dgram.DATA...)
	return res, nil
}

func (dgram *UDPDatagram) deserializeDstAddrAndPort(req []byte) error {
	atyp, err := dgram.deserializeAtyp(req)
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
	dgram.DST_ADDR = destAddr
	dgram.DST_PORT = uint16(req[dstAddrStartPos+nextByte])<<8 | uint16(req[dstAddrStartPos+nextByte+1])
	dgram.dataStartPos = dstAddrStartPos + nextByte + 2 // two because the 1st byte is part of the port, so we need the next one
	return nil
}
func (dgram *UDPDatagram) deserializeAtyp(req []byte) (uint16, error) {
	candidate := uint16(req[atypPos])
	if candidate != shared.ATYP_IPV4 && candidate != shared.ATYP_FQDN && candidate != shared.ATYP_IPV6 {
		return 0, &messages.InvalidAtypError{Atyp: candidate}
	}
	return candidate, nil
}
