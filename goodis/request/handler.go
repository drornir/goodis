package request

import (
	"fmt"
	"strings"
	"strconv"
	"io"
)

type Handler struct {
	buffer string
}

func (h *Handler) Handle(req string) (string, error) {
	h.buffer += req
	cmds, err := h.unmarshal()

	response := ""
	for _, cmd := range cmds {
		resp, err := handleCmd(cmd)
		if err == nil {
			response += resp
		} else {
			return "", fmt.Errorf("error while handling cmd '%v': %v", cmd, err)
		}
	}

	return response, nil
}

func (h *Handler) unmarshal() ([][]string, error) {
	cmds := make([][]string, 0)

	for {
		header, err := h.consumeLine()
		switch err {
		case io.EOF:
			// ignore
		default:
			return cmds, err
		}

		if header[0] != '*' {
			return cmds, nil
		}

		numArgs, _ := strconv.Atoi(header[1:])

		cmd := make([]string, 0)
		var line string
		for i := 0; i < numArgs; i++ {
			line, lines = lines[0], lines[1:]
			if line[0] != '$' {
				return "", &ProtocolError{"$", string(line[0])}
			}
			length, _ := strconv.Atoi(line[1:])
			line, lines = lines[0], lines[1:]

			value := line[:length]

			cmd = append(cmd, value)
		}

	}
}

func handleCmd(cmd []string) (string, error) {
	switch strings.ToLower(cmd[0]) {
	case "ping":
		return "+PONG\r\n", nil
	case "echo":
		return fmt.Sprintf("$%v\r\n%v\r\n", len(cmd[1]), cmd[1]), nil
	default:
		return fmt.Sprintf("-Error parsing %v", cmd), nil
	}
}

// consumes the buffer
func (h *Handler) consumeLine() (string, error) {
	splitted := strings.SplitN(h.buffer, "\r\n", 1)
	if len(splitted) == 1 {
		return "", io.EOF
	}
	line, remainder := splitted[0], splitted[1]
	h.buffer = remainder
	return line, nil
}
