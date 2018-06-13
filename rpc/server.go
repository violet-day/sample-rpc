package rpc

import (
	"github.com/gjvnq/go-logger"
	"io"
	"net"
	"os"
)

var log, _ = logger.New("rpc/server", 1, os.Stdout)

type Server struct {
	Listener *net.Listener
}

func (s *Server) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.DebugF("server listen %s\n", l.Addr().String())

	s.Listener = &l

	for {
		con, err := l.Accept()
		if err != nil {
			log.Error(err.Error())
		}
		log.DebugF("accept connection from %s", con.RemoteAddr().String())
		go handler(con)
	}
}

func handler(c net.Conn) {
	t := NewTransport(c)
	p := NewProtocol(t)

	processor := NewProcessor(t, p)

	for {
		message, err := p.ReadMessage()

		log.DebugF("receive message %#v from %s", message, c.RemoteAddr())

		if err != nil {
			switch err {
			case io.EOF:
				log.DebugF("remote connection %s closed", c.RemoteAddr())
			default:
				log.ErrorF("read message got error: %s", err.Error())
			}
			c.Close()
			return
		}

		go processor.Process(message)
	}
}
