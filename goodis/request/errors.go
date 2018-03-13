package request

import "fmt"

type ProtocolError struct {
	Exp, Got string
}

func (err *ProtocolError) Error() string {
	return fmt.Sprintf("protocol error: expected '%v', got '%v'", err.Exp, err.Got)
}
