package logger

import (
	"fmt"
)

func Infof(m string, p ...string) {
	fmt.Println(fmt.Sprintf("INFO: %s", m))
}

func Errorf(err error, m string, p ...string) {
	m = fmt.Sprintf(m, p)
	fmt.Println(fmt.Errorf("ERROR: %s: %w", m, err))
}
