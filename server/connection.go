package server

import (
	"fmt"
	"io"
	"net"
)

type Connection struct {
	conn       *net.Conn
	bufferSize uint
}

func NewConnection(conn *net.Conn) Connection {
	return Connection{conn, 1020}
}

func (c Connection) Read() (string, error) {
	connection, buffer := c.prepare()
	bytesRead, err := connection.Read(buffer)
	switch err {
	case io.EOF, nil:
		return string(buffer[:bytesRead]), err
	default:
		return "", fmt.Errorf("read failed: read %v bytes and got %v: ", bytesRead, err)
	}
}

func (c Connection) Send(string) error {
	bf := c.newBuffer()
	cnn := *c.conn
	if nn, err := cnn.Write(bf); err != nil {
		return fmt.Errorf("write failed: wrote %v characters and got %v: ", nn, err)
	}
	return cnn.Close()
}

func (c Connection) newBuffer() []byte {
	return make([]byte, c.bufferSize)
}

func (c Connection) prepare() (cnn net.Conn, bf []byte) {
	return *c.conn, c.newBuffer()
}
