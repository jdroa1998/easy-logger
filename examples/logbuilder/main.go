package main

import (
	"errors"

	"github.com/jdroa1998/easy-logger/logger"
)

func main() {
	// Create a default logger
	log := logger.New(logger.DefaultConfig())

	// 1. Using the builder pattern with structured logging
	log.Info().
		Str("service", "api").
		Int("port", 8080).
		Bool("active", true).
		Msg("Starting server on port %d", 8080)

	// 2. With an error included
	err := errors.New("connection refused")
	log.Error().
		WithError(err).
		Str("address", "192.168.1.1").
		Msg("Error connecting: %v", err)

	// 3. Using simplified methods
	log.InfoMsg("User admin authenticated successfully")
	log.Info().Msg("User %s authenticated successfully", "admin")

	// 4. With a custom logger using builder pattern
	customLogger := logger.NewBuilder().
		WithLevel(logger.DebugLevel).
		WithPrettyPrint(true).
		WithServiceName("auth-service").
		Build()

	customLogger.Info().
		Str("user_id", "12345").
		Int("login_attempts", 3).
		Msg("Processing login request")

	// 5. Mixed use with both API styles in the same program
	// Traditional style
	log.Info().
		Str("method", "GET").
		Int("status", 200).
		Msg("Request completed")

	// With formatted message
	log.Info().
		Str("method", "POST").
		Int("status", 201).
		Msg("Resource created with id %d", 12345)
}
