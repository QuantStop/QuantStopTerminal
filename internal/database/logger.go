package database

import (
	"github.com/quantstop/quantstopterminal/internal/log"
)

// Logger implements io.Writer interface to redirect SQL debug output to logger
type Logger struct{}

// Write takes input and sends to logger
func (l Logger) Write(p []byte) (n int, err error) {
	log.Debugf(log.Database, "SQL: %s", p)
	return 0, nil
}
