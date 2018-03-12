package goodis

import "fmt"

type Handler func(string) string

func AppHandler() Handler {
	return func(_ string) string {
		return fmt.Sprintf("+PONG\r\n")
	}
}
