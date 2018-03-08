package main

import (
	"github.com/drornir/goodis/server"
	"log"
	"github.com/drornir/goodis/handling"
)

func main(){
	srv := NewServer(6390)
	e := srv.Listen()
	if e != nil {
		log.Fatal(e)
	}
	srv.Close()
}



func NewServer(port int) Server {
	return server.New(port, handling.New())
}
