package main

import (
	"fmt"
	"github.com/drornir/goodis/testutils"
	"github.com/go-redis/redis"
	"io"
	"log"
	"testing"
)

func TestEcho(t *testing.T) {
	srv, clt, err := before(t)
	if err != nil {
		t.Fatal(testutils.Fatal(t, err))
		return
	}
	defer closeAndLog(srv)
	defer clt.Close()

	expect := fmt.Sprintf("hi\nman")

	resp, err := clt.Echo(expect).Result()
	if err != nil {
		t.Error(err)
	}

	if resp != expect {
		t.Fatal(testutils.Expected(t, expect, resp))
	}

	closeAndLog(srv)
}

func TestPing(t *testing.T) {
	srv, clt, err := before(t)
	if err != nil {
		t.Fatal(testutils.Fatal(t, err))
		return
	}
	defer closeAndLog(srv)
	defer clt.Close()

	expect := "PONG"

	resp, err := clt.Ping().Result()
	if err != nil {
		t.Error(err)
	}

	if expect != resp {
		t.Fatal(testutils.Expected(t, expect, resp))
	}
}
func before(t *testing.T) (io.Closer, *redis.Client, error) {
	srv, err := NewServer(testutils.ServerAddr)
	if err != nil {
		return nil, nil, err
	}
	clt := testutils.NewRedisClient(testutils.ServerAddr)
	return srv, clt, nil
}

func closeAndLog(s io.Closer) {
	closeErr := s.Close()
	log.Printf("server exit with '%v'\n", closeErr)
}
