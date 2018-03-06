package goodis

import (
	"fmt"
	"github.com/drornir/goodis/server"
	"github.com/go-redis/redis"
	"net"
	"testing"
	"time"
)

const testPort = 6380

var s *server.Server

func TestRespondToPing(t *testing.T) {
	serverUP(testPort)
	defer serverDown()

	err := waitForOpenPort()
	if err != nil {
		t.Fatalf("server port is not open: %v", err)
	}

	got, err := client().Ping().Result()
	if err != nil {
		t.Fatalf("pinging failed: %v", err)
	}

	expected := "PONG"
	if got != expected {
		t.Errorf("%v\nexpected: %v\ngot: %v", t.Name(), expected, got)
	}
}

func client() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%v", testPort),
	})
	return client
}

func serverUP(port int) {
	go func() {
		s := server.New(port)
		s.Listen()
	}()
}

func serverDown() {
	if s != nil {
		s.Close()
		s = nil
	}
}

func waitForOpenPort() error {
	var backoff time.Duration = 1

	timer := time.NewTimer(1 * time.Second)
	giveUp := false
	go func() {
		<-timer.C
		giveUp = true
	}()

	for !giveUp {
		addr := fmt.Sprintf("localhost:%v", testPort)
		conn, err := net.DialTimeout("tcp", addr, backoff*time.Millisecond)

		if err == nil {
			conn.Close()
			return nil
		} else {
			backoff *= 2
		}
	}
	return fmt.Errorf("wait for open port: timeout error")
}
