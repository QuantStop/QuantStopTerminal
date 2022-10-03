package coinbasepro

import "fmt"

// Capture ensures errors from deferred funcs are captured when an error has not already been set.
func Capture(capture *error, deferred error) {
	if *capture == nil {
		*capture = deferred
	}
}

type Error struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.StatusCode)
}
