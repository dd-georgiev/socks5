# Overview
The goal of this exercise is to explore the protocol by monitoring the network traffic from client, through proxy server and to mocked TCP/UDP server. In the end we will have a series of pcaps outlining the different SOCKS5 commands. Authentication and traffic rules(i.e. the proxy server allowing/denying certain requests is out of scope.
# Tooling
1. [[https://github.com/v-byte-cpu/wirez][Wirez]] is a tool for redirecting TCP and UDP traffic to SOCKS5 proxy. It is written in Golang. It supports the `CONNECT` and `UDP ASSOCIATE` commands
2. [[https://linux.die.net/man/1/socat][Socat]] is a tool for relaying information between two bidirectional byte streams. It supports the `BIND` command
3. [[https://www.wireshark.org/][Wireshark]] is tool for capturing and analyzing network traffic.
4. [[https://www.inet.no/dante/][Dante]] is a free SOCKS5 server
5. [[https://linux.die.net/man/1/nc][netcat]] is tool for working with TCP and UDP connections.
# CONCEPTS
## SOCKS5 commands
In [[https://datatracker.ietf.org/doc/html/rfc1928][RFC1928]], the following commands are outlined:
- CONNECT - This is the basic forwarding of TCP segments from the client to the server.
- BIND - The BIND command is used in scenarios, where the server will attempt to connect back to the client. A typical usage of BIND is for P2P network protocols or FTP([example case](https://stackoverflow.com/questions/25092819/when-should-an-ftp-server-connect-to-ftp-client-after-port-command).
- UDP_ASSOCIATE - The proxy is relaying UDP datagrams to server and the responses from the server are relayed back to the client. The connection is started via TCP, but later on the proxy server must offer an UDP listener, on which datagrams with specific format are being sent.
## SOCKS5 connection flow
In [[https://datatracker.ietf.org/doc/html/rfc1928][RFC1928]], the connection flow is outlined as follows:
1. The `client` sends a message containing the available authentication methods
2. The `server` picks authentication method, or returns `NO ACCEPTED METHODS`
3. If the response is `NO ACCEPTED METHODS` the `client` terminates the connection
4. If the response contains a desired authentication method, the `client` and the `server` enter a method-dependent sub-negotiation
5. Once the authentication is completed, the `client` sends a command - either `CONNECT`, `BIND` or `UDP_ASSOCIATE`
6. The server evaluates the requests and returns a response, either indicating success or failure.
7. In case of success, the proxying starts.

DIAGRAM_HERE
## SOCKS5 messages formats
The messages of SOCKS5 can be roughly divided in `requests` and `responses`. Every message starts with a field called `VER`, which in SOCKS5 is always the value of 5 (or 0x05). At the end of this lab, we will have an example of each message which the protocol provides. [RFC1928](https://datatracker.ietf.org/doc/html/rfc1928) provides some examples.
# Experiments
## Prerequisites
### Compiling wirez
As [[https://github.com/v-byte-cpu/wirez][Wirez]] doesn't provide ready to use binary, we must compile it. The instructions below are for Linux. You can compile the program by following the instructions in the repository.
1. Install go 1.20
#+BEGIN_SRC bash
$ go install golang.org/dl/go1.20.9@latest
$ go1.20.9 download
#+END_SRC
2. Clone the repository
#+BEGIN_SRC bash
$ git clone https://github.com/v-byte-cpu/wirez.git
#+END_SRC
3. Compile
#+BEGIN_SRC bash
$ cd wirez
$ go build
#+END_SRC
### Starting Dante on local machine with Docker
Running [[https://www.inet.no/dante/][Dante]] locally with Docker is a convenient way, as the setup doesn't focus on the specifics of the proxy server. 
1. [[https://docs.docker.com/engine/install/][Install docker]]
2. Pull the following image: ~wernight/dante~
3. To start the container and bind port ~1080~ to it run ~docker run -p 1080:1080 wernight/dante~
4. If everything was successful so far, you must have the following log line:
#+BEGIN_SRC
Jan 24 19:42:40 (1737747760.895891) sockd[7]: info: Dante/server[1/1] v1.4.2 running
#+END_SRC
### Installing other tools
The rest of the tooling (~netcat~, ~socat~, and ~Wireshark~) is widely available, consult official guidelines for your platform.
## Observing CONNECT command
## Observing BIND command
## Observing UDP_ASSOCIATE command
# Results
## Message references
## Connection flows
### TCP via CONNECT
### TCP via BIND
### UDP via UDP_ASSOCIATE
# Readings
1. https://datatracker.ietf.org/doc/html/rfc1928
2. https://go.dev/doc/manage-install
