package goodis

import "github.com/drornir/goodis/server"

type Server interface {
	Listen()
	Close()
}

func NewServer(port int) Server {
	return server.New(port)
}
