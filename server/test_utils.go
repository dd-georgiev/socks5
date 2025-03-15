package server

import (
	"net"
)

func startSocks5Server() (string, int) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	//	srv := Socks5Server{Listener: listener}
	go Start(listener)
	addr := listener.Addr().(*net.TCPAddr).IP.String()
	port := listener.Addr().(*net.TCPAddr).Port
	return addr, port
}
