package shared

import (
	"errors"
	"fmt"
	"net"
	"socks5_server/messages"
)

const messageFqdnLengthIndex = 1
const ipv4Format = "%d.%d.%d.%d"
const ipv4Size = 4

const ipv6Format = "%x:%x:%x:%x:%x:%x:%x:%x"
const ipv6Size = 16

// Provides name->int mapping for the different address types as defined in RFC1928
const (
	ATYP_IPV4 = 0x01
	ATYP_FQDN = 0x03
	ATYP_IPV6 = 0x04
)

type UnknownATYP struct {
	AddrType uint16
}

func (e UnknownATYP) Error() string {
	return fmt.Sprintf("Unknown ATYP %d", e.AddrType)
}

type DstAddr struct {
	Type  uint16
	Value string
}

func (addr *DstAddr) ToBytes() ([]byte, error) {

	if addr.Type == ATYP_FQDN {
		bin := make([]byte, 0)
		bin = append(bin, byte(len(addr.Value)))
		bin = append(bin, addr.Value...)
		return bin, nil
	}

	ip := net.ParseIP(addr.Value)
	if ip == nil {
		return nil, errors.New("invalid ip address " + addr.Value)
	}

	if addr.Type == ATYP_IPV6 {
		return ip.To16(), nil
	}
	return ip.To4(), nil
}

func (addr *DstAddr) Deserialize(buf []byte, addrType uint16) (int, error) {
	if len(buf) < 4 {
		return 0, messages.MalformedMessageError{}
	}
	addr.Type = addrType
	switch addrType {
	case ATYP_IPV4:
		addr.Value = fmt.Sprintf(ipv4Format, buf[0], buf[1], buf[2], buf[3])
		return ipv4Size, nil
	case ATYP_FQDN:
		fqdnSize := int(buf[0])
		if fqdnSize > len(buf) {
			return 0, messages.MalformedMessageError{}
		}
		addr.Value = string(buf[messageFqdnLengthIndex : messageFqdnLengthIndex+fqdnSize])
		return messageFqdnLengthIndex + fqdnSize, nil
	case ATYP_IPV6:
		if len(buf) < 16 {
			return 0, messages.MalformedMessageError{}
		}
		addr.Value = fmt.Sprintf(ipv6Format, buf[0:2], buf[2:4], buf[4:6], buf[6:8], buf[8:10], buf[10:12], buf[12:14], buf[14:16])
		return ipv6Size, nil
	}
	return 0, UnknownATYP{AddrType: addrType}
}
