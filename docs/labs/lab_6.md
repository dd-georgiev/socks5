- [Lab 6 - Implementing SOCKS5 Server CONNECT Command](#lab-6---implementing-socks5-server-connect-command)
    * [Overview](#overview)
    * [Experiments](#experiments)
        + [Setting up listener and session handler](#setting-up-listener-and-session-handler)
        + [Handling authentication](#handling-authentication)
        + [Handling CONNECT command](#handling-connect-command)
        + [Simple TCP Proxy](#simple-tcp-proxy)
        + [Modifying the client test to work with our server](#modifying-the-client-test-to-work-with-our-server)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

# Lab 6 - Implementing SOCKS5 Server CONNECT Command
## Overview
The goal is to use the test written in Lab 3, but instant of connecting to Dante to connect to a server written in Golang
## Experiments
### Setting up listener and session handler
```go
type SessionState uint16
// The states are with values 10, 20 30 so that in the future more states can be added.  For example between PendingAuthMethods Authenticated there may be a need to have more states
// NOTE: There the "ERRORED" state is not presented. In real implementation it should be!
const PendingAuthMethods SessionState = 10
const Authenticated SessionState = 20
const Proxying SessionState = 30
// I call a client connected to the proxy "session", as it passes though multiple states and its stateful(that is if a client reconnects the whole process must start again)
type Session struct {
	state SessionState
	conn  net.Conn
	err   error
}
// Basically the example from https://pkg.go.dev/net#Listener
func Start(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		session := Session{state: PendingAuthMethods, conn: conn}
		go session.handler()
	}
}
// This function invokes a specific handler depending on the state of the client 
func (session *Session) handler() {
	for {
		switch session.state {
            case PendingAuthMethods:
                session.handleAuth()
            case Authenticated:
                session.handleCmd()
            case Proxying:
                return
		}
	}
}
```
### Handling authentication
```go
func (session *Session) handleAuth() {
	// read the message send by the client in a buffer, as the message may be invalid this is still a message candidate and not a message.
	authMethodCandidate := make([]byte, 1024)
	n, err := session.conn.Read(authMethodCandidate)
	if err != nil {
		log.Fatal(err)
	}
	// extract the auth methods which the client supports
	authMethods := available_auth_methods.AvailableAuthMethods{}
	err = authMethods.Deserialize(authMethodCandidate[:n])
	if err != nil {
		session.setError(err)
	}
	msg := accept_auth_method.AcceptAuthMethod{} // we will return this to the client (potentially, if NoAuth is in the client methods)
	// make sure that the client is happy with NoAuth method. More complex auth logic must go here
	if slices.Contains(authMethods.Methods(), shared.NoAuthRequired) { 
		err := msg.SetMethod(shared.NoAuthRequired)
		if err != nil {
			session.setError(err)
		}
	}
	// the client is not accepting the NoAuth method and the server cannot process the client request. As such we will return NoAcceptableMethods to it(the client)
	err = msg.SetMethod(shared.NoAcceptableMethods) 
	if err != nil {
		session.setError(err)
	}
	session.conn.Write(msg.ToBytes())
	// !!!!!
	session.state = Authenticated // we set the session to authenticated. In real server we probably shouldn't.
	// !!!!!
}
```
### Handling CONNECT command
```go
func (session *Session) handleCmd() {
	commandCandidate := make([]byte, 1024)
	n, err := session.conn.Read(commandCandidate)
	if err != nil {
               session.setError(err)
        }

	cmd := command_request.CommandRequest{}
	err = cmd.Deserialize(commandCandidate[:n])
	if err != nil {
		session.setError(err)
	}

	switch cmd.CMD {
	case shared.CONNECT:
		session.handleConnectCmd(cmd)
		return
	}
// the command is not supported, so we return such message
	resp := command_response.CommandResponse{}
	resp.Status = command_response.CommandNotSupported
	bytes, err := resp.ToBytes()
	if err != nil {
		session.setError(err)
	}
	session.conn.Write(bytes)
	session.setError(errors.New("unsupported command"))
	return

}

func (session *Session) handleConnectCmd(cmd command_request.CommandRequest) {
	proxyErrors := make(chan error) // The code counts on centralized error handler(see below).
	remoteAddr := fmt.Sprintf("%s:%d", cmd.DST_ADDR.Value, cmd.DST_PORT)
	proxy, err := proxies.NewConnectProxy(remoteAddr, session.conn)
	if err != nil {
		session.setError(err)
	}
	// The proxy is described in the next section
	go proxy.Start(proxyErrors) 
	if err != nil {
		session.setError(err)
	}
	// start the error handler
	go session.proxyErrorHandler(proxyErrors, proxy)
    // Notify the client that the proxy is running and set the state to proxying.
	resp := command_response.CommandResponse{}
	resp.Status = command_response.Success
	resp.BND_PORT = 0
	resp.BND_ADDR = shared.DstAddr{Value: "0.0.0.0", Type: shared.ATYP_IPV4}
	bytes, _ := resp.ToBytes()
	session.conn.Write(bytes)
	session.state = Proxying
}
// Centralized error handler, if anything anywhere goes wrong it will be sent and handled here.
func (session *Session) proxyErrorHandler(errors chan error, proxy proxies.Proxy) {
	err := <-errors
	session.setError(err)
	proxy.Stop()
}
```
### Simple TCP Proxy
NOTE: The following video goes indepth about the code below. [https://www.youtube.com/watch?v=J4J-A9tcjcA](https://www.youtube.com/watch?v=J4J-A9tcjcA)
#### Splicing two io.ReadWriter
Splicing means joining two things together. In this context, we basically start sending the data from one `ReadWriter` to the other and vice-versa
```go
func SpliceConnections(client io.ReadWriter, server io.ReadWriter, errors chan error) {
	go func() {
		_, err := io.Copy(server, cli`ent)
		if err != nil {
			errors <- err
		}
	}()
	_, err := io.Copy(client, server)
	if err != nil {
		errors <- err
	}
}
```
#### TCP proxy
```go
type TCPProxy struct {
	server io.ReadWriteCloser
	client io.ReadWriteCloser
}

func NewConnectProxy(addr string, client io.ReadWriteCloser) (*TCPProxy, error) {
	server, err := net.DialTimeout("tcp", addr, time.Duration(1)*time.Second)
	if err != nil {
		return nil, err
	}

	return &TCPProxy{server: server, client: client}, nil
}

func (proxy *TCPProxy) Start(errors chan error) error {
	SpliceConnections(proxy.server, proxy.client, errors)
	return nil
}

func (proxy *TCPProxy) Stop() {
	proxy.server.Close()
	proxy.client.Close()
}
```

### Modifying the client test to work with our server
```go
func Test_Client_Connect(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*50)
	proxyAddr, proxyPort := startSocks5Server()
	socks5SrvAddr := fmt.Sprintf("%s:%d", proxyAddr, proxyPort)
// ..... Not different than the one from lab 3
}

func startSocks5Server() (string, int) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go Start(listener)
	addr := listener.Addr().(*net.TCPAddr).IP.String()
	port := listener.Addr().(*net.TCPAddr).Port
	return addr, port
}

```