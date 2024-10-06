package log

import (
	"context"

	"github.com/rizanw/go-log/logger"
)

type Config struct {
	AppName     string
	Environment string
	Level       Level
	TimeFormat  string
	WithCaller  bool
	CallerSkip  int
	FilePath    string
	Engine      Engine
}

// SetConfig is function to customize log configuration
func SetConfig(config *Config) error {
	var (
		err          error
		newLogger    logger.ILogger
		configLogger logger.Config
		engineLogger logger.Engine
	)

	if config != nil {
		configLogger = logger.Config{
			AppName:       config.AppName,
			Environment:   config.Environment,
			IsDevelopment: false,
			File:          config.FilePath,
			TimeFormat:    config.TimeFormat,
			Level:         config.Level,
			CallerSkip:    config.CallerSkip,
			WithCaller:    config.WithCaller,
		}
		engineLogger = config.Engine
	}

	newLogger, err = NewLogger(configLogger, engineLogger)
	if err != nil {
		return err
	}
	rlogger = newLogger
	return nil
}

// Debug prints log on debug level
func Debug(ctx context.Context, err error, metadata KV, message string) {
	rlogger.Debug(buildFields(ctx, metadata), err, message)
}

// Info prints log on info level
func Info(ctx context.Context, err error, metadata KV, message string) {
	rlogger.Info(buildFields(ctx, metadata), err, message)
}

// Warn prints log on warn level
func Warn(ctx context.Context, err error, metadata KV, message string) {
	rlogger.Warn(buildFields(ctx, metadata), err, message)
}

// Error prints log on error level
func Error(ctx context.Context, err error, metadata KV, message string) {
	rlogger.Error(buildFields(ctx, metadata), err, message)
}

// Fatal prints log on fatal level
func Fatal(ctx context.Context, err error, metadata KV, message string) {
	rlogger.Fatal(buildFields(ctx, metadata), err, message)
}

// Debugf prints log on debug level like fmt.Printf
func Debugf(ctx context.Context, err error, metadata KV, formatedMsg string, args ...interface{}) {
	rlogger.Debugf(buildFields(ctx, metadata), err, formatedMsg, args...)
}

// Infof prints log on info level like fmt.Printf
func Infof(ctx context.Context, err error, metadata KV, formatedMsg string, args ...interface{}) {
	rlogger.Infof(buildFields(ctx, metadata), err, formatedMsg, args...)
}

// Warnf prints log on warn level like fmt.Printf
func Warnf(ctx context.Context, err error, metadata KV, formatedMsg string, args ...interface{}) {
	rlogger.Warnf(buildFields(ctx, metadata), err, formatedMsg, args...)
}

// Errorf prints log on error level like fmt.printf
func Errorf(ctx context.Context, err error, metadata KV, formatedMsg string, args ...interface{}) {
	rlogger.Errorf(buildFields(ctx, metadata), err, formatedMsg, args...)
}

// Fatalf prints log on fatal level like fmt.printf
func Fatalf(ctx context.Context, err error, metadata KV, formatedMsg string, args ...interface{}) {
	rlogger.Fatalf(buildFields(ctx, metadata), err, formatedMsg, args...)
}
