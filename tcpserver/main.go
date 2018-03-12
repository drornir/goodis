package tcpserver

import (
	"fmt"
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
	Handle func(string) string
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
	bf := make([]byte, 1024)
	_, readErr := conn.Read(bf)
	if readErr != nil {
		log.Printf("read from connection '%v': %v ", conn.RemoteAddr(), readErr)
		bf = []byte{}
	}
	read := string(bf)

	resp := h.Handle(read)

	bf = []byte(resp)
	_, writeErr := conn.Write(bf)
	if writeErr != nil {
		log.Printf("write '%v' to connection '%v': %v ", resp, conn.RemoteAddr(), writeErr)
		return
	}
}
