package rpc

import (
	"sync"
)

type MessageType uint32

const (
	REQUEST MessageType = 1 << iota
	RESPONSE
)

type Message struct {
	Seq     uint32
	Type    MessageType
	Payload string
}

type Request struct {
	Message
	resp chan *Message
}

type Protocol interface {
	ReadMessage() (*Message, error)
	WriteMessage(message *Message)
}

type protocol struct {
	rl        sync.Mutex
	wl        sync.Mutex
	transport Transport
}

func (p *protocol) ReadMessage() (*Message, error) {
	p.rl.Lock()
	defer p.rl.Unlock()

	seq, err := p.transport.ReadI32()
	if err != nil {
		return nil, err
	}
	messageType, _ := p.transport.ReadI32()

	if err != nil {
		return nil, err
	}

	payload, err := p.transport.ReadString()

	if err != nil {
		return nil, err
	}

	return &Message{
		Seq:     seq,
		Type:    MessageType(messageType),
		Payload: payload,
	}, nil
}

func (p *protocol) WriteMessage(message *Message) {
	p.wl.Lock()
	defer p.wl.Unlock()

	p.transport.WriteI32(message.Seq)
	p.transport.WriteI32(uint32(message.Type))
	p.transport.WriteString(message.Payload)
}

func NewProtocol(transport Transport) Protocol {
	return &protocol{
		transport: transport,
	}
}
