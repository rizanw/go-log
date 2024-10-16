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

	if config.IsDevelopment {
		configEncoder = zap.NewDevelopmentEncoderConfig()
	}

	// set zap config
	configEncoder.TimeKey = "timestamp"
	configEncoder.EncodeTime = zapcore.RFC3339TimeEncoder
	if config.TimeFormat != "" {
		configEncoder.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeFormat)
	}
	configEncoder.CallerKey = "line"
	callerSkipFrameCount := 2 + config.CallerSkip
	configEncoder.StacktraceKey = "stack"
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
	zapLogger = zap.New(zapCore,
		zap.Fields(initialFields...),
		zap.WithCaller(config.WithCaller),
		zap.AddCallerSkip(callerSkipFrameCount),
	)

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

func buildFields(field logger.Field, err error) []zap.Field {
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
		zapFields = append(zapFields, zap.Any(logger.FieldNameMetadata, field.Metadata))
	}

	return zapFields
}

func (l *Logger) Debug(field logger.Field, err error, message string) {
	l.logger.Debug(message, buildFields(field, err)...)
}

func (l *Logger) Info(field logger.Field, err error, message string) {
	l.logger.Info(message, buildFields(field, err)...)
}

func (l *Logger) Warn(field logger.Field, err error, message string) {
	l.logger.Warn(message, buildFields(field, err)...)
}

func (l *Logger) Error(field logger.Field, err error, message string) {
	l.logger.Error(message, buildFields(field, err)...)
}

func (l *Logger) Fatal(field logger.Field, err error, message string) {
	l.logger.Fatal(message, buildFields(field, err)...)
}

func (l *Logger) Debugf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...), buildFields(field, err)...)
}

func (l *Logger) Infof(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...), buildFields(field, err)...)
}

func (l *Logger) Warnf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...), buildFields(field, err)...)
}

func (l *Logger) Errorf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...), buildFields(field, err)...)
}

func (l *Logger) Fatalf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...), buildFields(field, err)...)
}
