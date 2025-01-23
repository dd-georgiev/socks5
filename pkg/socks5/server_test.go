package socks5

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"testing"
)

func TestTCPListenerMustPickNoAuthIfPresented(t *testing.T) {
	expectedResp := []byte{0x05, 0x00}
	err := Start(":1081")
	if err != nil {
		t.Fatal(err)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1081")
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	req := []byte{0x05, 0x02, 0x00, 0x02}
	conn.Write(req)
	reply := make([]byte, 2)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
	}
	if !bytes.Equal(reply, expectedResp) {
		fmt.Printf("Expected: %v, actual: %v\n", expectedResp, reply)
		t.Fatal("Reply does not match expected response")
	}

}

func TestTCPListenerMustPickNoAcceptableMethodIfNoAuthIsNotPresented(t *testing.T) {
	expectedResp := []byte{0x05, 0xFF}
	err := Start(":1081")
	if err != nil {
		t.Fatal(err)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1081")
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	req := []byte{0x05, 0x02, 0x01, 0x02}
	conn.Write(req)
	reply := make([]byte, 2)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
	}
	if !bytes.Equal(reply, expectedResp) {
		fmt.Printf("Expected: %v, actual: %v\n", expectedResp, reply)
		t.Fatal("Reply does not match expected response")
	}
}
