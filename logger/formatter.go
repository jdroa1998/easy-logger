package logger

import (
	"io"

	"github.com/rs/zerolog"
)

// Formatter defines an interface for custom log formatting.
type Formatter interface {
	Format(w io.Writer) io.Writer
}

// JSONFormatter formats logs in JSON format.
type JSONFormatter struct{}

// Format returns a writer that formats logs in JSON.
func (f JSONFormatter) Format(w io.Writer) io.Writer {
	return w
}

// PrettyFormatter formats logs in a human-readable format.
type PrettyFormatter struct {
	// NoColor disables colors in the output
	NoColor bool
	// TimeFormat sets the format for timestamps
	TimeFormat string
}

// Format returns a writer that formats logs in a pretty, human-readable format.
func (f PrettyFormatter) Format(w io.Writer) io.Writer {
	output := zerolog.ConsoleWriter{
		Out:        w,
		NoColor:    f.NoColor,
		TimeFormat: f.TimeFormat,
	}
	return output
}

// DefaultJSONFormatter returns a new JSONFormatter with default settings.
func DefaultJSONFormatter() Formatter {
	return JSONFormatter{}
}

// DefaultPrettyFormatter returns a new PrettyFormatter with default settings.
func DefaultPrettyFormatter() Formatter {
	return PrettyFormatter{
		NoColor:    false,
		TimeFormat: zerolog.TimeFieldFormat,
	}
}
