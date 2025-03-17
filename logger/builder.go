package logger

import (
	"io"
	"time"
)

// LoggerBuilder provides a builder pattern for constructing a logger
type LoggerBuilder struct {
	config Config
}

// NewBuilder creates a new LoggerBuilder with default configuration
func NewBuilder() *LoggerBuilder {
	return &LoggerBuilder{
		config: DefaultConfig(),
	}
}

// WithLevel sets the logging level
func (b *LoggerBuilder) WithLevel(level Level) *LoggerBuilder {
	b.config.Level = level
	return b
}

// WithPrettyPrint enables or disables pretty format for terminal output
func (b *LoggerBuilder) WithPrettyPrint(enabled bool) *LoggerBuilder {
	b.config.Pretty = enabled
	return b
}

// WithCaller enables or disables caller information
func (b *LoggerBuilder) WithCaller(enabled bool) *LoggerBuilder {
	b.config.WithCaller = enabled
	return b
}

// WithOutput sets the destination for log output
func (b *LoggerBuilder) WithOutput(output io.Writer) *LoggerBuilder {
	b.config.Output = output
	return b
}

// WithTimeFormat sets the time format
func (b *LoggerBuilder) WithTimeFormat(format string) *LoggerBuilder {
	b.config.TimeFormat = format
	return b
}

// WithServiceName sets the service name to identify logs
func (b *LoggerBuilder) WithServiceName(name string) *LoggerBuilder {
	b.config.ServiceName = name
	return b
}

// Development configures the builder with optimal settings for development
func (b *LoggerBuilder) Development() *LoggerBuilder {
	b.config.Level = DebugLevel
	b.config.Pretty = true
	b.config.WithCaller = true
	b.config.TimeFormat = time.RFC3339Nano
	return b
}

// Production configures the builder with optimal settings for production
func (b *LoggerBuilder) Production() *LoggerBuilder {
	b.config.Level = InfoLevel
	b.config.Pretty = false
	b.config.WithCaller = false
	b.config.TimeFormat = time.RFC3339
	return b
}

// Build constructs and returns the configured logger
func (b *LoggerBuilder) Build() *Logger {
	return New(b.config)
}

// BuildAndSetAsDefault builds the logger and returns it
// NOTE: This function has been modified to not use global state
func (b *LoggerBuilder) BuildAndSetAsDefault() *Logger {
	logger := b.Build()
	return logger
}
