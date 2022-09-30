package log

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

var (
	errEmptyLoggerName = errors.New("cannot have empty logger name")

	// ErrSubLoggerAlreadyRegistered Returned when a sub logger is registered multiple times
	ErrSubLoggerAlreadyRegistered = errors.New("sub logger already registered")
)

const (
	timestampFormat = " 02/01/2006 15:04:05 "
	spacer          = " | "
	// DefaultMaxFileSize for logger rotation file
	DefaultMaxFileSize int64 = 100

	defaultCapacityForSliceOfBytes = 100
)

var (

	// logger is the global log instance
	logger = Logger{}

	// FileLoggingConfiguredCorrectly flag set during config check if file logging meets requirements
	FileLoggingConfiguredCorrectly bool

	// GlobalLogConfig holds global configuration options for logger
	GlobalLogConfig = NewConfig()

	// GlobalLogFile hold global configuration options for file logger
	GlobalLogFile = &Rotate{}

	eventPool = &sync.Pool{
		New: func() interface{} {
			sliceOBytes := make([]byte, 0, defaultCapacityForSliceOfBytes)
			return &sliceOBytes
		},
	}

	// FilePath system path to store log files in
	FilePath string

	// RWM read/write mutex for logger
	RWM = &sync.RWMutex{}

	// SubLoggers map of global SubLoggers
	SubLoggers = map[string]*SubLogger{}

	// Global SubLoggers
	Global    *SubLogger
	Database  *SubLogger
	Webserver *SubLogger
)

func init() {

	// register global sub loggers
	Global = registerNewSubLogger("ENGINE")
	Database = registerNewSubLogger("DATABASE")
	Webserver = registerNewSubLogger("WEBSERVER")

	// setup default global log config
	RWM.Lock()
	GlobalLogConfig = NewConfig()
	RWM.Unlock()

	// setup default global logger
	if err := SetupGlobalLogger(); err != nil {
		log.Panicf("Unable to setup default global logger. Error: %s\n", err)
	}

	// setup default global sub loggers
	if err := SetupSubLoggers(GlobalLogConfig.SubLoggers); err != nil {
		log.Panicf("Unable to setup default global sub loggers. Error: %s\n", err)
	}
}

// Logger represents a single Logger instance with settings
type Logger struct {
	ShowLogSystemName                                bool
	Timestamp                                        string
	InfoHeader, ErrorHeader, DebugHeader, WarnHeader string
	Spacer                                           string
}

func newLogger(c *Config) Logger {
	return Logger{
		Timestamp:         c.AdvancedSettings.TimeStampFormat,
		Spacer:            c.AdvancedSettings.Spacer,
		ErrorHeader:       c.AdvancedSettings.Headers.Error,
		InfoHeader:        c.AdvancedSettings.Headers.Info,
		WarnHeader:        c.AdvancedSettings.Headers.Warn,
		DebugHeader:       c.AdvancedSettings.Headers.Debug,
		ShowLogSystemName: *c.AdvancedSettings.ShowLogSystemName,
	}
}

func (log *Logger) newLogEvent(data, header, slName string, w io.Writer) error {
	if w == nil {
		return errors.New("io.Writer not set")
	}

	pool, ok := eventPool.Get().(*[]byte)
	if !ok {
		return errors.New("unable to type assert slice of bytes pointer")
	}

	*pool = append(*pool, header...)
	if log.ShowLogSystemName {
		*pool = append(*pool, log.Spacer...)
		*pool = append(*pool, slName...)
	}
	*pool = append(*pool, log.Spacer...)
	if log.Timestamp != "" {
		*pool = time.Now().AppendFormat(*pool, log.Timestamp)
	}
	*pool = append(*pool, log.Spacer...)
	*pool = append(*pool, data...)
	if data == "" || data[len(data)-1] != '\n' {
		*pool = append(*pool, '\n')
	}
	_, err := w.Write(*pool)
	*pool = (*pool)[:0]
	eventPool.Put(pool)

	return err
}

// CloseLogger is called on shutdown of application
func CloseLogger() error {
	return GlobalLogFile.Close()
}

// Level retries the current sub logger levels
func Level(name string) (Levels, error) {
	RWM.RLock()
	defer RWM.RUnlock()
	subLogger, found := SubLoggers[name]
	if !found {
		return Levels{}, fmt.Errorf("logger %s not found", name)
	}
	return subLogger.levels, nil
}

// SetLevel sets sub logger levels
func SetLevel(s, level string) (Levels, error) {
	RWM.Lock()
	defer RWM.Unlock()
	subLogger, found := SubLoggers[s]
	if !found {
		return Levels{}, fmt.Errorf("sub logger %v not found", s)
	}
	subLogger.SetLevels(splitLevel(level))
	return subLogger.levels, nil
}
