package requests

import (
	"errors"
	"socks/pkg/socks5/auth"
	"socks/pkg/socks5/shared"
)

type ConnectRequest struct {
	Version  uint16
	NMethods uint16
	Methods  []uint16
}

// fixme: return error codes suitable for socks5 protocol and not golang errors
// fixme: turn all the ifs into stateful compotations which modify a parameter to remove the if-err bloat
func NewConnectRequest(req []byte) (*ConnectRequest, error) {
	initReq := ConnectRequest{}
	initReq.Version = uint16(req[0])
	if initReq.Version != shared.PROTOCOL_VERSION {
		return nil, errors.New("unsupported protocol version or malformed request")
	}
	initReq.NMethods = uint16(req[1])
	if initReq.NMethods == 0 {
		return nil, errors.New("invalid auth method - auth methods count cannot be 0")
	}
	for i := 0; i < int(initReq.NMethods); i += 1 {
		method := uint16(req[2+i])
		if !auth.IsValidAuthMethod(method) {
			return nil, errors.New("invalid auth method %d")
		}
		initReq.Methods = append(initReq.Methods, method)
	}
	return &initReq, nil
}

func (req *ConnectRequest) Contains(method uint16) bool {
	for _, m := range req.Methods {
		if m == method {
			return true
		}
	}
	return false
}
