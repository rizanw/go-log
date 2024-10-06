package zap

import (
	"context"
	"os"

	"github.com/rizanw/go-log/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger struct
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

	if config.DebugLog {
		configEncoder = zap.NewDevelopmentEncoderConfig()
	}

	// set zap config
	configEncoder.TimeKey = "timestamp"
	configEncoder.EncodeTime = zapcore.RFC3339TimeEncoder
	if config.TimeFormat != "" {
		configEncoder.EncodeTime = zapcore.TimeEncoderOfLayout(config.TimeFormat)
	}
	configEncoder.CallerKey = "line"
	callerSkipFrameCount := 4 + config.CallerSkip
	configEncoder.StacktraceKey = "stack"
	if !config.WithCaller {
		configEncoder.CallerKey = zapcore.OmitKey
	}

	// set output log
	zapEncoder := zapcore.NewConsoleEncoder(configEncoder)
	writer := zapcore.AddSync(os.Stderr)

	file, err := config.OpenLogFile()
	if err != nil {
		return nil, err
	}
	if file != nil {
		zapEncoder = zapcore.NewJSONEncoder(configEncoder)
		writer = zapcore.AddSync(file)
	}

	zapCore := zapcore.NewCore(zapEncoder, writer, setLevel(config.Level))
	zapLogger = zap.New(zapCore,
		zap.Fields(
			zap.String("app", config.AppName),
			zap.String("env", config.Environment),
		),
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

func buildFields(ctx context.Context, field logger.Field, err error) []zap.Field {
	zapFields := make([]zap.Field, 0)

	if field.RequestID != "" {
		zapFields = append(zapFields, zap.String(logger.FieldNameRequestID, field.RequestID))
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

func (l *Logger) Debug(ctx context.Context, field logger.Field, err error, message string) {
	l.logger.Debug(message, buildFields(ctx, field, err)...)
}
