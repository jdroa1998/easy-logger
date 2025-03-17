package logger

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// Level represents a logging level.
type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = Level(zerolog.DebugLevel)
	// InfoLevel defines info log level.
	InfoLevel Level = Level(zerolog.InfoLevel)
	// WarnLevel defines warn log level.
	WarnLevel Level = Level(zerolog.WarnLevel)
	// ErrorLevel defines error log level.
	ErrorLevel Level = Level(zerolog.ErrorLevel)
	// FatalLevel defines fatal log level.
	FatalLevel Level = Level(zerolog.FatalLevel)
	// PanicLevel defines panic log level.
	PanicLevel Level = Level(zerolog.PanicLevel)
	// TraceLevel defines trace log level.
	TraceLevel Level = Level(zerolog.TraceLevel)
)

// ParseLevel converts a level string to a Level.
// Returns an error if the level string is invalid.
func ParseLevel(levelStr string) (Level, error) {
	levelStr = strings.ToLower(levelStr)
	switch levelStr {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "panic":
		return PanicLevel, nil
	case "trace":
		return TraceLevel, nil
	}
	return InfoLevel, fmt.Errorf("invalid log level: %s", levelStr)
}

// String converts a Level to a string.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case TraceLevel:
		return "trace"
	}
	return "unknown"
}
