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
	go func() {
		_, err := io.Copy(proxy.server, proxy.client)
		if err != nil {
			errors <- err
		}
	}()
	_, err := io.Copy(proxy.client, proxy.server)
	if err != nil {
		errors <- err
	}
	return nil
}

func (proxy *TCPProxy) Stop() {
	proxy.server.Close()
	proxy.client.Close()
}
