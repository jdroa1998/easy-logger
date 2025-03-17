package logger

import (
	"io"
	"time"
)

// Option is a function that configures a Logger.
type Option func(*Config)

// WithLevel sets the minimum log level.
func WithLevel(level Level) Option {
	return func(c *Config) {
		c.Level = level
	}
}

// WithPrettyPrint enables or disables pretty printing.
func WithPrettyPrint(enabled bool) Option {
	return func(c *Config) {
		c.Pretty = enabled
	}
}

// WithCaller enables or disables including the caller in log entries.
func WithCaller(enabled bool) Option {
	return func(c *Config) {
		c.WithCaller = enabled
	}
}

// WithOutput sets the output writer for the logger.
func WithOutput(w io.Writer) Option {
	return func(c *Config) {
		c.Output = w
	}
}

// WithTimeFormat sets the time format for the logger.
func WithTimeFormat(format string) Option {
	return func(c *Config) {
		c.TimeFormat = format
	}
}

// NewWithOptions creates a new logger with the provided options.
func NewWithOptions(opts ...Option) *Logger {
	cfg := DefaultConfig()

	// Apply all options
	for _, opt := range opts {
		opt(&cfg)
	}

	return New(cfg)
}

// Default creates a logger with reasonable default settings:
// - Info level
// - Pretty printing enabled
// - Caller information included
// - Output to stderr
// - RFC3339 time format
func Default() *Logger {
	return NewWithOptions(
		WithLevel(InfoLevel),
		WithPrettyPrint(true),
		WithCaller(true),
		WithTimeFormat(time.RFC3339),
	)
}

// Production creates a logger with production-suitable settings:
// - Info level
// - JSON format (no pretty printing)
// - No caller information (for performance)
// - Output to stderr
// - RFC3339 time format
func Production() *Logger {
	return NewWithOptions(
		WithLevel(InfoLevel),
		WithPrettyPrint(false),
		WithCaller(false),
	)
}

// Development creates a logger with development-suitable settings:
// - Debug level
// - Pretty printing enabled
// - Caller information included
// - Output to stderr
// - Time format that includes milliseconds
func Development() *Logger {
	return NewWithOptions(
		WithLevel(DebugLevel),
		WithPrettyPrint(true),
		WithCaller(true),
		WithTimeFormat(time.RFC3339Nano),
	)
}
