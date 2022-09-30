package system

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	signalChannel = make(chan os.Signal, 1)
)

func init() {
	signals := []os.Signal{
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	}
	signal.Notify(signalChannel, signals...)
}

// WaitForInterrupt waits until an os.Signal is
// received and returns the result
func WaitForInterrupt() os.Signal {
	return <-signalChannel
}
