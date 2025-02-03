package proxies

import (
	"io"
	"net"
	"time"
)

type ConnectProxy struct {
	server io.ReadWriteCloser
	client io.ReadWriteCloser
}

func NewConnectProxy(addr string, client io.ReadWriteCloser) (*ConnectProxy, error) {
	server, err := net.DialTimeout("tcp", addr, time.Duration(5)*time.Second)
	if err != nil {
		return nil, err
	}

	return &ConnectProxy{server: server, client: client}, nil
}

func (p *ConnectProxy) Start(errors chan error) error {
	go func() {
		_, err := io.Copy(p.server, p.client)
		if err != nil {
			errors <- err
		}
	}()
	_, err := io.Copy(p.client, p.server)
	if err != nil {
		errors <- err
	}
	return nil
}

func (p *ConnectProxy) Stop() {
	p.server.Close()
	p.client.Close()
}
