package logger

import (
	"os"
	"path/filepath"
	"strings"
)

type (
	// Level of log
	Level int

	// Engine of logger
	Engine string

	// ILogger interface
	ILogger interface {
		Debug(field Field, err error, message string)
		Debugf(field Field, err error, format string, args ...interface{})
		Info(field Field, err error, message string)
		Infof(field Field, err error, format string, args ...interface{})
		Warn(field Field, err error, message string)
		Warnf(field Field, err error, format string, args ...interface{})
		Error(field Field, err error, message string)
		Errorf(field Field, err error, format string, args ...interface{})
		Fatal(field Field, err error, message string)
		Fatalf(field Field, err error, format string, args ...interface{})
	}
)

// list of log level
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

const (
	EngineZap     Engine = "zap"
	EngineZerolog Engine = "zerolog"
)

type Config struct {
	AppName              string
	Environment          string
	IsDevelopment        bool
	TimeFormat           string
	Level                Level
	WithCaller           bool
	CallerSkip           int
	WithStack            bool
	StackLevel           Level
	StackMarshaller      func(err error) interface{}
	SensitiveFields      map[string]struct{}
	SensitiveFieldMasker func(value string) string
	UseJSON              bool
	UseColor             bool
	UseMultiWriters      bool
	File                 string
}

// OpenLogFile will open log file or generate it if not exist
func (c *Config) OpenLogFile() (*os.File, error) {
	if c.File == "" {
		return nil, nil
	}

	err := os.MkdirAll(filepath.Dir(c.File), 0755)
	if err != nil && err != os.ErrExist {
		return nil, err
	}

	return os.OpenFile(c.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}

// MaskSensitiveData recursively masks sensitive data in the map
func (c *Config) MaskSensitiveData(m map[string]interface{}) {
	if c.SensitiveFieldMasker == nil {
		c.SensitiveFieldMasker = maskFieldStar
	}
	for key, value := range m {
		switch v := value.(type) {
		case string:
			// Mask specific keys
			if _, ok := c.SensitiveFields[key]; ok {
				m[key] = c.SensitiveFieldMasker(v)
			}
		case map[string]interface{}:
			// Recursively call the function for nested maps
			c.MaskSensitiveData(v)
		}
	}
}

func maskFieldStar(s string) string {
	return strings.Repeat("*", len(s))
}
