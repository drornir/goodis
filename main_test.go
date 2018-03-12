package main

import (
	"github.com/drornir/goodis/testutils"
	"log"
	"testing"
)

func TestPing(t *testing.T) {
	s, err := NewServer(testutils.ServerAddr)
	if err != nil {
		t.Error(err)
	}

	client := testutils.NewRedisClient(testutils.ServerAddr)
	defer client.Close()

	expect := "PONG"

	resp, err := client.Ping().Result()
	if err != nil {
		t.Error(err)
	}

	if expect != resp {
		t.Fatal(testutils.Expected(expect, resp))
	}

	closeErr := s.Close()

	log.Printf("server exit with '%v'\n", closeErr)
}
