package server

import (
	"context"
	"fmt"
	"net"
	"socks5_server/client"
	"socks5_server/client/sockstests"
	"socks5_server/messages/encapsulation/udp"
	"socks5_server/messages/shared"
	"testing"
	"time"
)

const dataSendToUDPEcho = "HELLO_RANDOM"

func Test_Client_UDP_Associate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*50)
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
	addr, port := sockstests.UdpEchoServer()
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
	// send udp associate command request
	srvIp, srvPort, err := socks5client.UDPAssociateRequest("0.0.0.0", 0)
	if err != nil {
		t.Fatalf("Failed sending UDP associate request. Reason %v", err)
	}
	// connect to UDP listener of the proxy
	udpAddrStr := fmt.Sprintf("%s:%d", srvIp, srvPort)
	udpAddr, err := net.ResolveUDPAddr("udp", udpAddrStr)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Fatalf("Failed connecting to UDP. Reason %v", err)
	}
	defer conn.Close()
	// encapsulate data and conn. info for the server
	for i := 0; i < 10; i++ {
		msg := udp.UDPDatagram{}
		msg.Frag = 0
		msg.DST_ADDR = shared.DstAddr{Value: addr, Type: shared.ATYP_IPV4}
		msg.DST_PORT = port
		msg.DATA = []byte(dataSendToUDPEcho)
		data, err := msg.ToBytes()
		if err != nil {
			t.Fatalf("Failed converting message to bytes. Reason %v", err)
		}
		conn.Write(data)
		// read response of the datagram
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatalf("Failed reading from UDP. Reason %v", err)
		}
		// deserialize and access the response
		response := udp.UDPDatagram{}
		err = response.Deserialize(buf[:n])
		if err != nil {
			t.Fatalf("Failed reading from UDP. Reason %v", err)
		}
		if string(response.DATA) != dataSendToUDPEcho {
			t.Fatalf("Response doesn't match send data. expected %v got %v", dataSendToUDPEcho, response.DATA)
		}
	}
	socks5client.Close()
}
