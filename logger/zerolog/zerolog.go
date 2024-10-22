package zerolog

import (
	"io"
	"os"
	"time"

	"github.com/rizanw/go-log/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	logger *zerolog.Logger
	config *logger.Config
}

func New(config *logger.Config) (*Logger, error) {
	var (
		zeroLogger zerolog.Logger
		err        error
		writer     io.Writer
		timeFormat string = time.RFC3339
	)

	// set zerolog config
	if config.TimeFormat != "" {
		timeFormat = config.TimeFormat
	}
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = timeFormat
	zerolog.CallerFieldName = "line"
	callerSkipFrameCount := 4 + config.CallerSkip
	zerolog.ErrorStackFieldName = "stacktrace"
	if config.WithStack {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		if config.StackMarshaller != nil {
			zerolog.ErrorStackMarshaler = config.StackMarshaller
		}
	}

	// set output log
	writer = os.Stderr

	if config.IsDevelopment && config.UseColor {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: timeFormat,
		}
	} else if !config.UseJSON {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    true,
			TimeFormat: timeFormat,
		}
	}

	file, err := config.OpenLogFile()
	if err != nil {
		return nil, err
	}
	if file != nil {
		writer = file
	}

	if config.IsDevelopment {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if config.UseMultiWriters {
		writer = zerolog.MultiLevelWriter(file, os.Stdout)
	}

	zeroLogger = zerolog.New(writer).With().Timestamp().Logger().Level(setLevel(config.Level))
	if config.WithCaller {
		zeroLogger = zeroLogger.With().CallerWithSkipFrameCount(callerSkipFrameCount).Logger()
	}
	if config.AppName != "" {
		zeroLogger = zeroLogger.With().Str("app", config.AppName).Logger()
	}
	if config.Environment != "" {
		zeroLogger = zeroLogger.With().Str("env", config.Environment).Logger()
	}

	return &Logger{
		logger: &zeroLogger,
		config: config,
	}, nil
}

func setLevel(level logger.Level) zerolog.Level {
	switch level {
	case logger.DebugLevel:
		return zerolog.DebugLevel
	case logger.InfoLevel:
		return zerolog.InfoLevel
	case logger.WarnLevel:
		return zerolog.WarnLevel
	case logger.ErrorLevel:
		return zerolog.ErrorLevel
	case logger.FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.DebugLevel
	}
}

func buildFields(field logger.Field) map[string]interface{} {
	mapFields := make(map[string]interface{})

	if field.RequestID != "" {
		mapFields[logger.FieldNameRequestID] = field.RequestID
	}
	if field.Source != nil {
		mapFields[logger.FieldNameSource] = field.Source
	}
	if field.UserInfo != nil {
		mapFields[logger.FieldNameUserInfo] = field.UserInfo
	}

	for key, value := range field.Fields {
		mapFields[key] = value
	}

	if len(field.Metadata) > 0 {
		mapFields[logger.FieldNameMetadata] = field.Metadata
	}

	return mapFields
}

func (l *Logger) Debug(field logger.Field, err error, message string) {
	l.logger.Debug().Fields(buildFields(field)).Stack().Err(err).Msg(message)
}

func (l *Logger) Info(field logger.Field, err error, message string) {
	l.logger.Info().Fields(buildFields(field)).Stack().Err(err).Msg(message)
}

func (l *Logger) Warn(field logger.Field, err error, message string) {
	l.logger.Warn().Fields(buildFields(field)).Stack().Err(err).Msg(message)
}

func (l *Logger) Error(field logger.Field, err error, message string) {
	l.logger.Error().Fields(buildFields(field)).Stack().Err(err).Msg(message)
}

func (l *Logger) Fatal(field logger.Field, err error, message string) {
	l.logger.Fatal().Fields(buildFields(field)).Stack().Err(err).Msg(message)
}

func (l *Logger) Debugf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Debug().Fields(buildFields(field)).Stack().Err(err).Msgf(format, args...)
}

func (l *Logger) Infof(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Info().Fields(buildFields(field)).Stack().Err(err).Msgf(format, args...)
}

func (l *Logger) Warnf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Warn().Fields(buildFields(field)).Stack().Err(err).Msgf(format, args...)
}

func (l *Logger) Errorf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Error().Fields(buildFields(field)).Stack().Err(err).Msgf(format, args...)
}

func (l *Logger) Fatalf(field logger.Field, err error, format string, args ...interface{}) {
	l.logger.Fatal().Fields(buildFields(field)).Stack().Err(err).Msgf(format, args...)
}
