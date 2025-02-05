package server

import (
	"context"
	"fmt"
	"socks5_server/client"
	"socks5_server/client/sockstests"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
	"strconv"
	"testing"
	"time"
)

const serverRequestResponse = "OK"
const serverResponse = "TEST"

func Test_Client_Bind(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
	var serverConnectedWithPort uint16 // this variable is altered by the BindServer to indicate which local port is being used for the server->proxy->client connection
	addr, port := sockstests.BindServer(serverRequestResponse, serverResponse, &serverConnectedWithPort)

	connectClient := openConnectCmd(ctx, proxyAddr, uint16(proxyPort), addr, port)
	bindClient, err := client.NewSocks5Client(ctx, socks5SrvAddr)
	if err != nil {
		t.Fatal("Failed connecting to Dante")
	}
	// authenticate
	err = bindClient.Connect([]uint16{shared.NoAuthRequired})
	if err != nil {
		t.Fatalf("Failed sending authentication request. Reason %v", err)
	}
	if bindClient.State() != client.Authenticated {
		t.Fatalf("Failed authentication")
	}
	// send bind request
	addrProxy, portProxy, err := bindClient.BindRequest(addr, port)
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
	err = msg.Deserialize(buf2[:n])
	if err != nil {
		t.Fatalf("Failed deserializing command. Reason %v", err)
	}
	if msg.BND_PORT != serverConnectedWithPort {
		t.Fatalf("Expected port send by the proxy server: %v, expected: %v", msg.BND_PORT, serverConnectedWithPort)
	}
	if msg.BND_ADDR.Value != "127.0.0.1" {
		t.Fatalf("Expected IP end by the proxy server: %v, expected: 127.0.0.1", msg.BND_PORT)
	}
	// Read actual data
	//// NOTE: Sometimes, the response from the successful connection and the data its self come in the same read operation
	//// Because there is no way to differentiate between them specified in the RFC, the simplest yet reliable way I found is to check if
	//// the message was the whole buffer, if not then at least part of the data arrived with the socks5 message.
	//// this is happening only with the golang server implementation and not with Dante or at least didn't happen during the 100 times I executed this test against Dante
	respLength, _ := msg.ToBytes()
	if n > len(respLength) {
		if string(buf2[len(respLength):n]) != serverResponse {
			t.Fatalf("Expected server to send with TEST, got %v", string(buf2[len(respLength):n]))
		}
		return
	}
	n, err = rwBind.Read(buf2)
	if string(buf2[:n]) != serverResponse {
		t.Fatalf("Expected server to send with TEST, got %v", string(buf2[:n]))
	}
	bindClient.Close()
}
func openConnectCmd(ctx context.Context, ip string, port uint16, dstSrv string, dstPort uint16) *client.Socks5Client {
	addr := fmt.Sprintf("%s:%d", ip, port)
	clt, err := client.NewSocks5Client(ctx, addr)
	if err != nil {
		panic(err)
	}
	err = clt.Connect([]uint16{shared.NoAuthRequired})
	if err != nil {
		panic(err)
	}
	_, _, err = clt.ConnectRequest(dstSrv, dstPort)
	if err != nil {
		panic(err)
	}
	return clt
}
