package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Handler struct {
	client net.Conn
	buffer []byte
}

func NewHandler(client net.Conn) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) Handle() {
	defer h.client.Close()

	clientPrependedWithPrevious := io.MultiReader(bytes.NewReader(h.buffer), h.client)
	reader := bufio.NewReader(clientPrependedWithPrevious)

	cmds, unprocessed, err := h.unmarshal(reader)
	if err == nil || err == io.EOF {
		h.buffer = []byte(unprocessed)
	} else {
		log.Fatalln(err)
	}

	var resp string
	for _, cmd := range cmds {
		c := strings.ToLower(string(cmd[0]))
		switch c {
		case "ping":
			resp += "+PONG\r\n"
		case "echo":
			resp += fmt.Sprintf("$%v\r\n%v\r\n", len(cmd[1]), string(cmd[1]))
		}
	}

	h.client.Write([]byte(resp))
}
func (h *Handler) unmarshal(reader *bufio.Reader) (cmds [][]string, remainder string, e error) {
	for {
		header, err := reader.ReadString('\n')
		switch err {
		case io.EOF:
			return cmds, header, io.EOF
		case nil:
		default:
			return cmds, "", err
		}

		if header[0] != '*' {
			return cmds, "", &ProtocolError{Expected: "*", Got: string(header[0])}
		}

		content := string(header[1:len(header)-2])
		numArgs, err := strconv.Atoi(content)
		if err != nil {
			return cmds, "", &ProtocolError{Expected: "<number>", Got: content}
		}

		cmd, rest, err := unmarshalOne(reader, numArgs)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		switch err {
		case io.EOF:
			return cmds, rest, io.EOF
		case nil:
		default:
			return cmds, "", err
		}
	}
}

func unmarshalOne(reader *bufio.Reader, numArgs int) ([]string, string, error) {
	var cmd []string
	for i := 0; i < numArgs; i++ {
		firstRune, _, err := reader.ReadRune()
		switch err {
		case io.EOF:
			return cmd, string(firstRune), io.EOF
		case nil:
		default:
			return cmd, "", err
		}

		if firstRune != '$' {
			return cmd, "", &ProtocolError{Expected: "$", Got: string(firstRune)}
		}

		line, err := reader.ReadString('\n')
		switch err {
		case io.EOF:
			return cmd, line, io.EOF
		case nil:
		default:
			return cmd, "", err
		}

		content := string(line[:len(line)-2])
		length, err := strconv.Atoi(content)
		if err != nil {
			return cmd, "", err
		}

		b := make([]byte, 0, length+2)
		n, err := reader.Read(b)
		switch err {
		case io.EOF:
			return cmd, string(b[:n]), io.EOF
		case nil:
		default:
			return cmd, "", err
		}
		if n != length+2 {
			return cmd, string(b[:n]), nil
		}

		cmd = append(cmd, string(b[:length]))
	}
	return cmd, "", nil
}
