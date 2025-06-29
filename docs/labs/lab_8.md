- [Lab 8 - Implementing SOCKS5 Server UDP ASSOCIATE Command](#lab-8---implementing-socks5-server-udp-associate-command)
    * [Overview](#overview)
    * [Experiments](#experiments)
        + [Handling UDP ASSOCIATE command](#handling-udp-associate-command)
            - [Modifying the handleCmd method](#modifying-the-handlecmd-method)
            - [Adding UDP Associate specific handler](#adding-udp-associate-specific-handler)
        + [UDP Proxy](#udp-proxy)
        + [Modifying the test](#modifying-the-test)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

# Lab 8 - Implementing SOCKS5 Server UDP ASSOCIATE Command
## Overview
In this lab, we extend the server code to handle the `UDP ASSOCIATE` SOCKS5 command.
## Experiments
### Handling UDP ASSOCIATE command
#### Modifying the handleCmd method
```go
func (session *Session) handleCmd() {
// .................
    switch cmd.CMD {
        case shared.CONNECT:
            session.handleConnectCmd(cmd)
            return
            case shared.BIND:
            session.handleBindCmd(cmd)
            return
        case shared.UDP_ASSOCIATE: // this is the change
            session.handleUdpAssociateCmd(cmd)
            return
    }
// ...............
}
```
#### Adding UDP Associate specific handler
```go
func (session *Session) handleUdpAssociateCmd() {
	proxyErrors := make(chan error)
	proxy, err := proxies.NewUDPProxy()
	if err != nil {
		return
	}
	go proxy.Start(proxyErrors)

	resp := command_response.CommandResponse{}
	resp.Status = command_response.Success
	resp.BND_PORT = proxy.Port
	resp.BND_ADDR = shared.DstAddr{Value: session.conn.LocalAddr().(*net.TCPAddr).IP.String(), Type: shared.ATYP_IPV4}
	bytes, _ := resp.ToBytes()
	session.conn.Write(bytes)
	session.state = Proxying
	go session.proxyErrorHandler(proxyErrors, proxy)
}
```
### UDP Proxy
```go
// The biggest difference between this proxy and the other two, is that this one must handle the packet encapsulation.
// This probably can be delegate to the server & let the proxy deal only with data but it makes the code looks more complicated that it is.
package proxies

import (
	"fmt"
	"net"
	"socks5_server/messages/encapsulation/udp"
)

type UDPProxy struct {
	server *net.UDPConn
	Port   uint16
	Addr   string
}

func NewUDPProxy() (*UDPProxy, error) {
	udpServer, err := startUdpListener()
	if err != nil {
		return nil, err
	}

	port := uint16(udpServer.LocalAddr().(*net.UDPAddr).Port)
	addr := udpServer.LocalAddr().String()

	return &UDPProxy{server: udpServer, Port: port, Addr: addr}, nil
}

func (proxy *UDPProxy) Start(errors chan error) error {
	go func() {
		for {
			addrClient, dgram, err := proxy.receiveRequest()
			if err != nil {
				errors <- err
				return

			}
			// concatIpAndPort(dgram.DST_ADDR.Value, dgram.DST_PORT) translate the information from the client in a way the go std can understand. Nothing fancy here
			responseData, err := sendToRemote(dgram.DATA, concatIpAndPort(dgram.DST_ADDR.Value, dgram.DST_PORT))
			if err != nil {
				errors <- err
				return
			}

			respDgram := encapsulateResponse(dgram, responseData)
			resp, err := respDgram.ToBytes()
			if err != nil {
				errors <- err
				return
			}

			_, err = proxy.server.WriteToUDP(resp, addrClient)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
	return nil
}

func (proxy *UDPProxy) Stop() {
	proxy.server.Close()
}

func sendToRemote(data []byte, addr string) ([]byte, error) {
	udpAddrSrv, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, udpAddrSrv)
	if err != nil {
		return nil, err
	}

	n, err := conn.Write(data)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func (proxy *UDPProxy) receiveRequest() (*net.UDPAddr, *udp.UDPDatagram, error) {
	var buf = make([]byte, 65535) // max udp?
	n, addrClient, err := proxy.server.ReadFromUDP(buf)
	if err != nil {
		return nil, nil, err
	}
	dgram := udp.UDPDatagram{}
	err = dgram.Deserialize(buf[:n])
	if err != nil {
		return nil, nil, err
	}
	return addrClient, &dgram, nil
}

func encapsulateResponse(requestDatagram *udp.UDPDatagram, data []byte) *udp.UDPDatagram {
	respDgram := udp.UDPDatagram{}
	respDgram.DST_ADDR = requestDatagram.DST_ADDR
	respDgram.DST_PORT = requestDatagram.DST_PORT
	respDgram.Frag = requestDatagram.Frag
	respDgram.DATA = data
	return &respDgram
}

func concatIpAndPort(addr string, port uint16) string {
	return fmt.Sprintf("%s:%d", addr, port)
}

func startUdpListener() (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}
	udpServer, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return nil, err
	}

	return udpServer, nil
}

```
### Modifying the test
```go
func Test_Client_UDP_Associate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*50)
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
    // ......................
}
```