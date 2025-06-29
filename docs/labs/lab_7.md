- [Lab 7 - Implementing SOCKS5 Server BIND Command](#lab-7---implementing-socks5-server-bind-command)
    * [Overview](#overview)
    * [Experiments](#experiments)
        + [Handling BIND command](#handling-bind-command)
            - [Modifying the handleCmd method](#modifying-the-handlecmd-method)
            - [Creating handler for the BIND command](#creating-handler-for-the-bind-command)
        + [Creating BIND proxy](#creating-bind-proxy)
        + [Modifying the client test](#modifying-the-client-test)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

# Lab 7 - Implementing SOCKS5 Server BIND Command
## Overview
In this lab, we extend the server code to handle the `BIND` SOCKS5 command. 
## Experiments
### Handling BIND command
#### Modifying the handleCmd method
```go
func (session *Session) handleCmd() {
	commandCandidate := make([]byte, 1024)
	n, err := session.conn.Read(commandCandidate)
// ............. no changes until the switch cmd.CMD statement

	switch cmd.CMD {
// ............. Add another case calling a specific handler for the bind command
	case shared.BIND:
		session.handleBindCmd(cmd)
		return
// ............. no changes
}
```
#### Creating handler for the BIND command
```go
// No much difference from the CONNECT command one.
func (session *Session) handleBindCmd(cmd command_request.CommandRequest) {
	proxyErrors := make(chan error)
	remoteAddr := fmt.Sprintf("%s:%d", cmd.DST_ADDR.Value, cmd.DST_PORT)
	proxy, err := proxies.NewBindProxy(session.conn, remoteAddr) // IMPORTANT, the proxy is different from the CONNECT one
	if err != nil {
		session.setError(err)
	}
	go proxy.Start(proxyErrors)
	if err != nil {
		session.setError(err)
	}
	go session.proxyErrorHandler(proxyErrors, proxy)

	resp := command_response.CommandResponse{}
	resp.Status = command_response.Success
	resp.BND_PORT = proxy.ListeningPort
	resp.BND_ADDR = shared.DstAddr{Value: proxy.ListeningIp, Type: shared.ATYP_IPV4}
	bytes, _ := resp.ToBytes()
	session.conn.Write(bytes)
	session.state = Proxying
}
```
### Creating BIND proxy
```go
package proxies

import (
	"io"
	"net"
	"socks5_server/messages/responses/command_response"
	"socks5_server/messages/shared"
)

type BindProxy struct {
	server        net.Listener
	client        io.ReadWriteCloser
	ListeningPort uint16
	ListeningIp   string
}

func NewBindProxy(client io.ReadWriteCloser, addr string) (*BindProxy, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:8877")
	if err != nil {
		return nil, err
	}

	return &BindProxy{client: client, server: listener, ListeningPort: 8877, ListeningIp: "127.0.0.1"}, nil
}
func (proxy *BindProxy) Start(errors chan error) error {
	go func() {
		in, err := proxy.server.Accept()
		if err != nil {
			errors <- err
			return
		}

		err = proxy.notifyClientAboutIncomingConnection(in)
		if err != nil {
			errors <- err
			return
		}

		SpliceConnections(in, proxy.client, errors)
	}()

	return nil
}

func (proxy *BindProxy) Stop() {
	proxy.server.Close()
	proxy.client.Close()
}

// Here we force the proxy to speak the SOCKS5 protocol. Which is not great design. 
// This most likely can be moved to the server and use the function below to notify the server to notify the client... But it gets more complicated. 
func (proxy *BindProxy) notifyClientAboutIncomingConnection(in net.Conn) error {
	addr := in.RemoteAddr().(*net.TCPAddr).IP.String()
	port := in.RemoteAddr().(*net.TCPAddr).Port
	reqSourceMsg := command_response.CommandResponse{Status: command_response.Success, BND_ADDR: shared.DstAddr{Value: addr, Type: shared.ATYP_IPV4}, BND_PORT: uint16(port)}

	bytes, err := reqSourceMsg.ToBytes()
	_, err = proxy.client.Write(bytes)

	return err
}
```

### Modifying the client test
```go
func Test_Client_Bind(t *testing.T) { 
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
	var serverConnectedWithPort uint16 // this variable is altered by the BindServer to indicate which local port is being used for the server->proxy->client connection
	addr, port := sockstests.BindServer(serverRequestResponse, serverResponse, &serverConnectedWithPort)
    // ................

    // Read actual data
	//// NOTE: Sometimes, the response from the successful connection and the data its self come in the same read operation
	//// Because there is no way to differentiate between them specified in the RFC, the simplest yet reliable way I found is to check if
	//// the message was the whole buffer, if not then at least part of the data arrived with the socks5 message.
	//// this is happening only with the golang server implementation and not with Dante or at least didn't happen during the 100 times I executed this test against Dante
	respLength, _ := msg.ToBytes()
	if n > len(respLength) {
		if string(buf2[len(respLength):n]) != serverResponse {
			t.Fatalf("Expected server to send with TEST, got %v", string(buf2[len(respLength):n]))
		}
		return
	}
	n, err = rwBind.Read(buf2)
	if string(buf2[:n]) != serverResponse {
		t.Fatalf("Expected server to send with TEST, got %v", string(buf2[:n]))
	}
	bindClient.Close()
}
```