package server

import (
	"fmt"
	"net"
	"reflect"
	"socks5_server/messages/requests/available_auth_methods"
	"socks5_server/messages/responses/accept_auth_method"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
	"testing"
)

func Fuzz_Server_Preauth(f *testing.F) {
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
	f.Add([]byte{0x05})
	f.Fuzz(func(t *testing.T, data []byte) {
		if reflect.DeepEqual(data, []byte("")) {
			t.Skip("Skipping test because of empty data, the connections will hang on read")
		}
		tcpAddr, err := net.ResolveTCPAddr("tcp", socks5SrvAddr)
		if err != nil {
			t.Fatal(err)
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			t.Fatal(err)
		}
		_, err = conn.Write(data)
		if err != nil {
			t.Fatalf("Failed to write with error %v and data: %v", err, data)
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		respMsg := accept_auth_method.AcceptAuthMethod{}
		err = respMsg.Deserialize(buf[:n])
		if err != nil {
			t.Fatalf("Failed decoding response. Reason: %v with data: %v", err, data)
		}
		conn.Close()
	})
}

func Fuzz_Server_Postauth(f *testing.F) {
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
	f.Add([]byte{0x05})

	f.Fuzz(func(t *testing.T, data []byte) {
		if reflect.DeepEqual(data, []byte("")) {
			t.Skip("Skipping test because of empty data, the connections will hang on read")
		}
		tcpAddr, err := net.ResolveTCPAddr("tcp", socks5SrvAddr)
		if err != nil {
			t.Fatal(err)
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			t.Fatal(err)
		}
		authMsg := available_auth_methods.AvailableAuthMethods{}
		authMsg.AddMethod(shared.NoAuthRequired)
		_, err = conn.Write(authMsg.ToBytes())
		if err != nil {
			t.Fatalf("Failed to write with error %v and data: %v", err, data)
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal("Failed to authenticate")
		}
		_, err = conn.Write(data)
		if err != nil {
			t.Fatalf("Failed to write with error %v and data: %v", err, data)
		}
		buf = make([]byte, 1024)
		n, err = conn.Read(buf)
		respMsg := command_response.CommandResponse{}
		err = respMsg.Deserialize(buf[:n])
		if err != nil {
			t.Fatalf("Failed decoding response. Reason: %v with data: %v", err, data)
		}
		conn.Close()
	})
}
