package main

import (
	"errors"
	"os"
	"time"

	"github.com/jdroa1998/easy-logger/logger"
)

func main() {
	// Create our logger instance
	log := logger.New(logger.DefaultConfig())

	// Example 1: Using basic methods
	log.InfoMsg("This is a simplified information message")
	log.Debug().Msg("This is a formatted message: %d", 42)

	// Example 2: Changing level to show more information
	log.SetLevel(logger.DebugLevel)
	log.DebugMsg("Now this debug message will appear")

	// Example 3: Adding structured fields
	log.Info().
		Str("service", "api").
		Int("port", 8080).
		Bool("debug", true).
		Msg("Server started")

	// Example 4: Creating a logger with builder pattern
	customLogger := logger.NewBuilder().
		WithLevel(logger.DebugLevel).
		WithPrettyPrint(true).
		WithServiceName("my-service").
		WithCaller(true).
		Build()

	customLogger.InfoMsg("Logger created with builder pattern")
	customLogger.Info().Msg("The service name is: %s", customLogger.ServiceName())

	// Example 5: Handling errors with formatted messages
	err := errors.New("something went wrong")
	log.Error().Msg("Error during processing: %v", err)

	// Traditional API for errors also works
	log.Error().
		WithError(err).
		Str("context", "processing file").
		Msg("Error during processing")

	// Example 6: Using WithFields for context with new methods
	requestLogger := log.WithFields(map[string]any{
		"request_id": "12345",
		"user_id":    "user-abc",
		"ip":         "192.168.1.1",
	})

	requestLogger.InfoMsg("Processing request")
	requestLogger.Warn().Msg("High response time: %dms", 500)

	// Example 7: Using predefined configurations with builder
	prodLogger := logger.NewBuilder().
		Production().
		WithServiceName("api-production").
		Build()

	prodLogger.InfoMsg("Application started in production mode")
	prodLogger.Info().Msg("Version: %s, Timestamp: %d", "1.0.0", time.Now().Unix())

	// Example 8: Complete configuration through builder
	advLogger := logger.NewBuilder().
		WithLevel(logger.WarnLevel).
		WithPrettyPrint(true).
		WithCaller(true).
		WithOutput(os.Stdout).
		WithTimeFormat(time.RFC822).
		WithServiceName("advanced-service").
		Build()

	advLogger.WarnMsg("This is a warning message with custom format")
	advLogger.InfoMsg("This info message won't appear due to the level set")

	// Example 9: Mixing styles in the same logger
	mixedLogger := logger.NewBuilder().
		Development().
		WithServiceName("mixed-service").
		Build()

	// Basic style
	mixedLogger.InfoMsg("Simple message")

	// Formatted style
	mixedLogger.Info().Msg("Formatted message: %s, %d", "hello", 123)

	// Complete structured style
	mixedLogger.Info().
		Str("key1", "value1").
		Int("key2", 456).
		Bool("key3", true).
		Msg("Structured message")
}
