package socks5

import (
	"fmt"
	"io"
	"log"
	"net"
	"socks/pkg/socks5/auth"
	"socks/pkg/socks5/proxy"
	"socks/pkg/socks5/requests"
	"socks/pkg/socks5/responses"
	"socks/pkg/socks5/shared"
)

type SOCKS5_SERVER struct {
	tcpListener net.Listener
}

func Start(address string) error {
	server := &SOCKS5_SERVER{}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	server.tcpListener = listener
	server.handleConnect()
	return nil
}

func (server *SOCKS5_SERVER) Stop() error {
	return server.tcpListener.Close()
}

func (server *SOCKS5_SERVER) handleConnect() {
	for {
		conn, err := server.tcpListener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go func() {
			startConnection(conn)
		}()
	}
}

func startConnection(conn net.Conn) {
	connReq, err := initConnection(conn)
	if err != nil {
		notifyClientForError(responses.GENERIC_SERVER_FAILURE, conn)
		return
	}
	err = handleAuth(conn, connReq)
	if err != nil {
		notifyClientForError(responses.GENERIC_SERVER_FAILURE, conn)
		return
	}
	proxyReq, proxyErr := receiveProxyRequest(conn)
	if proxyErr != 0 {
		notifyClientForError(proxyErr, conn)
		return
	}
	resp, err, errors := handleRequest(proxyReq, conn)
	if err != nil {
		notifyClientForError(responses.GoErrorToSocksError(err), conn)
		return
	}
	_, err = conn.Write(resp)
	if err != nil {
		notifyClientForError(responses.GENERIC_SERVER_FAILURE, conn)
		conn.Close()
		return
	}
	for {
		select {
		case errDuringProxying := <-errors:
			notifyClientForError(responses.GoErrorToSocksError(errDuringProxying), conn)
			conn.Close()
		}
	}
}
func handleAuth(conn io.ReadWriter, connReq *requests.ConnectRequest) error {
	var response []byte
	var err error
	if !connReq.Contains(auth.NO_AUTH) {
		response, err = responses.SelectionMessage(auth.NO_ACCEPTABLE_METHOD)
		if err != nil {
			return err
		}
	} else {
		response, err = responses.SelectionMessage(auth.NO_AUTH)
		if err != nil {
			return err
		}
	}
	_, err = conn.Write(response)
	return err
}
func initConnection(conn io.ReadWriter) (*requests.ConnectRequest, error) {
	initReq := make([]byte, 10)
	_, err := conn.Read(initReq)
	if err != nil {
		return &requests.ConnectRequest{}, err
	}
	return requests.NewConnectRequest(initReq)
}
func handleRequest(connReq *requests.ProxyRequest, client net.Conn) ([]byte, error, chan error) {
	errors := make(chan error)
	addr := fmt.Sprintf("%s:%d", connReq.DST_ADDR.Value, connReq.DST_PORT)
	var errProxy error
	if connReq.CMD == requests.CONNECT {
		errProxy = proxy.StartProxy(client, addr, errors)
	} else if connReq.CMD == requests.UDP_ASSOCIATE {
		errProxy = proxy.StartProxyUdp(client, addr, errors)
	} else {
		errProxy = proxy.StartProxyBind(client, addr, errors)
	}
	if errProxy != nil {
		return []byte{}, errProxy, nil
	}
	// fixme: when the command is bind two responses must be send

	if connReq.CMD == requests.CONNECT {
		response, err := responses.NewSucceeded(connReq.ATYP, &connReq.DST_ADDR, connReq.DST_PORT).ToBinary()
		return response, err, errors
	} else if connReq.CMD == requests.UDP_ASSOCIATE {
		response, err := responses.NewSucceeded(connReq.ATYP, &shared.DstAddr{Type: shared.ATYP_IPV4, Value: "127.0.0.1"}, 8878).ToBinary()
		return response, err, errors
	} else {
		response, err := responses.NewSucceeded(connReq.ATYP, &connReq.DST_ADDR, connReq.DST_PORT).ToBinary()
		client.Write(response)
		response, err = responses.NewSucceeded(connReq.ATYP, &shared.DstAddr{Type: shared.ATYP_IPV4, Value: "127.0.0.1"}, 8877).ToBinary()
		return response, err, errors
	}
}

func notifyClientForError(errType int, client io.ReadWriter) {
	var res []byte
	res, err := responses.NewFailureBinary(errType)
	if err != nil {
		res = responses.NewGenericServerFailureBinary()

	}
	if _, err = client.Write(res); err != nil {
		log.Fatalf("Failed writing error to client: %v", err)
	}
}

func receiveProxyRequest(conn io.ReadWriter) (*requests.ProxyRequest, int) {
	proxyRequest := make([]byte, 1024)
	_, err := conn.Read(proxyRequest)
	if err != nil {
		log.Printf("Error reading proxy request: %v", err)
		return nil, responses.GENERIC_SERVER_FAILURE
	}
	return requests.NewProxyRequest(proxyRequest)
}
