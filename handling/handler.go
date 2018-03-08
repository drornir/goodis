package handling

import (
	"fmt"
)

type Handler interface {
	Handle(cmd string) (res string, e error)
}

type handler struct{}

func New() Handler {
	return new(handler)
}

func (h *handler) Handle(cmd string) (res string, e *Error) {
	_ = cmd
	res = fmt.Sprintf("+PONG\r\n")
	return res, nil
}
