package zap

import (
	"fmt"
	"os"

	"github.com/rizanw/go-log/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	config *logger.Config
}

func New(config *logger.Config) (*Logger, error) {
	var (
		configEncoder = zap.NewProductionEncoderConfig()
		zapLogger     *zap.Logger
		err           error
	)

	// set zap config
	configEncoder.MessageKey = "message"
	configEncoder.LevelKey = "level"
	configEncoder.TimeKey = "timestamp"
	configEncoder.EncodeTime = zapcore.RFC3339TimeEncoder
	if config.TimeFormat != "" {
		configEncoder.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeFormat)
	}
	configEncoder.StacktraceKey = "stacktrace"
	configEncoder.CallerKey = "line"
	callerSkipFrameCount := 2 + config.CallerSkip
	if !config.WithCaller {
		configEncoder.CallerKey = zapcore.OmitKey
	}

	// set output log
	zapEncoder := zapcore.NewJSONEncoder(configEncoder)
	writer := zapcore.AddSync(os.Stderr)

	if config.IsDevelopment {
		zapEncoder = zapcore.NewConsoleEncoder(configEncoder)
	}
	if !config.UseJSON {
		zapEncoder = zapcore.NewConsoleEncoder(configEncoder)
	}

	file, err := config.OpenLogFile()
	if err != nil {
		return nil, err
	}
	if file != nil {
		zapEncoder = zapcore.NewJSONEncoder(configEncoder)
		writer = zapcore.AddSync(file)
	}

	initialFields := make([]zap.Field, 0)
	if config.AppName != "" {
		initialFields = append(initialFields, zap.String("app", config.AppName))
	}
	if config.Environment != "" {
		initialFields = append(initialFields, zap.String("env", config.Environment))
	}

	zapCore := zapcore.NewCore(zapEncoder, writer, setLevel(config.Level))
	if config.UseMultiWriters {
		zapCore = zapcore.NewTee(
			zapcore.NewCore(zapEncoder, zapcore.Lock(file), setLevel(config.Level)),
			zapcore.NewCore(zapEncoder, zapcore.Lock(os.Stdout), setLevel(config.Level)),
		)
	}

	zapLogger = zap.New(zapCore,
		zap.Fields(initialFields...),
	)

	if config.WithCaller {
		zapLogger = zapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(callerSkipFrameCount))
	}

	if config.WithStack {
		zapLogger = zapLogger.WithOptions(zap.AddStacktrace(setLevel(config.StackLevel)))
	}

	defer zapLogger.Sync()
	return &Logger{
		logger: zapLogger,
		config: config,
	}, nil
}

func setLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.DebugLevel:
		return zap.DebugLevel
	case logger.InfoLevel:
		return zap.InfoLevel
	case logger.WarnLevel:
		return zap.WarnLevel
	case logger.ErrorLevel:
		return zap.ErrorLevel
	case logger.FatalLevel:
		return zap.FatalLevel
	default:
		return zap.DebugLevel
	}
}

func buildFields(cfg *logger.Config, field logger.Field, err error) []zap.Field {
	zapFields := make([]zap.Field, 0)

	if field.RequestID != "" {
		zapFields = append(zapFields, zap.String(logger.FieldNameRequestID, field.RequestID))
	}

	if field.Source != nil {
		zapFields = append(zapFields, zap.Any(logger.FieldNameSource, field.Source))
	}

	if field.UserInfo != nil {
		zapFields = append(zapFields, zap.Any(logger.FieldNameUserInfo, field.UserInfo))
	}

	if err != nil {
		zapFields = append(zapFields, zap.Error(err))
	}

	for key, value := range field.Fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	if len(field.Metadata) > 0 {
		metadata := field.Metadata
		if len(cfg.SensitiveFields) > 0 {
			cfg.MaskSensitiveData(metadata)
		}
		zapFields = append(zapFields, zap.Any(logger.FieldNameMetadata, metadata))
	}

	return zapFields
}

func (l *Logger) Debug(field logger.Field, err error, message string) {
	l.logger.Debug(message, buildFields(l.config, field, err)...)
}

func (l *Logger) Info(field logger.Field, err error, message string) {
	l.logger.Info(message, buildFields(l.config, field, err)...)
}

func (l *Logger) Warn(field logger.Field, err error, message string) {
	l.logger.Warn(message, buildFields(l.config, field, err)...)
}

func (l *Logger) Error(field logger.Field, err error, message string) {
	l.logger.Error(message, buildFields(l.config, field, err)...)
}

func (l *Logger) Fatal(field logger.Field, err error, message string) {
	l.logger.Fatal(message, buildFields(l.config, field, err)...)
}

func (l *Logger) Debugf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...), buildFields(l.config, field, err)...)
}

func (l *Logger) Infof(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...), buildFields(l.config, field, err)...)
}

func (l *Logger) Warnf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...), buildFields(l.config, field, err)...)
}

func (l *Logger) Errorf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...), buildFields(l.config, field, err)...)
}

func (l *Logger) Fatalf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...), buildFields(l.config, field, err)...)
}
