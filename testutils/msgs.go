package testutils

import "fmt"

func Expected(exp, got string) string {
	return fmt.Sprintf("expected: %v\ngot: %v\n", exp, got)
}
