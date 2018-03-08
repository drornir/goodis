package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"log"
	"github.com/drornir/goodis/handling"
)

type Server interface {
	Listen() error
	Close()
}

type server struct {
	host     string
	port     string
	listener net.Listener
	handler  handling.Handler
}

func New(port int, handler handling.Handler) *server {
	s := new(server)
	s.host = ""
	s.port = strconv.Itoa(port)
	s.handler = handler
	registerHttpHandlers(s)

	return s
}

func (s *server) Listen() (err error) {
	if listener, err := net.Listen("tcp", s.addr()); err != nil {
		return fmt.Errorf("listen on %v failed: %v", s.addr(), err)
	} else {
		s.listener = listener
	}
	log.Printf("listening on '%v'\n", s.addr())
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			err = fmt.Errorf("can't accept connection: %v", err)
			break
		}
		connection := NewConnection(&conn)

		cmd, err := connection.Read()
		switch err {
		case io.EOF, nil:
		default:
			err = fmt.Errorf("can't read: %v", err)
			break
		}

		resp, err := s.handler.Handle(cmd)
		if err != nil {
			err = fmt.Errorf("can't handle command '%s': %v", cmd, err)
			break
		}
		connection.Send(resp)
	}
	if err != nil {
		return fmt.Errorf("listen stopped with error: %v", err)
	} else {
		return nil
	}

}
func (s *server) Close() {
	_ = s.listener.Close()
	s.listener = nil
}
func registerHttpHandlers(s *server) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response, err := s.handler.Handle(r.URL.Path[1:])
		if err != nil {
			http.Error(w, fmt.Sprintf("request %v: %v", r, err), err.Status())
		}
		fmt.Fprint(w, response)
	})
}

func (s *server) addr() string {
	return fmt.Sprintf("%v:%v", s.host, s.port)
}
