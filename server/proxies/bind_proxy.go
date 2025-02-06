package proxies

import (
	"io"
	"net"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
)

// A proxy which starts a listener and any accepted traffic is send to the client.
// When a connection is accepted, the client is firstly notified about it via standard CommandResponse message
// The ListeningPort and ListeningIp fields are used to return a message back to the client.
type BindProxy struct {
	server        net.Listener
	client        io.ReadWriteCloser
	ListeningPort uint16
	ListeningIp   string
}

func NewBindProxy(client io.ReadWriteCloser, addr string) (*BindProxy, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:8877")
	if err != nil {
		return nil, err
	}

	return &BindProxy{client: client, server: listener, ListeningPort: 8877, ListeningIp: "127.0.0.1"}, nil
}
func (proxy *BindProxy) Start(errors chan error) error {
	go func() {
		in, err := proxy.server.Accept()
		if err != nil {
			errors <- err
			return
		}

		err = proxy.notifyClientAboutIncomingConnection(in)
		if err != nil {
			errors <- err
			return
		}

		SpliceConnections(in, proxy.client, errors)
	}()

	return nil
}

func (proxy *BindProxy) Stop() {
	proxy.server.Close()
	proxy.client.Close()
}
func (proxy *BindProxy) notifyClientAboutIncomingConnection(in net.Conn) error {
	addr := in.RemoteAddr().(*net.TCPAddr).IP.String()
	port := in.RemoteAddr().(*net.TCPAddr).Port
	reqSourceMsg := command_response.CommandResponse{Status: command_response.Success, BND_ADDR: shared.DstAddr{Value: addr, Type: shared.ATYP_IPV4}, BND_PORT: uint16(port)}

	bytes, err := reqSourceMsg.ToBytes()
	_, err = proxy.client.Write(bytes)

	return err
}
