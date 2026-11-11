package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

// ErrMultiWritersWithoutFile is returned when UseMultiWriters is set without a log file path.
var ErrMultiWritersWithoutFile = errors.New("go-log: UseMultiWriters requires FilePath")

// Validate checks the configuration for combinations that can't work.
func (c *Config) Validate() error {
	if c.UseMultiWriters && c.File == "" {
		return ErrMultiWritersWithoutFile
	}
	return nil
}

// OpenLogFile will open log file or generate it if not exist
func (c *Config) OpenLogFile() (*os.File, error) {
	if c.File == "" {
		return nil, nil
	}

	if err := os.MkdirAll(filepath.Dir(c.File), 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(c.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}

// IsSensitive reports whether a field key should be masked (case-insensitive).
func (c *Config) IsSensitive(key string) bool {
	_, ok := c.SensitiveFields[strings.ToLower(key)]
	return ok
}

// Mask returns a copy of value with every sensitive key masked, at any depth.
// Maps (including named map types such as log.KV) and slices are copied, so
// the caller's data is never modified. Non-sensitive values are kept as is.
func (c *Config) Mask(value interface{}) interface{} {
	if len(c.SensitiveFields) == 0 || value == nil {
		return value
	}
	return c.maskValue(reflect.ValueOf(value))
}

// MaskSensitiveData returns a masked copy of m. See Mask.
func (c *Config) MaskSensitiveData(m map[string]interface{}) map[string]interface{} {
	if len(c.SensitiveFields) == 0 || m == nil {
		return m
	}
	masked, _ := c.maskValue(reflect.ValueOf(m)).(map[string]interface{})
	return masked
}

func (c *Config) maskValue(v reflect.Value) interface{} {
	if !v.IsValid() {
		return nil
	}
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return v.Interface()
		}
		if v.Kind() == reflect.Pointer && v.Elem().Kind() != reflect.Map && v.Elem().Kind() != reflect.Slice {
			return v.Interface() // leave pointers to structs and scalars untouched
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return v.Interface()
		}
		out := make(map[string]interface{}, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			if c.IsSensitive(key) {
				out[key] = c.maskLeaf(iter.Value())
				continue
			}
			out[key] = c.maskValue(iter.Value())
		}
		return out
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			return v.Interface()
		}
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return v.Interface() // []byte
		}
		out := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			out[i] = c.maskValue(v.Index(i))
		}
		return out
	default:
		return v.Interface()
	}
}

func (c *Config) maskLeaf(v reflect.Value) string {
	masker := c.SensitiveFieldMasker
	if masker == nil {
		masker = maskFieldStar
	}
	if !v.IsValid() {
		return masker("")
	}
	return masker(fmt.Sprint(v.Interface()))
}

func maskFieldStar(s string) string {
	return strings.Repeat("*", len(s))
}
