package requests

import (
	"socks/pkg/socks5/responses"
	"socks/pkg/socks5/shared"
)

const (
	CONNECT       = 1
	BIND          = 2
	UDP_ASSOCIATE = 3
)

type ProxyRequest struct {
	Version  uint16
	CMD      uint16
	RSV      uint16
	ATYP     uint16
	DST_ADDR shared.DstAddr
	DST_PORT uint16
}

/*
Format:
+------+-------+---------+-------+------------+----------+
| VER  |  CMD  |   RSV   |  ATYP |  DST.ADDR  | DST.PORT |
+------+-------+---------+-------+------------+----------+
*/

// fixme: turn all the ifs into stateful compotations which modify a parameter to remove the if-err bloat
func NewProxyRequest(req []byte) (*ProxyRequest, int) {
	initReq := ProxyRequest{}
	initReq.Version = uint16(req[0])
	if initReq.Version != shared.PROTOCOL_VERSION {
		return nil, responses.GENERIC_SERVER_FAILURE
	}
	initReq.CMD = uint16(req[1])
	if initReq.CMD < 1 || initReq.CMD > 3 {
		return nil, responses.COMMAND_NOT_SUPPORTED
	}
	initReq.RSV = uint16(req[2])
	if initReq.RSV != 0x00 {
		return nil, responses.GENERIC_SERVER_FAILURE
	}
	initReq.ATYP = uint16(req[3])
	if initReq.ATYP < 0 && initReq.ATYP > 3 { // make configurable
		return nil, responses.ADDRESS_TYPE_NOT_SUPPORTED
	}
	destAddr := shared.DstAddr{}
	nextByte, err := destAddr.DstAddrFromBytes(req, int(initReq.ATYP))
	initReq.DST_ADDR = destAddr
	if err != nil {
		return nil, responses.GENERIC_SERVER_FAILURE
	}

	initReq.DST_PORT = uint16(req[nextByte])<<8 | uint16(req[nextByte+1])
	return &initReq, 0
}
