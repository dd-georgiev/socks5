package server

import (
	"log"
	"net"
	"time"
)

type SessionState uint16

const PendingAuthMethods SessionState = 10
const Authenticated SessionState = 20
const Proxying SessionState = 30

type Session struct {
	state SessionState
	conn  net.Conn
	err   error
}

func (session *Session) setError(err error) {
	session.RespondToClientDependingOnState()
	session.err = err
	go session.closeClientAfter5Seconds()
}
func Start(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		session := Session{state: PendingAuthMethods, conn: conn}
		go session.handler()
	}
}

func (session *Session) handler() {
	for {
		switch session.state {
		case PendingAuthMethods:
			session.handleAuth()
		case Authenticated:
			session.handleCommand()
		case Proxying:
			return
		}
	}
}

func (session *Session) closeClientAfter5Seconds() {
	time.Sleep(5 * time.Second)
	session.conn.Close()
}
