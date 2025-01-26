package proxy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"socks/pkg/socks5/shared"
	"time"
)

// TODO: Either turn those into proper pure functions or introduce proper state management for the connections
// TODO: Split into 3 different proxy classes - CONNECT, BIND and ASSOCIATE
func StartProxy(client io.ReadWriter, addr string, errors chan error) error {
	server, err := net.DialTimeout("tcp", addr, time.Duration(5)*time.Second) // fixme: move timeout(5 seconds) to config
	// fixme: server is not closed
	if err != nil {
		log.Println(err)
		return err
	}
	go spliceConnections(client, server, errors)
	return nil
}

func spliceConnections(client io.ReadWriter, server io.ReadWriter, errors chan error) {
	go func() {
		_, err := io.Copy(server, client)
		if err != nil {
			errors <- err
		}
	}()
	_, err := io.Copy(client, server)
	if err != nil {
		errors <- err
	}
}

func StartProxyBind(client io.ReadWriter, addr string, errors chan error) error {
	server, err := net.DialTimeout("tcp", addr, time.Duration(5)*time.Second) // fixme: move timeout(5 seconds) to config
	// fixme: server is not closed
	if err != nil {
		log.Println(err)
		return err
	}
	go func() {
		listener, _ := net.Listen("tcp", "127.0.0.1:8877")

		go spliceConnections(client, server, errors)
		in, _ := listener.Accept()
		go spliceConnections(in, client, errors)
	}()

	return nil
}

// fixme: for some reason this work only for single client<->server exchange
func StartProxyUdp(_ io.ReadWriter, addr string, errors chan error) error {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8878")
	if err != nil {
		return err
	}
	udpServer, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	go func() {
		for {
			var buf = make([]byte, 1024)
			readBytes, addrClient, err := udpServer.ReadFromUDP(buf[0:])
			if err != nil {
				fmt.Println(err)
				return
			}
			dgram, err := NewFromRequest(buf, readBytes)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(dgram)
			addr := fmt.Sprintf("%s:%d", dgram.addr.Value, dgram.port)
			udpAddrSrv, err := net.ResolveUDPAddr("udp", addr)

			if err != nil {
				fmt.Println(err)
			}
			conn, err := net.DialUDP("udp", nil, udpAddrSrv)
			_, err = conn.Write(dgram.data)
			data, err := bufio.NewReader(conn).ReadString('\n')
			fmt.Println(string(data))
			dgram.data = []byte(data) // altering this dgram is not great idea
			resp, _ := dgram.ToBytes()
			udpServer.WriteToUDP(resp, addrClient)
		}
	}()
	return nil
}

// fix me - private fields, bad naming, bad file
type UpdDatagramRequest struct {
	rsv   uint16
	frag  uint16
	atype uint16
	addr  shared.DstAddr
	port  uint16
	data  []byte
}

func NewFromRequest(buf []byte, size int) (*UpdDatagramRequest, error) {
	dgram := &UpdDatagramRequest{}
	rsv := uint16(buf[0])<<8 | uint16(buf[1])
	dgram.rsv = rsv
	frag := uint16(buf[2])
	dgram.frag = frag
	atyp := uint16(buf[3])
	dgram.atype = atyp
	addr := shared.DstAddr{}
	bytes, err := addr.DstAddrFromBytes(buf, int(atyp))
	if err != nil {
		return nil, err
	}
	dgram.addr = addr
	port := uint16(buf[bytes])<<8 | uint16(buf[bytes+1])
	dgram.port = port
	dgram.data = buf[bytes+2 : size]
	return dgram, nil
}

func (dgram *UpdDatagramRequest) ToBytes() ([]byte, error) {
	addr, err := dgram.addr.ToBinary()
	if err != nil {
		return nil, err
	}

	bin := make([]byte, 0)
	bin = append(bin, byte(dgram.rsv>>8), byte(dgram.rsv))
	bin = append(bin, byte(dgram.frag))
	bin = append(bin, byte(dgram.atype))
	bin = append(bin, addr...)
	bin = append(bin, byte(dgram.port>>8), byte(dgram.port))
	bin = append(bin, dgram.data...)

	return bin, nil
}
