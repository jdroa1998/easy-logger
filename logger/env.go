package logger

import (
	"os"
	"strings"
)

const (
	// EnvLogLevel is the environment variable to configure log level
	EnvLogLevel = "LOG_LEVEL"
	// EnvLogFormat is the environment variable to configure log format
	EnvLogFormat = "LOG_FORMAT"
	// EnvLogCaller is the environment variable to configure if caller info is included
	EnvLogCaller = "LOG_CALLER"
	// EnvServiceName is the environment variable for service name
	EnvServiceName = "SERVICE_NAME"
)

// GetEnvStr returns the value of an environment variable or a default value
func GetEnvStr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvBool returns the boolean value of an environment variable
func GetEnvBool(key string, defaultValue bool) bool {
	if valueStr, exists := os.LookupEnv(key); exists {
		valueStr = strings.ToLower(valueStr)
		return valueStr == "true" || valueStr == "yes" || valueStr == "1"
	}
	return defaultValue
}

// NewFromEnv creates a new Logger configured with environment variables
func NewFromEnv() *Logger {
	// Get configuration from environment variables
	logLevel := GetEnvStr(EnvLogLevel, "info")
	logFormat := GetEnvStr(EnvLogFormat, "json")
	logCallerEnabled := GetEnvBool(EnvLogCaller, true)
	serviceName := GetEnvStr(EnvServiceName, "")

	// Determine log level
	level, err := ParseLevel(logLevel)
	if err != nil {
		level = InfoLevel
	}

	// Create configuration
	cfg := Config{
		Level:       level,
		Pretty:      logFormat == "pretty" || logFormat == "console",
		WithCaller:  logCallerEnabled,
		TimeFormat:  "2006-01-02T15:04:05.000Z07:00", // RFC3339 with milliseconds
		ServiceName: serviceName,
	}

	// Create logger
	logger := New(cfg)

	return logger
}

// InitGlobalFromEnv initializes a logger from environment variables and returns it
func InitGlobalFromEnv() *Logger {
	logger := NewFromEnv()
	return logger
}
