package server

import (
	"errors"
	"fmt"
	"net"
	"socks5_server/messages/requests/command_request"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
	"socks5_server/server/proxies"
)

func (session *Session) handleCommand() {
	commandCandidate := make([]byte, 1024)
	n, err := session.conn.Read(commandCandidate)
	if err != nil {
		session.setError(err)
	}

	cmd := command_request.CommandRequest{}
	err = cmd.Deserialize(commandCandidate[:n])
	if err != nil {
		session.setError(err)
	}

	switch cmd.CMD {
	case command_request.CONNECT:
		session.handleConnectCmd(cmd)
		return
	case command_request.BIND:
		session.handleBindCmd(cmd)
		return
	case command_request.UDP_ASSOCIATE:
		session.handleUdpAssociateCmd()
		return
	default:
		session.setError(errors.New("unknown command"))
		return
	}
}
func (session *Session) handleConnectCmd(cmd command_request.CommandRequest) {
	proxyErrors := make(chan error)
	remoteAddr := fmt.Sprintf("%s:%d", cmd.DST_ADDR.Value, cmd.DST_PORT)
	proxy, err := proxies.NewConnectProxy(remoteAddr, session.conn)
	if err != nil {
		session.setError(err)
	}
	go proxy.Start(proxyErrors)
	if err != nil {
		session.setError(err)
	}
	go session.proxyErrorHandler(proxyErrors, proxy)

	resp := command_response.CommandResponse{}
	resp.Status = command_response.Success
	resp.BND_PORT = 0
	resp.BND_ADDR = shared.DstAddr{Value: "0.0.0.0", Type: shared.ATYP_IPV4}
	bytes, _ := resp.ToBytes()
	session.conn.Write(bytes)
	session.state = Proxying
}
func (session *Session) handleUdpAssociateCmd() {
	proxyErrors := make(chan error)
	proxy, err := proxies.NewUDPProxy()
	if err != nil {
		return
	}
	go proxy.Start(proxyErrors)

	resp := command_response.CommandResponse{}
	resp.Status = command_response.Success
	resp.BND_PORT = proxy.Port
	resp.BND_ADDR = shared.DstAddr{Value: session.conn.LocalAddr().(*net.TCPAddr).IP.String(), Type: shared.ATYP_IPV4}
	bytes, _ := resp.ToBytes()
	session.conn.Write(bytes)
	session.state = Proxying
	go session.proxyErrorHandler(proxyErrors, proxy)
}
func (session *Session) handleBindCmd(cmd command_request.CommandRequest) {
	proxyErrors := make(chan error)
	remoteAddr := fmt.Sprintf("%s:%d", cmd.DST_ADDR.Value, cmd.DST_PORT)
	proxy, err := proxies.NewBindProxy(session.conn, remoteAddr)
	if err != nil {
		session.setError(err)
	}
	go proxy.Start(proxyErrors)
	if err != nil {
		session.setError(err)
	}
	go session.proxyErrorHandler(proxyErrors, proxy)

	resp := command_response.CommandResponse{}
	resp.Status = command_response.Success
	resp.BND_PORT = proxy.ListeningPort
	resp.BND_ADDR = shared.DstAddr{Value: proxy.ListeningIp, Type: shared.ATYP_IPV4}
	bytes, _ := resp.ToBytes()
	session.conn.Write(bytes)
	session.state = Proxying
}

func (session *Session) proxyErrorHandler(errors chan error, proxy proxies.Proxy) {
	err := <-errors
	session.setError(err)
	proxy.Stop()
}
