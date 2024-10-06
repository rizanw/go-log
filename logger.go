package log

import (
	"github.com/rizanw/go-log/logger"
	"github.com/rizanw/go-log/logger/zap"
)

type (
	// Level of log
	Level = logger.Level

	// Engine logger
	Engine = logger.Engine

	// Logger interface
	Logger = logger.ILogger
)

// Level options
const (
	DebugLevel = logger.DebugLevel
	InfoLevel  = logger.InfoLevel
	WarnLevel  = logger.WarnLevel
	ErrorLevel = logger.ErrorLevel
	FatalLevel = logger.FatalLevel
)

// Engine options
const (
	Zap Engine = logger.EngineZap
)

var (
	rlogger, _ = NewLogger(logger.Config{IsDevelopment: true}, logger.EngineZap)
)

// NewLogger creates a logger instance based on selected logger engine
func NewLogger(config logger.Config, engine logger.Engine) (Logger, error) {
	var (
		err error
		l   Logger
	)

	switch engine {
	case logger.EngineZap:
		l, err = zap.New(&config)
	default:
		l, err = zap.New(&config)
	}
	if err != nil {
		return nil, err
	}

	return l, err
}
