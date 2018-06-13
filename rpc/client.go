package rpc

import (
	"github.com/gjvnq/go-logger"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

var clientLog, _ = logger.New("rpc/client", 1, os.Stdout)

type Client interface {
	Connect(addr string) error
	Do(payload string) (string, error)
}

type client struct {
	seq       uint32
	Con       net.Conn
	Transport Transport
	Protocol  Protocol

	resolves *sync.Map
	ch       chan *Message
}

func (c *client) Connect(addr string) error {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		clientLog.Error(err.Error())
		return err
	}
	c.Con = con
	c.Transport = NewTransport(c.Con)
	c.Protocol = NewProtocol(c.Transport)

	go c.receive()

	return nil
}

func (c *client) newSeq() uint32 {
	return atomic.AddUint32(&c.seq, 1)
}

func (c *client) Do(payload string) (string, error) {
	seq := c.newSeq()

	message := Message{
		Seq:     seq,
		Type:    REQUEST,
		Payload: payload,
	}
	req := &Request{
		Message: message,
		resp:    make(chan *Message),
	}

	clientLog.DebugF("send message %#v", message)
	c.Protocol.WriteMessage(&message)

	c.resolves.Store(seq, req)

	for {
		select {
		case resp := <-req.resp:
			log.DebugF("receive from channel %#v", resp)
			c.resolves.Delete(seq)
			return resp.Payload, nil
		}
	}
}

func (c *client) receive() {
	for {
		resp, err := c.Protocol.ReadMessage()

		if err != nil {
			switch err {
			case io.EOF:
				continue
			default:
				log.Error(err.Error())
				return
			}
		}
		log.DebugF("got resp : %#v", resp)

		if req, ok := c.resolves.Load(resp.Seq); ok {
			req.(*Request).resp <- resp
		}
	}
}

func NewClient() Client {
	return &client{
		resolves: &sync.Map{},
		ch:       make(chan *Message),
	}
}
