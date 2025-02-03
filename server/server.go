package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"slices"
	"socks5_server/messages/requests/available_auth_methods"
	"socks5_server/messages/requests/command_request"
	"socks5_server/messages/responses/accept_auth_method"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
	"socks5_server/server/proxies"
)

type SessionState uint16

const PendingAuthMethods SessionState = 10
const Authenticated SessionState = 20
const Proxying SessionState = 30

type Session struct {
	state SessionState
	conn  net.Conn
	err   error
}

func (session *Session) setError(err error) {
	session.err = err
}
func Start(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		session := Session{state: PendingAuthMethods, conn: conn}
		go session.handler()
	}
}

func (session *Session) handler() {
	for {
		switch session.state {
		case PendingAuthMethods:
			session.handleAuth()
		case Authenticated:
			session.handleCmd()
		case Proxying:
			return
		}
	}
}

func (session *Session) handleAuth() {
	authMethodCandidate := make([]byte, 1024)
	n, err := session.conn.Read(authMethodCandidate)
	if err != nil {
		log.Fatal(err)
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
	}
	err = msg.SetMethod(shared.NoAcceptableMethods)
	if err != nil {
		session.setError(err)
	}
	session.conn.Write(msg.ToBytes())
	session.state = Authenticated
}

func (session *Session) handleCmd() {
	commandCandidate := make([]byte, 1024)
	n, err := session.conn.Read(commandCandidate)
	if err != nil {
		log.Fatal(err)
	}

	cmd := command_request.CommandRequest{}
	err = cmd.Deserialize(commandCandidate[:n])
	if err != nil {
		session.setError(err)
	}

	resp := command_response.CommandResponse{}
	if cmd.CMD != shared.CONNECT {
		resp.Status = command_response.CommandNotSupported
		bytes, err := resp.ToBytes()
		if err != nil {
			session.setError(err)
		}
		session.conn.Write(bytes)
		session.setError(errors.New("unsupported command"))
		return
	}

	proxyErrors := make(chan error)
	remoteAddr := fmt.Sprintf("%s:%d", cmd.DST_ADDR.Value, cmd.DST_PORT)
	proxy, err := proxies.NewConnectProxy(remoteAddr, session.conn)
	if err != nil {
		session.setError(err)
	}
	go proxy.Start(proxyErrors) // for some reason doesn't proxy
	if err != nil {
		session.setError(err)
	}
	go session.proxyErrorHandler(proxyErrors, proxy)

	resp.Status = command_response.Success
	resp.BND_PORT = 0
	resp.BND_ADDR = shared.DstAddr{Value: "0.0.0.0", Type: shared.ATYP_IPV4}
	bytes, _ := resp.ToBytes()
	session.conn.Write(bytes)
	session.state = Proxying
}
func (session *Session) proxyErrorHandler(errors chan error, proxy *proxies.ConnectProxy) {
	err := <-errors
	session.setError(err)
	proxy.Stop()
}
