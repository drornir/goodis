package tcpserver

import (
	"fmt"
	"io"
	"log"
	"net"
)

type TCPServer interface {
	Close() error
}

type tcp struct {
	handler  Handler
	listener net.Listener
}

type Handler struct {
	Handle func(string) (string, error)
}

func (t *tcp) Close() error {
	return t.listener.Close()
}

func NewListeningTCPServer(addr string, h Handler) (TCPServer, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcp listen on '%v' failed: %v", addr, err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Printf("server tcp accept error: %v ", err)
				break
			}

			go handleConn(conn, h)
		}
	}()

	return &tcp{handler: h, listener: l}, nil
}

func handleConn(conn net.Conn, h Handler) {
	defer conn.Close()

	bytesRead, readErr := read(conn)
	if readErr != nil {
		log.Printf("read from '%v' finished: %v", conn, readErr)
		return
	}

	resp, _ := h.Handle(string(bytesRead))

	bf := []byte(resp)
	_, writeErr := conn.Write(bf)
	if writeErr != nil {
		log.Printf("write '%v' to connection '%v': %v ", resp, conn.RemoteAddr(), writeErr)
		return
	}
}

func read(conn net.Conn) ([]byte, error) {
	bf := make([]byte, 1024)
	n, readErr := conn.Read(bf)

	switch readErr {
	case nil:
		return bf[:n], nil
	case io.EOF:
		if n > 0 {
			return bf[:n], nil
		} else {
			return []byte{}, io.EOF
		}
	default:
		return []byte{}, fmt.Errorf("read from conn '%v' failed : %v", conn, readErr)
	}
}
