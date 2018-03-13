package main

import (
	"github.com/drornir/goodis/goodis"
	"github.com/drornir/goodis/tcpserver"
	"io"
	"log"
)

func main() {
	server, err := NewServer(":6379")
	if err != nil {
		log.Fatalf("error creating server: %v ", err)
	}
	defer server.Close()
	for {
		continue
	}
}

func NewServer(addr string) (io.Closer, error) {
	return tcpserver.NewListeningTCPServer(addr, AppHandler())
}

func AppHandler() tcpserver.Handler {
	return tcpserver.Handler{Handle: goodis.New().Handle}
}
