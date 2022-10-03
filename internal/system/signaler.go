package system

import (
	"os"
	"os/signal"
	"syscall"
)

// SIGTERM gets sent as the generic software termination signal.
// This is sent for almost all shutdown events (except for those below).

// SIGKILL gets sent as a termination signal that is sent for “quit immediately” events.
// This generally should not be interfered with.

// SIGINT gets sent when user inputs an interrupt signal (such as Ctrl+C).
// This is similar to SIGTERM, but for user events

// SIGQUIT gets sent when user inputs a quit signal (Ctrl+D).
// This is similar to SIGKILL, but for user events (such as force quit)

var (
	shutdownChannel = make(chan os.Signal, 1)
)

func init() {
	signals := []os.Signal{
		os.Interrupt,
		os.Kill,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	}
	signal.Notify(shutdownChannel, signals...)
}

// WaitForInterrupt waits until an os.Signal is
// received and returns the result
func WaitForInterrupt() os.Signal {
	return <-shutdownChannel
}
