package server

import (
	"slices"
	"socks5_server/messages/requests/available_auth_methods"
	"socks5_server/messages/responses/accept_auth_method"
	"socks5_server/messages/shared"
)

func (session *Session) handleAuth() {
	authMethodCandidate := make([]byte, 1024)
	n, err := session.conn.Read(authMethodCandidate)
	if err != nil {
		session.setError(err)
	}

	authMethods := available_auth_methods.AvailableAuthMethods{}
	err = authMethods.Deserialize(authMethodCandidate[:n])
	if err != nil {
		session.setError(err)
	}

	msg := accept_auth_method.AcceptAuthMethod{}
	if slices.Contains(authMethods.Methods(), shared.NoAuthRequired) {
		err := msg.SetMethod(shared.NoAuthRequired)
		if err != nil {
			session.setError(err)
		}
	} else {
		err = msg.SetMethod(shared.NoAcceptableMethods)
		if err != nil {
			session.setError(err)
		}
	}

	_, err = session.conn.Write(msg.ToBytes())
	if err != nil {
		session.setError(err)
	}
	session.state = Authenticated
}
