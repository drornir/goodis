package main

import (
	"fmt"
	"github.com/drornir/goodis/testutils"
	"io"
	"net"
	"testing"
	"time"
)

func TestEcho(t *testing.T) {
	WithServer(t, func() {
		clt := testutils.NewRedisClient()
		defer clt.Close()
		expect := "echo\nme\nback"

		resp, err := clt.Echo(expect).Result()
		if err != nil {
			t.Fatalf("client error: %v", err)
		}

		if expect != resp {
			t.Fatal(testutils.Expected(t, expect, resp))
		}
	})
}

func TestPing(t *testing.T) {
	WithServer(t, func() {
		clt := testutils.NewRedisClient()
		defer clt.Close()

		expect := "PONG"

		resp, err := clt.Ping().Result()
		if err != nil && err != io.EOF {
			t.Fatalf("client error: %v", err)
		}

		if expect != resp {
			t.Fatal(testutils.Expected(t, expect, resp))
		}
	})
}

func WithServer(t *testing.T, cb func()) {
	srv := New(testutils.ServerAddr)
	defer srv.Close()
	go func() {
		_ = srv.Listen()
	}()

	if err := waitUntilOpen(testutils.ServerAddr); err != nil {
		t.Fatalf("'%v' is not open: %v", testutils.ServerAddr, err)
	}

	cb()
}

func waitUntilOpen(addr string) error {
	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timeout")
		default:
			if isOpen(addr) {
				return nil
			} else {
				time.Sleep(time.Millisecond)
			}
		}
	}
}
func isOpen(addr string) bool {
	l, err := net.Dial("tcp", addr)
	if l != nil {
		defer l.Close()
	}

	if err != nil {
		return false
	} else {
		return true
	}
}
