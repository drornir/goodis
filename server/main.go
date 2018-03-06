package server

import (
	"net"
	"strconv"
	"fmt"
)

type Server struct {
	host     string
	port     string
	listener net.Listener
}

func New(port int) *Server {
	s := new(Server)
	s.host = "localhost"
	s.port = strconv.Itoa(port)
	return s
}

func (s *Server) Listen() error {

	addr := fmt.Sprintf("%v:%v", s.host, s.port)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("tcp listen failed for address `%v`:\n%v", addr, err)
	}

	s.listener = l
	return nil
}

func (s *Server) Close() {
 	_ = s.listener.Close()
 	s.listener = nil
}
