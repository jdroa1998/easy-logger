// Package logger provides a simple, efficient logging library based on zerolog.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger to provide additional functionality.
type Logger struct {
	zl          zerolog.Logger
	serviceName string
}

// LogBuilder provides a fluid interface for creating logs with formatted messages.
type LogBuilder struct {
	logger *Logger
	event  *zerolog.Event
	err    error
}

// Config contains configuration options for the logger.
type Config struct {
	// Level sets the minimum level of log messages to output
	Level Level
	// Pretty enables pretty, more human-readable logging format
	Pretty bool
	// WithCaller adds the caller information (file and line) to log entries
	WithCaller bool
	// Output is where log entries will be written. Defaults to os.Stderr if nil
	Output io.Writer
	// TimeFormat specifies the format for timestamps
	TimeFormat string
	// ServiceName identifies the service that generated the log
	ServiceName string
}

// DefaultConfig returns a default configuration for the logger.
func DefaultConfig() Config {
	return Config{
		Level:       InfoLevel,
		Pretty:      false,
		WithCaller:  true,
		Output:      os.Stderr,
		TimeFormat:  time.RFC3339,
		ServiceName: "",
	}
}

// New creates a new Logger with the given configuration.
func New(cfg Config) *Logger {
	output := cfg.Output
	if output == nil {
		output = os.Stderr
	}

	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = "UNKNOWN-SERVICE"
	}

	zctx := zerolog.New(output).
		Level(zerolog.Level(cfg.Level)).
		With()

	zctx = zctx.Timestamp()

	zctx = zctx.Str("service", serviceName)

	if cfg.WithCaller {
		zctx = zctx.Caller()
	}

	var zl zerolog.Logger
	if cfg.Pretty {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: cfg.TimeFormat,
		}
		zl = zctx.Logger().Output(consoleWriter)
	} else {
		zl = zctx.Logger()
	}

	zerolog.TimeFieldFormat = cfg.TimeFormat

	return &Logger{
		zl:          zl,
		serviceName: serviceName,
	}
}

// ServiceName returns the name of the service used by this logger
func (l *Logger) ServiceName() string {
	return l.serviceName
}

// SetServiceName sets the name of the service used by this logger
func (l *Logger) SetServiceName(name string) {
	l.serviceName = name
}

// With adds key-value pairs to the logger context
func (l *Logger) With() zerolog.Context {
	return l.zl.With()
}

// WithFields returns a new logger with the given fields added to the context.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	ctx := l.zl.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &Logger{
		zl:          ctx.Logger(),
		serviceName: l.serviceName,
	}
}

// SetLevel changes the log level of the logger
func (l *Logger) SetLevel(level Level) {
	l.zl = l.zl.Level(zerolog.Level(level))
}

// NewLogBuilder creates a new log builder instance
func (l *Logger) newLogBuilder(event *zerolog.Event) *LogBuilder {
	return &LogBuilder{
		logger: l,
		event:  event,
	}
}

// WithError adds an error to the log builder
func (lb *LogBuilder) WithError(err error) *LogBuilder {
	lb.event.Err(err)
	return lb
}

// Field adds a generic field to the log
func (lb *LogBuilder) AddField(key string, value any) *LogBuilder {
	lb.event.Interface(key, value)
	return lb
}

// Str adds a string field to the log
func (lb *LogBuilder) Str(key string, value string) *LogBuilder {
	lb.event.Str(key, value)
	return lb
}

// Int adds an integer field to the log
func (lb *LogBuilder) Int(key string, value int) *LogBuilder {
	lb.event.Int(key, value)
	return lb
}

// Bool adds a boolean field to the log
func (lb *LogBuilder) Bool(key string, value bool) *LogBuilder {
	lb.event.Bool(key, value)
	return lb
}

// Debug creates a debug level log
func (l *Logger) Debug() *LogBuilder {
	return l.newLogBuilder(l.zl.Debug())
}

// Debug creates a info level log
func (l *Logger) Info() *LogBuilder {
	return l.newLogBuilder(l.zl.Info())
}

// Warn creates a warn level log
func (l *Logger) Warn() *LogBuilder {
	return l.newLogBuilder(l.zl.Warn())
}

// Error creates an error level log
func (l *Logger) Error() *LogBuilder {
	return l.newLogBuilder(l.zl.Error())
}

// Fatal creates a fatal level log
func (l *Logger) Fatal() *LogBuilder {
	return l.newLogBuilder(l.zl.Fatal())
}

// Panic creates a panic level log
func (l *Logger) Panic() *LogBuilder {
	return l.newLogBuilder(l.zl.Panic())
}

// Trace creates a trace level log
func (l *Logger) Trace() *LogBuilder {
	return l.newLogBuilder(l.zl.Trace())
}

// Msg finalizes the log with a message
func (lb *LogBuilder) Msg(msg string, values ...any) {
	lb.event.Msgf(msg, values...)
}

// DebugMsg logs a simple message at debug level
func (l *Logger) DebugMsg(msg string, values ...any) {
	l.Debug().Msg(msg, values...)
}

// InfoMsg logs a simple message at info level
func (l *Logger) InfoMsg(msg string, values ...any) {
	l.Info().Msg(msg, values...)
}

// WarnMsg logs a simple message at warn level
func (l *Logger) WarnMsg(msg string, values ...any) {
	l.Warn().Msg(msg, values...)
}

// ErrorMsg logs a simple message at error level
func (l *Logger) ErrorMsg(msg string, values ...any) {
	l.Error().Msg(msg, values...)
}

// FatalMsg logs a simple message at fatal level, then calls os.Exit(1)
func (l *Logger) FatalMsg(msg string, values ...any) {
	l.Fatal().Msg(msg, values...)
}

// PanicMsg logs a simple message at panic level, then panics
func (l *Logger) PanicMsg(msg string, values ...any) {
	l.Panic().Msg(msg, values...)
}

// TraceMsg logs a simple message at trace level
func (l *Logger) TraceMsg(msg string, values ...any) {
	l.Trace().Msg(msg, values...)
}
