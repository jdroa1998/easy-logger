package logger

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
)

// Benchmark to measure the performance of zerolog directly
func BenchmarkZerologSimple(b *testing.B) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(io.Discard).With().Timestamp().Logger()

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		log.Info().Msg("This is a simple log message")
	}
}

// Benchmark to measure the performance of zerolog with structured fields
func BenchmarkZerologStructured(b *testing.B) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(io.Discard).With().Timestamp().Logger()

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		log.Info().
			Str("key1", "value1").
			Int("key2", 123).
			Bool("key3", true).
			Msg("This is a structured log message")
	}
}

// Benchmark to measure the performance of zerolog with caller enabled
func BenchmarkZerologWithCaller(b *testing.B) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(io.Discard).With().Timestamp().Caller().Logger()

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		log.Info().Msg("Log with caller information")
	}
}
