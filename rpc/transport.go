package rpc

import (
	"encoding/binary"
	"io"
)

const (
	READ_BUFFER_SIZE  = 1024
	WRITE_BUFFER_SIZE = 512
)

type Transport interface {
	Read(len uint32) ([]byte, error)
	Write(bs []byte)

	ReadI32() (uint32, error)
	WriteI32(i32 uint32)

	ReadString() (string, error)
	WriteString(str string)
}

type transport struct {
	RW io.ReadWriter
}

func NewTransport(rw io.ReadWriter) Transport {
	return &transport{
		RW: rw,
	}
}

func (t *transport) Read(len uint32) ([]byte, error) {
	var buf = make([]byte, len)
	_, err := t.RW.Read(buf)
	return buf, err
}

func (t *transport) Write(bs []byte) {
	t.RW.Write(bs)
}

func (t *transport) ReadI32() (uint32, error) {
	buf, err := t.Read(4)
	if err != nil {
		return 0, err
	}
	i32 := binary.LittleEndian.Uint32(buf)
	return i32, nil
}

func (t *transport) WriteI32(i32 uint32) {
	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i32)
	t.RW.Write(buf)
}

func (t *transport) ReadString() (string, error) {
	strLen, err := t.ReadI32()
	if err != nil {
		return "", err
	}
	bs, err := t.Read(strLen)

	if err != nil {
		return "", err
	}
	return string(bs[:]), nil
}

func (t *transport) WriteString(str string) {
	bs := []byte(str)
	t.WriteI32(uint32(len(bs)))
	t.Write(bs)
}
