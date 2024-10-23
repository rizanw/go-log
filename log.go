package log

import (
	"context"

	"github.com/rizanw/go-log/logger"
)

// Config for Log configuration
type Config struct {
	// AppName is your application name
	// it will be printed as `app` in log
	AppName string

	// Environment is your application environment running on
	// it will be printed as `env` in log
	// `dev` | `development` | `local` will mark your app under development env
	Environment string

	// Level is minimum log level to be printed (default: DEBUG)
	Level Level

	// TimeFormat is for log time format (default: RFC3339)
	TimeFormat string

	// WithCaller toggle to print which line is calling the log (default: false)
	WithCaller bool

	// CallerSkip is offset number for which caller line you wants to be print (default: 0)
	CallerSkip int

	// WithStack is a toggle to print which stack trace error located (default: false)
	WithStack bool

	// StackLevel is minimum log level for zap stack trace (default: ERROR)
	StackLevel *Level

	// StackMarshaller, function to get and log the stack trace for zerolog (default: `zerolog/pkgerrors`)
	StackMarshaller func(err error) interface{}

	// MaskSensitiveData is keys of field to be masked
	MaskSensitiveData []string

	// SensitiveDataMasker, function to modify sensitive value into something (default: `*****`)
	SensitiveDataMasker func(value string) string

	// UseJSON is a toggle to format log as json (default: false)
	UseJSON bool

	// UseColor is a toggle to colorize your log console
	// note: it only works using `zerolog` engine and under `development` environment
	UseColor bool

	// UseMultiWriters is a toggle to print log into log file and log console
	// note: FilePath must be filled
	UseMultiWriters bool

	// FilePath a file path to write the log as a file
	// note: if you fill the file path, your console log will be empty.
	FilePath string

	// Engine is logger to be used
	Engine Engine
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
		isDevelopment := false
		if config.Environment == "development" || config.Environment == "local" || config.Environment == "dev" {
			isDevelopment = true
		}

		errStackLevel := ErrorLevel
		if config.StackLevel != nil {
			errStackLevel = *config.StackLevel
		}

		maskSensitiveData := make(map[string]struct{})
		for _, key := range config.MaskSensitiveData {
			maskSensitiveData[key] = struct{}{}
		}

		configLogger = logger.Config{
			AppName:              config.AppName,
			Environment:          config.Environment,
			IsDevelopment:        isDevelopment,
			TimeFormat:           config.TimeFormat,
			Level:                config.Level,
			WithCaller:           config.WithCaller,
			CallerSkip:           config.CallerSkip,
			WithStack:            config.WithStack,
			StackLevel:           errStackLevel,
			StackMarshaller:      config.StackMarshaller,
			UseJSON:              config.UseJSON,
			UseColor:             config.UseColor,
			SensitiveFields:      maskSensitiveData,
			SensitiveFieldMasker: config.SensitiveDataMasker,
			UseMultiWriters:      config.UseMultiWriters,
			File:                 config.FilePath,
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
