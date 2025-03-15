package sockstests

import (
	"fmt"
	"net"
)

func TcpEchoServer() (string, uint16) {
	srv, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	go func() {
		client, err := srv.Accept()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		if err != nil {
			panic(err)
		}
		client.Write(buf[:n])
		client.Close()
	}()
	addr := srv.Addr().(*net.TCPAddr).IP.String()
	port := srv.Addr().(*net.TCPAddr).Port
	return addr, uint16(port)
}

func UdpEchoServer() (string, uint16) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)

	go func() {
		for {
			buf := make([]byte, 1024)
			n, addr, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				fmt.Println(err)
				return
			}

			conn.WriteToUDP(buf[0:n], addr)
		}
	}()
	srvAddr := addr.IP.String()
	srvPort := addr.Port
	return srvAddr, uint16(srvPort)
}

func BindServer(serverRequestResponse string, serverResponse string, serverConnectedWithPort *uint16) (string, uint16) {
	srv, err := net.Listen("tcp4", "127.0.0.1:4440")
	if err != nil {
		panic(err)
	}

	go func() {
		client, err := srv.Accept()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		if err != nil {
			panic(err)
		}
		conn, err := connectBackToClient(string(buf[:n]), serverConnectedWithPort)
		if err != nil {
			panic(err)
		}
		client.Write([]byte(serverRequestResponse))
		conn.Write([]byte(serverResponse))
		conn.Close()
		client.Close()
	}()
	addr := srv.Addr().(*net.TCPAddr).IP.String()
	port := srv.Addr().(*net.TCPAddr).Port
	return addr, uint16(port)
}

func connectBackToClient(addr string, serverConnectedWithPort *uint16) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	*serverConnectedWithPort = uint16(conn.LocalAddr().(*net.TCPAddr).Port)
	return conn, nil
}
