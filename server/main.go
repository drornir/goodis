package server

type Server struct {
	port int
}

func New(port int) *Server {
	s := new(Server)
	s.port = port
	return s
}

func (s *Server) Listen() {

}

func (s *Server) Close() {

}
