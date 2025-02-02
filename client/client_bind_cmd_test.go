package client

import (
	"context"
	"net"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
	"strconv"
	"testing"
	"time"
)

const serverRequestResponse = "OK"
const serverResponse = "TEST"

var serverConnectedWithPort uint16 // this variable holds the local port used when connecting from the server, to the proxy server

func Test_Client_Bind(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	mockTcp := startServer()
	addr := mockTcp.(*net.TCPAddr).IP.String()
	port := mockTcp.(*net.TCPAddr).Port
	connectClient := openConnectCmd(ctx, addr, uint16(port))
	// auth
	bindClient, err := NewSocks5Client(ctx, "127.0.0.1:1080")

	err = bindClient.Connect([]uint16{shared.NoAuthRequired})
	if err != nil {
		t.Fatalf("Failed sending authentication request. Reason %v", err)
	}
	if bindClient.State() != Authenticated {
		t.Fatalf("Failed authentication")
	}
	// Bind
	addrProxy, portProxy, err := bindClient.BindRequest(addr, uint16(port))
	if err != nil {
		t.Fatalf("Failed sending bind request to Dante. Reason %v", err)
	}
	// Send proxy address to server
	rwConn, err := connectClient.GetReaderWriter()
	rwConn.Write([]byte(addrProxy + ":" + strconv.Itoa(int(portProxy))))
	buf := make([]byte, 1024)
	n, err := rwConn.Read(buf)
	if string(buf[:n]) != serverRequestResponse {
		t.Fatalf("Expected server to response with %v", serverRequestResponse)
	}
	// wait for server to connect to proxy and read data
	rwBind, err := bindClient.GetReaderWriter()
	// Once the server established connection with the proxy server, before sending the data, the proxy server sends information about the connection
	buf2 := make([]byte, 1024)
	n, err = rwBind.Read(buf2)
	msg := command_response.CommandResponse{}
	msg.Deserialize(buf2[:n])
	if msg.BND_PORT != serverConnectedWithPort {
		t.Fatalf("Expected port send by the proxy server: %v, expected: %v", msg.BND_PORT, serverConnectedWithPort)
	}
	if msg.BND_ADDR.Value != "127.0.0.1" {
		t.Fatalf("Expected IP end by the proxy server: %v, expected: 127.0.0.1", msg.BND_PORT)
	}
	// Read actual data
	n, err = rwBind.Read(buf2)
	if string(buf2[:n]) != serverResponse {
		t.Fatalf("Expected server to send with TEST, got %v", string(buf2[:n]))
	}
}

func openConnectCmd(ctx context.Context, addr string, port uint16) *Socks5Client {
	client, err := NewSocks5Client(ctx, "127.0.0.1:1080")
	if err != nil {
		panic(err)
	}
	err = client.Connect([]uint16{shared.NoAuthRequired})
	if err != nil {
		panic(err)
	}
	_, _, err = client.ConnectRequest(addr, port)
	if err != nil {
		panic(err)
	}
	return client
}

func startServer() net.Addr {
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
		conn, err := connectBackToClient(string(buf[:n]))
		if err != nil {
			panic(err)
		}
		client.Write([]byte(serverRequestResponse))
		conn.Write([]byte(serverResponse))
		conn.Close()
		client.Close()
	}()
	return srv.Addr()
}

func connectBackToClient(addr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	serverConnectedWithPort = uint16(conn.LocalAddr().(*net.TCPAddr).Port)
	return conn, nil
}
