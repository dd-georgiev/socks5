package client

import (
	"context"
	"socks5_server/client/sockstests"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
	"strconv"
	"testing"
	"time"
)

const serverRequestResponse = "OK"
const serverResponse = "TEST"

// this variable holds the local port used when connecting from the server, to the proxy server

func Test_Client_Bind(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	var serverConnectedWithPort uint16 // this variable is altered by the BindServer to indicate which local port is being used for the server->proxy->client connection
	addr, port := sockstests.BindServer(serverRequestResponse, serverResponse, &serverConnectedWithPort)

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
	// send bind command
	addrProxy, portProxy, err := bindClient.BindRequest(addr, uint16(port))
	if err != nil {
		t.Fatalf("Failed sending bind request to Dante. Reason %v", err)
	}
	// Send proxy address to server, so that the server can open connection to it
	// and transmit data back to the client. This is sent via the TCP session from the CONNECT command
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
