package proxies

import (
	"io"
	"net"
	"time"
)

type TCPProxy struct {
	server io.ReadWriteCloser
	client io.ReadWriteCloser
}

func NewConnectProxy(addr string, client io.ReadWriteCloser) (*TCPProxy, error) {
	server, err := net.DialTimeout("tcp", addr, time.Duration(1)*time.Second)
	if err != nil {
		return nil, err
	}

	return &TCPProxy{server: server, client: client}, nil
}

func (proxy *TCPProxy) Start(errors chan error) error {
	SpliceConnections(proxy.server, proxy.client, errors)
	return nil
}

func (proxy *TCPProxy) Stop() {
	proxy.server.Close()
	proxy.client.Close()
}
