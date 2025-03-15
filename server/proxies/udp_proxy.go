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
