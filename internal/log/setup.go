package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	errSubLoggerConfigIsNil  = errors.New("subLogger config is nil")
	errUnhandledOutputWriter = errors.New("unhandled output writer")
)

func getWriters(s *SubLoggerConfig) (io.Writer, error) {
	if s == nil {
		return nil, errSubLoggerConfigIsNil
	}

	outputWriters := strings.Split(s.Output, "|")
	writers := make([]io.Writer, 0, len(outputWriters))
	for x := range outputWriters {
		var writer io.Writer
		switch strings.ToLower(outputWriters[x]) {
		case "stdout", "console":
			writer = os.Stdout
		case "stderr":
			writer = os.Stderr
		case "file":
			if FileLoggingConfiguredCorrectly {
				writer = GlobalLogFile
			}
		default:
			// Note: Do not want to add a io.Discard here as this adds
			// additional routines for every write for no reason.
			return nil, fmt.Errorf("%w: %s", errUnhandledOutputWriter, outputWriters[x])
		}
		writers = append(writers, writer)
	}
	return multiWriter(writers...)
}

func configureSubLogger(subLogger, levels string, output io.Writer) error {
	RWM.Lock()
	defer RWM.Unlock()
	logPtr, found := SubLoggers[subLogger]
	if !found {
		return fmt.Errorf("sub logger %v not found", subLogger)
	}

	logPtr.SetOutput(output)
	logPtr.SetLevels(splitLevel(levels))
	SubLoggers[subLogger] = logPtr
	return nil
}

// SetupSubLoggers configure all sub loggers with provided configuration values
func SetupSubLoggers(s []SubLoggerConfig) error {
	for x := range s {
		output, err := getWriters(&s[x])
		if err != nil {
			return err
		}
		err = configureSubLogger(strings.ToUpper(s[x].Name), s[x].Level, output)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetupGlobalLogger set up the global loggers with the default global config values
func SetupGlobalLogger() error {
	RWM.Lock()
	defer RWM.Unlock()

	if FileLoggingConfiguredCorrectly {
		GlobalLogFile = &Rotate{
			FileName: GlobalLogConfig.LoggerFileConfig.FileName,
			MaxSize:  GlobalLogConfig.LoggerFileConfig.MaxSize,
			Rotate:   GlobalLogConfig.LoggerFileConfig.Rotate,
		}
	}

	for x := range SubLoggers {
		SubLoggers[x].SetLevels(splitLevel(GlobalLogConfig.Level))
		writers, err := getWriters(&GlobalLogConfig.SubLoggerConfig)
		if err != nil {
			return err
		}
		SubLoggers[x].SetOutput(writers)
	}
	logger = newLogger(GlobalLogConfig)
	return nil
}

func splitLevel(level string) (l Levels) {
	enabledLevels := strings.Split(level, "|")
	for x := range enabledLevels {
		switch level := enabledLevels[x]; level {
		case "DEBUG":
			l.Debug = true
		case "INFO":
			l.Info = true
		case "WARN":
			l.Warn = true
		case "ERROR":
			l.Error = true
		}
	}
	return
}

func registerNewSubLogger(subLogger string) *SubLogger {
	temp := &SubLogger{
		name:   strings.ToUpper(subLogger),
		output: os.Stdout,
		levels: splitLevel("INFO|WARN|DEBUG|ERROR"),
	}
	RWM.Lock()
	SubLoggers[subLogger] = temp
	RWM.Unlock()
	return temp
}
