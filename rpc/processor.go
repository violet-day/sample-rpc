package rpc

import (
	"fmt"
	"strconv"
)

type Processor interface {
	Process(message *Message)
}

type processor struct {
	Transport Transport
	Protocol  Protocol
}

func (p *processor) Process(message *Message) {

	v, _ := strconv.Atoi(message.Payload)
	payload := fmt.Sprintf("%d->%d", v, v*2)

	resp := &Message{
		Seq:     message.Seq,
		Type:    RESPONSE,
		Payload: payload,
	}

	p.Protocol.WriteMessage(resp)

	log.DebugF("send resp: %#v", resp)
}

func NewProcessor(t Transport, p Protocol) Processor {
	return &processor{
		Transport: t,
		Protocol:  p,
	}
}
