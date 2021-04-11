package logger

import (
	"fmt"
)

func Infof(m string, p ...string) {
	fmt.Printf("INFO: %s\n", fmt.Sprintf(m, p))
}

func Errorf(err error, m string, p ...string) {
	m = fmt.Sprintf(m, p)
	fmt.Println(fmt.Errorf("ERROR: %s: %w", m, err))
}
