# Lab 5 - Implementing SOCKS5 UDP ASSOCIATE Command
## Overview
In this lab the client will be extended to support the UDP ASSOCIATE Command

## Concepts
### UDP Associate command
Unlike the other commands, this one uses UDP. Before the command is accepted, the socks5 server must start a UDP listener, which will forward the datagram to the desired UDP server. Here the client is expected to encapsulate the data in a specific format, before sending it to the proxy server. The encapsulation format was observed in Lab 1.
## Experiments
### Handling UDP Associate command
```go
func (client *Socks5Client) UDPAssociateRequest(addr string, port uint16) (string, uint16, error) {
	if client.state != Authenticated {
		return "", 0, errors.New("client is not authenticated")
	}

	err := client.constructAndSendCommand(shared.UDP_ASSOCIATE, addr, shared.ATYP_IPV4, port)
	if err != nil {
		client.setError(err)
		return "", 0, err
	}

	addrProxy, portProxy, err := client.handleCommandResponse()
	if err != nil {
		client.setError(err)
		return "", 0, err
	}

	return addrProxy, portProxy, nil
}
func (client *Socks5Client) constructAndSendCommand(cmdType uint16, addr string, addrType uint16, port uint16) error {
	req, err := constructCommand(cmdType, addr, addrType, port)
	if err != nil {
		return err
	}

	_, err = client.tcpConn.Write(req)
	if err != nil {
		return err
	}

	client.setState(CommandRequested)
	return nil
}
```
### Mock UDP Server
```
func UdpEchoServer() (string, uint16) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9999") // Probably not a terrible idea to randomly pick address 
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)

	go func() {
		for {
			buf := make([]byte, 1024)
			n, addr, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				fmt.Println(err)
				return
			}

			conn.WriteToUDP(buf[0:n], addr)
		}
	}()
	srvAddr := addr.IP.String()
	srvPort := addr.Port
	return srvAddr, uint16(srvPort)
}
```
### Test
```go

const dataSendToUDPEcho = "HELLO_RANDOM"

func Test_Client_UDP_associate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	addr, port := sockstests.UdpEchoServer()
	client, err := NewSocks5Client(ctx, "127.0.0.1:1080")
	if err != nil {
		t.Fatal("Failed connecting to Dante")
	}
	// authenticate
	err = client.Connect([]uint16{shared.NoAuthRequired})
	if err != nil {
		t.Fatalf("Failed sending authentication request. Reason %v", err)
	}
	if client.State() != Authenticated {
		t.Fatalf("Failed authentication")
	}
	// send udp associate command request
	srvIp, srvPort, err := client.UDPAssociateRequest("0.0.0.0", 0)
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
	// ^^^^^^^^^^ important
	msg := udp.UDPDatagram{}
	msg.Frag = 0
	msg.DST_ADDR = shared.DstAddr{Value: addr, Type: shared.ATYP_IPV4}
	msg.DST_PORT = uint16(port)
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

```