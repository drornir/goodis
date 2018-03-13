package testutils

import (
	"fmt"
	"testing"
)

func Expected(t *testing.T, exp, got string) string {
	return fmt.Sprintf(
		"unexpected reponse in %v\n"+
			"\texpected: '%v'\n"+
			"\tgot:      '%v'\n",
		t.Name(), exp, got)
}

func Fatal(t *testing.T, err error) string {
	return fmt.Sprintf("fatal error in %v: %v", t.Name(), err)
}
