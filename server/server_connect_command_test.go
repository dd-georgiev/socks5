package server

import (
	"context"
	"fmt"
	"socks5_server/client"
	"socks5_server/client/sockstests"
	"socks5_server/messages/shared"
	"testing"
	"time"
)

func Test_Client_Connect(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*50)
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
	addr, port := sockstests.TcpEchoServer()
	socks5client, err := client.NewSocks5Client(ctx, socks5SrvAddr)
	if err != nil {
		t.Fatal("Failed connecting to Dante")
	}
	// authenticate
	err = socks5client.Connect([]uint16{shared.NoAuthRequired})
	if err != nil {
		t.Fatalf("Failed sending authentication request. Reason %v", err)
	}
	if socks5client.State() != client.Authenticated {
		t.Fatalf("Failed authentication")
	}
	// send connect request
	_, _, err = socks5client.ConnectRequest(addr, port)
	if err != nil {
		t.Fatalf("Failed sending connect request to Dante. Reason %v", err)
	}
	if socks5client.State() != client.CommandAccepted {
		t.Fatalf("Failed sending connect request to Dante")
	}
	// send and receive data
	rw, err := socks5client.GetReaderWriter()
	if err != nil {
		t.Fatalf("%v", err)
	}
	testString := "Hello"
	_, err = rw.Write([]byte(testString))
	if err != nil {
		t.Fatalf("Failed writing to mock server, reason: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := rw.Read(buf)
	if err != nil {
		t.Fatalf("Failed reading from mock server, reason: %v", err)
	}
	if string(buf[:n]) != testString {
		t.Fatalf("Expected '%s', got '%s'", testString, string(buf[:n]))
	}
	socks5client.Close()
}
