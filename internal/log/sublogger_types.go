package log

import (
	"io"
	"sync"
)

// Global vars related to the logger package
var (
	SubLoggers = map[string]*SubLogger{}

	Global            *SubLogger
	SubsystemLogger   *SubLogger
	ConnMonitor       *SubLogger
	DatabaseLogger    *SubLogger
	NTPLogger         *SubLogger
	TraderLogger      *SubLogger
	Webserver         *SubLogger
	SentimentAnalyzer *SubLogger
	ExchangeManager   *SubLogger
)

// SubLogger defines a sub logger can be used externally for packages wanted to
// leverage GCT library logger features.
type SubLogger struct {
	name   string
	levels Levels
	output io.Writer
	mtx    sync.RWMutex
}

// logFields is used to store data in a non-global and thread-safe manner
// so logs cannot be modified mid-log causing a data-race issue
type logFields struct {
	info   bool
	warn   bool
	debug  bool
	error  bool
	name   string
	output io.Writer
	logger Logger
}
