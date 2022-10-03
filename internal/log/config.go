package log

import "github.com/quantstop/quantstopterminal/internal/system/convert"

// Config holds configuration settings loaded from bot config
type Config struct {
	Enabled *bool `json:"enabled"`
	SubLoggerConfig
	LoggerFileConfig *loggerFileConfig `json:"fileSettings,omitempty"`
	AdvancedSettings advancedSettings  `json:"advancedSettings"`
	SubLoggers       []SubLoggerConfig `json:"subLoggers,omitempty"`
}

// SubLoggerConfig holds sub logger configuration settings loaded from bot config
type SubLoggerConfig struct {
	Name   string `json:"name,omitempty"`
	Level  string `json:"level"`
	Output string `json:"output"`
}

type loggerFileConfig struct {
	FileName string `json:"filename,omitempty"`
	Rotate   *bool  `json:"rotate,omitempty"`
	MaxSize  int64  `json:"maxsize,omitempty"`
}

type advancedSettings struct {
	ShowLogSystemName *bool   `json:"showLogSystemName"`
	Spacer            string  `json:"spacer"`
	TimeStampFormat   string  `json:"timeStampFormat"`
	Headers           headers `json:"headers"`
}

type headers struct {
	Info  string `json:"info"`
	Warn  string `json:"warn"`
	Debug string `json:"debug"`
	Error string `json:"error"`
}

// NewConfig returns a pointer to a Config struct with default settings
func NewConfig() *Config {
	return &Config{
		Enabled: convert.BoolPtr(true),
		SubLoggerConfig: SubLoggerConfig{
			Level:  "INFO|DEBUG|WARN|ERROR",
			Output: "console|file",
		},
		LoggerFileConfig: &loggerFileConfig{
			FileName: "log.txt",
			Rotate:   convert.BoolPtr(true),
			MaxSize:  100,
		},
		AdvancedSettings: advancedSettings{
			ShowLogSystemName: convert.BoolPtr(true),
			Spacer:            spacer,
			TimeStampFormat:   timestampFormat,
			Headers: headers{
				Info:  "[INFO]",
				Warn:  "[WARN]",
				Debug: "[DEBUG]",
				Error: "[ERROR]",
			},
		},
	}
}
