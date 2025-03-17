package logger

import (
	"io"
	"testing"
)

// Benchmark to measure performance and memory usage of logger creation
func BenchmarkLoggerCreation(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = New(DefaultConfig())
	}
}

// Benchmark to measure performance of the builder pattern
func BenchmarkBuilderPattern(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = NewBuilder().
			WithLevel(DebugLevel).
			WithPrettyPrint(false).
			WithServiceName("benchmark-service").
			Build()
	}
}

// Benchmark to measure performance of a simple log
func BenchmarkSimpleLog(b *testing.B) {
	logger := New(Config{
		Level:       InfoLevel,
		Pretty:      false,
		WithCaller:  false,
		Output:      io.Discard, // Using io.Discard to avoid measuring write time
		ServiceName: "benchmark-service",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.InfoMsg("This is a simple log message")
	}
}

// Benchmark to measure performance of a structured log
func BenchmarkStructuredLog(b *testing.B) {
	logger := New(Config{
		Level:       InfoLevel,
		Pretty:      false,
		WithCaller:  false,
		Output:      io.Discard,
		ServiceName: "benchmark-service",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.Info().
			Str("key1", "value1").
			Int("key2", 123).
			Bool("key3", true).
			Msg("This is a structured log message")
	}
}

// Benchmark to measure performance of a formatted log
func BenchmarkFormattedLog(b *testing.B) {
	logger := New(Config{
		Level:       InfoLevel,
		Pretty:      false,
		WithCaller:  false,
		Output:      io.Discard,
		ServiceName: "benchmark-service",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.InfoMsg("This is a formatted message with values: %s, %d", "string", 123)
	}
}

// Benchmark to measure creation and usage of logger with contextual fields
func BenchmarkLoggerWithContext(b *testing.B) {
	logger := New(Config{
		Level:       InfoLevel,
		Pretty:      false,
		WithCaller:  false,
		Output:      io.Discard,
		ServiceName: "benchmark-service",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		contextLogger := logger.WithFields(map[string]any{
			"request_id": "12345",
			"user_id":    "user-abc",
		})
		contextLogger.InfoMsg("Log with context")
	}
}

// Benchmark to measure performance of Pretty mode (formatted output)
func BenchmarkPrettyLog(b *testing.B) {
	logger := New(Config{
		Level:       InfoLevel,
		Pretty:      true, // Enable pretty mode
		WithCaller:  false,
		Output:      io.Discard,
		ServiceName: "benchmark-service",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.Info().
			Str("key1", "value1").
			Int("key2", 123).
			Msg("This is a pretty log message")
	}
}

// Benchmark to measure performance with caller information enabled
func BenchmarkLoggerWithCaller(b *testing.B) {
	logger := New(Config{
		Level:       InfoLevel,
		Pretty:      false,
		WithCaller:  true, // Enable caller info
		Output:      io.Discard,
		ServiceName: "benchmark-service",
	})

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		logger.InfoMsg("Log with caller information")
	}
}
