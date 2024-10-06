package logger

import (
	"context"
	"os"
	"path/filepath"
)

type (
	// Level of log
	Level int

	// Engine of logger
	Engine string

	// Logger interface
	Logger interface {
		Debug(ctx context.Context, field Field, err error, message string)
		Debugf(requestID string, err error, fields interface{}, format string, args ...interface{})
		Info(ctx context.Context, field Field, err error, message string)
		Infof(requestID string, err error, fields interface{}, format string, args ...interface{})
		Warn(ctx context.Context, field Field, err error, message string)
		Warnf(requestID string, err error, fields interface{}, format string, args ...interface{})
		Error(ctx context.Context, field Field, err error, message string)
		Errorf(requestID string, err error, fields interface{}, format string, args ...interface{})
		Fatal(ctx context.Context, field Field, err error, message string)
		Fatalf(requestID string, err error, fields interface{}, format string, args ...interface{})
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
	Zap Engine = "zap"
)

type Config struct {
	AppName     string
	Environment string
	File        string
	TimeFormat  string
	Level       Level
	CallerSkip  int
	WithCaller  bool
	DebugLog    bool
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
