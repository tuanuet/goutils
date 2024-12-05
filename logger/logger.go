package logger

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

// A global variable so that log functions can be directly accessed
var log Logger = DefaultLogger()

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// LoggerBackend reprents the int enum for backend of logger (exists for legacy compat reason)
type LoggerBackend int

const (
	// Debug has verbose message
	debugLvl = "debug"
	// Info is default log level
	infoLvl = "info"
	// Warn is for logging messages about possible issues
	warnLvl = "warn"
	// Error is for logging errors
	errorLvl = "error"
	// Fatal is for logging fatal messages. The sytem shutsdown after logging the message.
	fatalLvl = "fatal"
)

const (
	// LoggerBackendZap logging using Uber's zap backend
	LoggerBackendZap LoggerBackend = iota
	// LoggerBackendLogrus logging using logrus backend
	LoggerBackendLogrus
)

var (
	errInvalidLoggerInstance = errors.New("Invalid logger instance.")

	once sync.Once
)

// Logger is our contract for the logger
type Logger interface {
	Debug(msg string)
	Debugf(format string, args ...interface{})

	Info(msg string)
	Infof(format string, args ...interface{})

	Warn(msg string)
	Warnf(format string, args ...interface{})

	Error(msg string)
	Errorf(format string, args ...interface{})

	Fatal(msg string)
	Fatalf(format string, args ...interface{})

	Panic(msg string)
	Panicf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger

	GetDelegate() interface{}
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool   `name:"log-enable-console" help:"Enable log console" env:"LOG_ENABLE_CONSOLE" default:"true" yaml:"enable-console"`
	ConsoleJSONFormat bool   `name:"log-console-json-format" help:"Console to json format" env:"LOG_CONSOLE_JSON_FORMAT" default:"false" yaml:"console-json-format"`
	ConsoleLevel      string `name:"log-console-level" help:"Console log level" env:"LOG_CONSOLE_LEVEL" default:"info" enum:"debug, info, warn, error, fatal, panic" yaml:"console-level"`
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

// DefaultLogger creates default logger, which uses zap sugarlogger and outputs to console
func DefaultLogger() Logger {
	cfg := Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      "warn",
		EnableFile:        false,
		FileJSONFormat:    false,
	}
	logger, _ := newZapLogger(cfg)
	return logger
}

// InitLogger returns an instance of logger
func InitLogger(config Configuration, backend LoggerBackend) (Logger, error) {
	var err error
	once.Do(func() {
		switch backend {
		case LoggerBackendZap:
			log, err = NewLogger(config, backend)

		case LoggerBackendLogrus:
			log, err = NewLogger(config, backend)

		default:
			err = errInvalidLoggerInstance
		}
	})
	return log, err
}

func NewLogger(config Configuration, backend LoggerBackend) (Logger, error) {
	switch backend {
	case LoggerBackendZap:
		return newZapLogger(config)

	case LoggerBackendLogrus:
		return newLogrusLogger(config)

	default:
		return nil, errInvalidLoggerInstance
	}
}

func Debug(msg string) {
	log.Debugf(msg)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(msg string) {
	log.Infof(msg)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(msg string) {
	log.Warnf(msg)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(msg string) {
	log.Errorf(msg)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(msg string) {
	log.Fatalf(msg)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panic(msg string) {
	log.Panicf(msg)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// WithTraceIdPrefix add prefix with context as trace ID of span
func WithTraceIdPrefix(ctx context.Context) Logger {
	return log.WithFields(map[string]interface{}{
		"trace-id": trace.SpanContextFromContext(ctx).TraceID().String(),
	})
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

func GetDelegate() interface{} {
	return log.GetDelegate()
}

//go:generate mockery --case underscore --name=Logger --inpackage --filename=logger_mock.go
