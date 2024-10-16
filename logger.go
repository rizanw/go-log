package log

import (
	"github.com/rizanw/go-log/logger"
	"github.com/rizanw/go-log/logger/zap"
	"github.com/rizanw/go-log/logger/zerolog"
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
	Zap     Engine = logger.EngineZap
	Zerolog Engine = logger.EngineZerolog
)

var (
	rlogger, _ = NewLogger(logger.Config{IsDevelopment: true}, logger.EngineZerolog)
)

// NewLogger creates a logger instance based on selected logger engine
func NewLogger(config logger.Config, engine logger.Engine) (Logger, error) {
	var (
		err error
		l   Logger
	)

	switch engine {
	case logger.EngineZerolog:
		l, err = zerolog.New(&config)
	case logger.EngineZap:
		l, err = zap.New(&config)
	default:
		l, err = zerolog.New(&config)
	}
	if err != nil {
		return nil, err
	}

	return l, err
}
