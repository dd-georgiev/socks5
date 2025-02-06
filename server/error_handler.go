package server

import (
	"socks5_server/messages/responses/accept_auth_method"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
)

func (session *Session) RespondToClientDependingOnState() {
	switch session.state {
	case PendingAuthMethods:
		noAcceptableMethodMsg := accept_auth_method.AcceptAuthMethod{}
		noAcceptableMethodMsg.SetMethod(shared.NoAcceptableMethods)
		session.conn.Write(noAcceptableMethodMsg.ToBytes())
		session.conn.Close()
	case Authenticated:
		srvFailure := command_response.CommandResponse{}
		srvFailure.Status = command_response.SocksServerFailure
		srvFailure.BND_ADDR = shared.DstAddr{Value: "0.0.0.0", Type: shared.ATYP_IPV4}
		srvFailure.BND_PORT = 0
		srvFailureBytes, err := srvFailure.ToBytes()
		if err != nil {
			session.conn.Close()
		}
		session.conn.Write(srvFailureBytes)
	}
	if session.state == Proxying {
		session.conn.Close()
	}
}
