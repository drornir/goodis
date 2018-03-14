package main

import "fmt"

type ProtocolError struct{ Got, Expected string }

func (e ProtocolError) Error() string {
	return fmt.Sprintf("protocol error: expected '%v', got '%v' ", e.Expected, e.Got)
}
