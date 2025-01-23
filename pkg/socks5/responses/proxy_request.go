package responses

import (
	"errors"
	"socks/pkg/socks5/shared"
	"strings"
)

const SUCCESS_CODE = 0x00

type ProxyRequestResponse struct {
	Version  uint16
	Reply    uint16
	RSV      uint16
	ATYP     uint16
	BIN_ADDR shared.DstAddr
	BIN_PORT uint16
}

func (res *ProxyRequestResponse) ToBinary() ([]byte, error) {
	addr, err := res.BIN_ADDR.ToBinary()
	if err != nil {
		return nil, err
	}

	bin := make([]byte, 0)
	bin = append(bin, byte(res.Version))
	bin = append(bin, byte(res.Reply))
	bin = append(bin, byte(res.RSV))
	bin = append(bin, byte(res.ATYP))
	bin = append(bin, addr...)
	bin = append(bin, byte(res.BIN_PORT>>8), byte(res.BIN_PORT))

	return bin, nil
}

func NewSucceeded(atype uint16, addr *shared.DstAddr, port uint16) *ProxyRequestResponse {
	return &ProxyRequestResponse{
		Version:  shared.PROTOCOL_VERSION,
		Reply:    SUCCESS_CODE,
		RSV:      0x00,
		ATYP:     atype,
		BIN_ADDR: *addr,
		BIN_PORT: port,
	}
}

func NewFailure(failure int) (*ProxyRequestResponse, error) {
	if failure < 1 || failure > 8 {
		return nil, errors.New("unknown failure type")
	}
	return &ProxyRequestResponse{
		Version:  shared.PROTOCOL_VERSION,
		Reply:    uint16(failure),
		RSV:      0x00,
		ATYP:     0x01,
		BIN_ADDR: shared.DstAddr{Type: shared.ATYP_IPV4, Value: "0.0.0.0"},
		BIN_PORT: 0x00,
	}, nil
}

func NewFailureBinary(failure int) ([]byte, error) {
	if failure < 1 || failure > 8 {
		return nil, errors.New("unknown failure type")
	}
	fail := ProxyRequestResponse{
		Version:  shared.PROTOCOL_VERSION,
		Reply:    uint16(failure),
		RSV:      0x00,
		ATYP:     shared.ATYP_IPV4,
		BIN_ADDR: shared.DstAddr{Type: shared.ATYP_IPV4, Value: "0.0.0.0"},
		BIN_PORT: 0x00,
	}
	return fail.ToBinary()
}

func NewGenericServerFailureBinary() []byte {
	genericServerError := &ProxyRequestResponse{
		Version:  shared.PROTOCOL_VERSION,
		Reply:    uint16(GENERIC_SERVER_FAILURE),
		RSV:      0x00,
		ATYP:     shared.ATYP_IPV4,
		BIN_ADDR: shared.DstAddr{Type: shared.ATYP_IPV4, Value: "0.0.0.0"},
		BIN_PORT: 0x00,
	}
	binary, _ := genericServerError.ToBinary()
	return binary
}

func GoErrorToSocksError(err error) int {
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "refused"):
		return CONNECTION_REFUSED
	case strings.Contains(errStr, "network is unreachable"):
		return NETWORK_UNREACHABLE
	default:
		return HOST_UNREACHABLE
	}
}
