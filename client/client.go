package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"socks5_server/messages/requests/available_auth_methods"
	"socks5_server/messages/requests/command_request"
	"socks5_server/messages/responses/accept_auth_method"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
)

type ConnectionState int

// The states in which the client may be. They are in increasing order and a client cannot transition from higher to lower state.
const (
	PendingAuthMethods          ConnectionState = 10
	ExpectingAcceptedAuthMethod ConnectionState = 20
	PendingAuthentication       ConnectionState = 30
	Authenticated               ConnectionState = 40
	CommandRequested            ConnectionState = 50
	CommandAccepted             ConnectionState = 60
	Closed                      ConnectionState = 70
	Errored                     ConnectionState = 80
)

type Socks5Client struct {
	state   ConnectionState
	tcpConn net.Conn
	err     error
}

func (client *Socks5Client) State() ConnectionState {
	return client.state
}

// Transitions the client to a new state, if the previous state allows it. This is for internal usage only. If the new state is lesser than the current one the function will panic
func (client *Socks5Client) setState(newState ConnectionState) {
	if client.state > newState {
		panic(fmt.Sprintf("cannot transition from %v to %v", client.state, newState))
	}
	client.state = newState
}

// Sets the err and state field of the client struct.
func (client *Socks5Client) setError(err error) {
	client.err = err
	client.state = Errored
}

// NewSocks5Client Creates new client bound to context and connect to given proxy server. The connection is not start with the creation!
func NewSocks5Client(ctx context.Context, servAddr string) (*Socks5Client, error) {
	conn, err := openTcpConnection(servAddr)
	if err != nil {
		return nil, err
	}
	client := &Socks5Client{}
	client.state = PendingAuthMethods
	client.tcpConn = conn
	go func() {
		select {
		case <-ctx.Done():
			client.setState(Closed)
			_ = client.tcpConn.Close()
		}
	}()
	return client, nil
}

// Connect Start initial connection ot the proxy, by sending the authentication methods supported by the client. After this method is called the handleAuth method (which expects the response with the chose auth method) is called synchronously.
func (client *Socks5Client) Connect(authMethods []uint16) error {
	aam := available_auth_methods.AvailableAuthMethods{}

	if err := aam.AddMultipleMethods(authMethods); err != nil {
		client.setError(err)
		return err
	}

	_, err := client.tcpConn.Write(aam.ToBytes())
	if err != nil {
		client.setError(err)
		return err
	}
	client.setState(ExpectingAcceptedAuthMethod)
	return client.handleAuth()
}

// ConnectRequest Send a Connect command request to the proxy server
func (client *Socks5Client) ConnectRequest(addr string, port uint16) (string, uint16, error) {
	if client.state != Authenticated {
		return "", 0, errors.New("client is not authenticated")
	}
	commandRequest := command_request.CommandRequest{}
	commandRequest.CMD = shared.CONNECT
	commandRequest.DST_ADDR = shared.DstAddr{Value: addr, Type: shared.ATYP_IPV4}
	commandRequest.DST_PORT = port
	req, err := commandRequest.ToBytes()

	if err != nil {
		return "", 0, err
	}
	_, err = client.tcpConn.Write(req)
	if err != nil {
		client.setError(err)
		return "", 0, err
	}
	client.setState(CommandRequested)
	addrProxy, portProxy, err := client.handleCommandResponse()
	if err != nil {
		client.setError(err)
		return "", 0, err
	}
	return addrProxy, portProxy, nil
}
func (client *Socks5Client) BindRequest(addr string, port uint16) (string, uint16, error) {
	if client.state != Authenticated {
		return "", 0, errors.New("client is not authenticated")
	}
	commandRequest := command_request.CommandRequest{}
	commandRequest.CMD = shared.BIND
	commandRequest.DST_ADDR = shared.DstAddr{Value: addr, Type: shared.ATYP_IPV4}
	commandRequest.DST_PORT = port
	req, err := commandRequest.ToBytes()

	if err != nil {
		return "", 0, err
	}
	_, err = client.tcpConn.Write(req)
	if err != nil {
		client.setError(err)
		return "", 0, err
	}

	client.setState(CommandRequested)
	addrProxy, portProxy, err := client.handleCommandResponse()
	if err != nil {
		client.setError(err)
		return "", 0, err
	}
	client.setState(CommandAccepted)
	return addrProxy, portProxy, err
}

func (client *Socks5Client) Close() error {
	client.setState(Closed)
	return client.tcpConn.Close()
}

// GetReaderWriter Returns io.ReadWrite after the server has accepted a command request.
func (client *Socks5Client) GetReaderWriter() (io.ReadWriter, error) {
	if client.state != CommandAccepted {
		return nil, errors.New("the server has not accepted any command")
	}
	return client.tcpConn, nil
}

// Private

func (client *Socks5Client) handleAuth() error {
	if client.state != ExpectingAcceptedAuthMethod {
		return errors.New("client is not expecting accepted auth clients")
	}
	buf := make([]byte, 64)
	_, err := client.tcpConn.Read(buf)
	if err != nil {
		client.setError(err)
		return err
	}
	acceptedMethod := accept_auth_method.AcceptAuthMethod{}
	if err := acceptedMethod.Deserialize(buf); err != nil {
		client.setError(err)
		return err
	}
	if acceptedMethod.Method() != shared.NoAuthRequired {
		client.setState(PendingAuthentication)
	}
	client.setState(Authenticated)
	return nil
}
func (client *Socks5Client) handleCommandResponse() (string, uint16, error) {
	if client.State() != CommandRequested {
		return "", 0, errors.New("client is has not requested command")
	}
	commandResponse, err := waitForServerCommandResponse(client.tcpConn)
	if err != nil {
		client.setError(err)
		return "", 0, err
	}
	if err := isCommandSuccessful(commandResponse); err != nil {
		client.setError(err)
		return "", 0, err
	}

	client.setState(CommandAccepted)
	return commandResponse.BND_ADDR.Value, commandResponse.BND_PORT, nil
}

// utils

func isCommandSuccessful(cmd *command_response.CommandResponse) error {
	if cmd.Status != command_response.Success {
		errMsg := fmt.Sprintf("server didn't respond with success, responed with %v", cmd.Status)
		return errors.New(errMsg)
	}
	return nil
}

func waitForServerCommandResponse(client net.Conn) (*command_response.CommandResponse, error) {
	buf := make([]byte, 64)
	_, err := client.Read(buf)
	if err != nil {
		return nil, err
	}

	commandResponse := command_response.CommandResponse{}
	err = commandResponse.Deserialize(buf)
	if err != nil {
		return nil, err
	}
	return &commandResponse, nil
}
func openTcpConnection(servAddr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
