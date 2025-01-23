package shared

import (
	"errors"
	"fmt"
	"net"
)

const (
	ATYP_IPV4 = 0x01
	ATYP_FQDN = 0x03
	ATYP_IPV6 = 0x04
)

const PROTOCOL_VERSION = 0x05

type DstAddr struct {
	Type  int
	Value string
}

func (addr *DstAddr) ToBinary() ([]byte, error) {

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
func (addr *DstAddr) DstAddrFromBytes(req_bytes []byte, addrType int) (int, error) {
	switch addrType {
	case ATYP_IPV4:
		ipv4Addr := fmt.Sprintf("%d.%d.%d.%d", req_bytes[4], req_bytes[5], req_bytes[6], req_bytes[7])
		addr.Type = ATYP_IPV4
		addr.Value = ipv4Addr
		return 8, nil
	case ATYP_FQDN:
		fqdnSize := int(req_bytes[4])
		fqdnOffset := 5 // since the size doesn't include the offset from the beginning of the request, we compute it
		addr.Type = ATYP_FQDN
		addr.Value = string(req_bytes[5 : fqdnOffset+fqdnSize])
		return fqdnOffset + fqdnSize, nil
	case ATYP_IPV6:
		format := "%x:%x:%x:%x:%x:%x:%x:%x"
		ip := fmt.Sprintf(format, req_bytes[4:6], req_bytes[6:8], req_bytes[8:10], req_bytes[10:12], req_bytes[12:14], req_bytes[14:16], req_bytes[16:18], req_bytes[18:20])
		addr.Type = ATYP_IPV6
		addr.Value = ip
		return 20, nil
	}
	return 0, errors.New("invalid request - ATYP is unknown")
}
