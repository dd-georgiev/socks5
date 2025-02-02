package client

import (
	"context"
	"net"
	"socks5_server/messages/shared"
	"testing"
	"time"
)

func Test_Client_Connect(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	mockTcp := startEchoServer()
	client, err := NewSocks5Client(ctx, "127.0.0.1:1080")
	if err != nil {
		t.Fatal("Failed connecting to Dante")
	}
	err = client.Connect([]uint16{shared.NoAuthRequired})
	if err != nil {
		t.Fatalf("Failed sending authentication request. Reason %v", err)
	}
	if client.State() != Authenticated {
		t.Fatalf("Failed authentication")
	}
	addr := mockTcp.(*net.TCPAddr).IP.String()
	port := mockTcp.(*net.TCPAddr).Port
	_, _, err = client.ConnectRequest(addr, uint16(port))
	if err != nil {
		t.Fatalf("Failed sending connect request to Dante. Reason %v", err)
	}
	if client.State() != CommandAccepted {
		t.Fatalf("Failed sending connect request to Dante")
	}
	rw, err := client.GetReaderWriter()
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
	client.Close()
}

func startEchoServer() net.Addr {
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
	return srv.Addr()
}
