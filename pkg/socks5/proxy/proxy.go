package proxy

import (
	"io"
	"log"
	"net"
	"time"
)

// TODO: Either turn those into proper pure functions or introduce proper state management for the connections
func StartProxy(client io.ReadWriter, addr string, errors chan error) error {
	server, err := net.DialTimeout("tcp", addr, time.Duration(5)*time.Second) // fixme: move timeout(5 seconds) to config
	// fixme: server is not closed
	if err != nil {
		log.Println(err)
		return err
	}
	go spliceConnections(client, server, errors)
	return nil
}

func spliceConnections(client io.ReadWriter, server io.ReadWriter, errors chan error) {
	go func() {
		_, err := io.Copy(server, client)
		if err != nil {
			errors <- err
		}
	}()
	_, err := io.Copy(client, server)
	if err != nil {
		errors <- err
	}
}

func StartProxyBind(client io.ReadWriter, addr string, errors chan error) error {
	server, err := net.DialTimeout("tcp", addr, time.Duration(5)*time.Second) // fixme: move timeout(5 seconds) to config
	// fixme: server is not closed
	if err != nil {
		log.Println(err)
		return err
	}
	go func() {
		listener, _ := net.Listen("tcp", "127.0.0.1:8877")

		go spliceConnections(client, server, errors)
		in, _ := listener.Accept()
		go spliceConnections(in, client, errors)
	}()

	return nil
}
