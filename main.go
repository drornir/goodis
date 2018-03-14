package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	srv := New("localhost:6379")
	err := srv.Listen()
	fmt.Printf("server exited with %v", err)
}

type Server struct {
	addr     string
	listener net.Listener
}

func New(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = listener

	acceptedConnections := make(chan net.Conn)
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				acceptedConnections <- conn
			} else {
				log.Printf("server tcp accept error: %v ", err)
				close(acceptedConnections)
				return
			}
		}
	}()

	clients := map[net.Conn]*Handler{}
	for {
		select {
		case client := <-acceptedConnections:
			clients[client] = NewHandler(client)
			go clients[client].Handle()
		default:
			//
		}
	}

	return nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

//func (s *Server) handleConn(client net.Conn) {
//	defer client.Close()
//	reader := bufio.NewReader(client)
//
//	header, err := reader.ReadString('\n')
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	if header[0] != '*' {
//		return
//	}
//	header = header[1:len(header)-2]
//	numArgs, _ := strconv.Atoi(header)
//
//	var cmd [][]byte
//	for i := 0; i < numArgs; i++ {
//		line, _ := reader.ReadString('\n')
//		line = line[1:len(line)-2]
//		length, _ := strconv.Atoi(line)
//
//		b := make([]byte, length+2)
//		n, _ := io.ReadFull(reader, b)
//		if n != length+2 {
//			return
//		}
//		cmd = append(cmd, b[:length])
//	}
//
//	var resp string
//	c := strings.ToLower(string(cmd[0]))
//	switch c {
//	case "ping":
//		resp = "+PONG\r\n"
//	case "echo":
//		resp = fmt.Sprintf("$%v\r\n%v\r\n", len(cmd[1]), string(cmd[1]))
//	}
//
//	client.Write([]byte(resp))
//}
