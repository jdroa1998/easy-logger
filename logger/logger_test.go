package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

// TestLogLevels tests that log levels work correctly
func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger with custom output to capture messages
	log := New(Config{
		Level:      DebugLevel,
		Pretty:     false, // Use JSON for easier testing
		WithCaller: false,
		Output:     &buf,
	})

	// Test debug level
	log.Debug().Msg("debug message")
	assertLogContains(t, buf.String(), "debug message", "debug")
	buf.Reset()

	// Test info level
	log.Info().Msg("info message")
	assertLogContains(t, buf.String(), "info message", "info")
	buf.Reset()

	// Test warn level
	log.Warn().Msg("warn message")
	assertLogContains(t, buf.String(), "warn message", "warn")
	buf.Reset()

	// Test error level
	log.Error().Msg("error message")
	assertLogContains(t, buf.String(), "error message", "error")
	buf.Reset()

	// Change level and verify lower level messages are not logged
	log.SetLevel(WarnLevel)

	log.Debug().Msg("hidden debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not have been logged")
	}
	buf.Reset()

	log.Info().Msg("hidden info message")
	if buf.Len() > 0 {
		t.Error("Info message should not have been logged")
	}
	buf.Reset()

	log.Warn().Msg("visible warn message")
	if buf.Len() == 0 {
		t.Error("Warn message should have been logged")
	}
	buf.Reset()
}

// TestLogBuilder tests the new LogBuilder pattern
func TestLogBuilder(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Level:      DebugLevel,
		Pretty:     false, // Use JSON for easier testing
		WithCaller: false,
		Output:     &buf,
	})

	// Test creating a log with fields and a message
	log.Info().
		Str("key1", "value1").
		Int("key2", 42).
		Bool("key3", true).
		Msg("structured message")

	logData := buf.String()

	// Verify the message contains all fields
	assertLogContains(t, logData, "structured message", "info")
	assertLogContains(t, logData, "value1", "")
	assertLogContains(t, logData, "42", "")
	assertLogContains(t, logData, "true", "")
	buf.Reset()

	// Test with an error
	err := new(mockError)
	log.Error().
		WithError(err).
		Str("operation", "test").
		Msg("error occurred")

	logData = buf.String()
	assertLogContains(t, logData, "error occurred", "error")
	assertLogContains(t, logData, "operation", "")
	assertLogContains(t, logData, "error", "")
}

// TestWithFields verifies that the WithFields method works correctly
func TestWithFields(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Level:      InfoLevel,
		Pretty:     false,
		WithCaller: false,
		Output:     &buf,
	})

	// Create a logger with context fields
	contextLogger := log.WithFields(map[string]any{
		"request_id": "12345",
		"user_id":    "abc",
	})

	// Verify context fields appear in all messages
	contextLogger.Info().Msg("first message")
	assertLogContains(t, buf.String(), "12345", "")
	assertLogContains(t, buf.String(), "abc", "")
	buf.Reset()

	// Add more fields and test
	contextLogger.Error().
		Str("extra", "field").
		Msg("second message")

	assertLogContains(t, buf.String(), "12345", "")
	assertLogContains(t, buf.String(), "abc", "")
	assertLogContains(t, buf.String(), "field", "")
}

// TestFormattedLogs tests logging with formatted messages
func TestFormattedLogs(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Level:      DebugLevel,
		Pretty:     false,
		WithCaller: false,
		Output:     &buf,
	})

	// Test with formatting parameters
	log.Info().Msg("Value: %d", 42)

	assertLogContains(t, buf.String(), "Value: 42", "")
	buf.Reset()

	// Test with multiple values
	log.Debug().Msg("Values: %s, %d, %t", "test", 123, true)

	assertLogContains(t, buf.String(), "Values: test, 123, true", "debug")
}

// TestServiceNameField verifies that the service field is present in logs
func TestServiceNameField(t *testing.T) {
	var buf bytes.Buffer

	// Test default service name
	defaultLogger := New(Config{
		Level:      InfoLevel,
		Pretty:     false,
		WithCaller: false,
		Output:     &buf,
	})

	defaultLogger.InfoMsg("Message from default logger")

	var logData map[string]any
	if err := json.Unmarshal([]byte(buf.String()), &logData); err != nil {
		t.Errorf("Could not parse log as JSON: %v", err)
		return
	}

	service, ok := logData["service"].(string)
	if !ok {
		t.Error("Log does not contain 'service' field")
		return
	}

	if service != "UNKNOWN-SERVICE" {
		t.Errorf("Expected service name 'UNKNOWN-SERVICE', but got '%s'", service)
	}
	buf.Reset()

	// Test custom service name
	customLogger := New(Config{
		Level:       InfoLevel,
		Pretty:      false,
		WithCaller:  false,
		Output:      &buf,
		ServiceName: "test-service",
	})

	customLogger.InfoMsg("Message from custom logger")

	if err := json.Unmarshal([]byte(buf.String()), &logData); err != nil {
		t.Errorf("Could not parse log as JSON: %v", err)
		return
	}

	service, ok = logData["service"].(string)
	if !ok {
		t.Error("Log does not contain 'service' field")
		return
	}

	if service != "test-service" {
		t.Errorf("Expected service name 'test-service', but got '%s'", service)
	}
}

// Helper function to verify log contains specific text
func assertLogContains(t *testing.T, logLine, expectedText, expectedLevel string) {
	if !strings.Contains(logLine, expectedText) {
		t.Errorf("Log should contain '%s', but contains: %s", expectedText, logLine)
	}

	if expectedLevel != "" {
		var logData map[string]any
		if err := json.Unmarshal([]byte(logLine), &logData); err != nil {
			t.Errorf("Could not parse log as JSON: %v", err)
			return
		}

		level, ok := logData["level"].(string)
		if !ok {
			t.Error("Log does not contain 'level' field")
			return
		}

		if level != expectedLevel {
			t.Errorf("Expected level '%s', but got '%s'", expectedLevel, level)
		}
	}
}

// Mock error for testing error logging
type mockError struct{}

func (m *mockError) Error() string {
	return "mock error"
}

// TestBuilderPattern tests the builder pattern functionality
func TestBuilderPattern(t *testing.T) {
	var buf bytes.Buffer

	// Test basic builder creation
	builder := NewBuilder()
	if builder == nil {
		t.Error("NewBuilder should return a non-nil LoggerBuilder")
	}

	// Test WithLevel
	levelBuilder := builder.WithLevel(DebugLevel)
	if levelBuilder != builder {
		t.Error("WithLevel should return the same builder instance")
	}
	if builder.config.Level != DebugLevel {
		t.Errorf("Expected level %v, got %v", DebugLevel, builder.config.Level)
	}

	// Test WithPrettyPrint
	prettyBuilder := builder.WithPrettyPrint(true)
	if prettyBuilder != builder {
		t.Error("WithPrettyPrint should return the same builder instance")
	}
	if !builder.config.Pretty {
		t.Error("Pretty should be true after setting it")
	}

	// Test WithCaller
	callerBuilder := builder.WithCaller(false)
	if callerBuilder != builder {
		t.Error("WithCaller should return the same builder instance")
	}
	if builder.config.WithCaller {
		t.Error("WithCaller should be false after setting it to false")
	}

	// Test WithOutput
	outputBuilder := builder.WithOutput(&buf)
	if outputBuilder != builder {
		t.Error("WithOutput should return the same builder instance")
	}
	if builder.config.Output != &buf {
		t.Error("Output should be set to the provided buffer")
	}

	// Test WithTimeFormat
	timeFormat := "2006-01-02T15:04:05"
	timeFormatBuilder := builder.WithTimeFormat(timeFormat)
	if timeFormatBuilder != builder {
		t.Error("WithTimeFormat should return the same builder instance")
	}
	if builder.config.TimeFormat != timeFormat {
		t.Errorf("Expected TimeFormat %s, got %s", timeFormat, builder.config.TimeFormat)
	}

	// Test WithServiceName
	serviceName := "test-service"
	serviceNameBuilder := builder.WithServiceName(serviceName)
	if serviceNameBuilder != builder {
		t.Error("WithServiceName should return the same builder instance")
	}
	if builder.config.ServiceName != serviceName {
		t.Errorf("Expected ServiceName %s, got %s", serviceName, builder.config.ServiceName)
	}

	// Test Development
	devBuilder := NewBuilder().Development()
	if devBuilder.config.Level != DebugLevel || !devBuilder.config.Pretty || !devBuilder.config.WithCaller {
		t.Error("Development should set appropriate development settings")
	}

	// Test Production
	prodBuilder := NewBuilder().Production()
	if prodBuilder.config.Level != InfoLevel || prodBuilder.config.Pretty || prodBuilder.config.WithCaller {
		t.Error("Production should set appropriate production settings")
	}

	// Test Build
	logger := builder.Build()
	if logger == nil {
		t.Error("Build should return a non-nil Logger")
	}
	if logger.ServiceName() != serviceName {
		t.Errorf("Expected logger serviceName %s, got %s", serviceName, logger.ServiceName())
	}

	// Test BuildAndSetAsDefault
	defaultLogger := builder.BuildAndSetAsDefault()
	if defaultLogger == nil {
		t.Error("BuildAndSetAsDefault should return a non-nil Logger")
	}
}

// TestOtherLogLevels tests the log levels that aren't tested elsewhere
func TestOtherLogLevels(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Level:      TraceLevel, // Use lowest level to test all
		Pretty:     false,
		WithCaller: false,
		Output:     &buf,
	})

	// Test Trace
	log.Trace().Msg("trace message")
	assertLogContains(t, buf.String(), "trace message", "trace")
	buf.Reset()

	// Test Panic (use recover to prevent actual panic)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Panic level should have caused a panic")
		}
	}()

	log.Panic().Msg("panic message")
	assertLogContains(t, buf.String(), "panic message", "panic")
}

// TestMessageMethods tests the convenience message methods
func TestMessageMethods(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Level:      TraceLevel, // Use lowest level to test all
		Pretty:     false,
		WithCaller: false,
		Output:     &buf,
	})

	// DebugMsg
	log.DebugMsg("debug direct message")
	assertLogContains(t, buf.String(), "debug direct message", "debug")
	buf.Reset()

	// DebugMsg with formatting
	log.DebugMsg("debug %s message", "formatted")
	assertLogContains(t, buf.String(), "debug formatted message", "debug")
	buf.Reset()

	// WarnMsg
	log.WarnMsg("warn direct message")
	assertLogContains(t, buf.String(), "warn direct message", "warn")
	buf.Reset()

	// ErrorMsg
	log.ErrorMsg("error direct message")
	assertLogContains(t, buf.String(), "error direct message", "error")
	buf.Reset()

	// TraceMsg
	log.TraceMsg("trace direct message")
	assertLogContains(t, buf.String(), "trace direct message", "trace")
	buf.Reset()
}

// TestEnvironmentFunctions tests environment-related functions
func TestEnvironmentFunctions(t *testing.T) {
	// Test GetEnvStr with existing variable
	origEnv := os.Getenv("TEST_ENV_VAR")
	defer os.Setenv("TEST_ENV_VAR", origEnv)

	os.Setenv("TEST_ENV_VAR", "test-value")
	result := GetEnvStr("TEST_ENV_VAR", "default")
	if result != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", result)
	}

	// Test GetEnvStr with non-existing variable
	result = GetEnvStr("NON_EXISTENT_VAR", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}

	// Test GetEnvBool with true values
	trueValues := []string{"true", "TRUE", "True", "yes", "YES", "1"}
	for _, val := range trueValues {
		os.Setenv("TEST_BOOL_VAR", val)
		boolResult := GetEnvBool("TEST_BOOL_VAR", false)
		if !boolResult {
			t.Errorf("Expected true for '%s', got false", val)
		}
	}

	// Test GetEnvBool with false values
	falseValues := []string{"false", "FALSE", "no", "NO", "0", "anything-else"}
	for _, val := range falseValues {
		os.Setenv("TEST_BOOL_VAR", val)
		boolResult := GetEnvBool("TEST_BOOL_VAR", true)
		if boolResult {
			t.Errorf("Expected false for '%s', got true", val)
		}
	}

	// Test GetEnvBool with non-existing variable
	os.Unsetenv("TEST_BOOL_VAR")
	boolResult := GetEnvBool("TEST_BOOL_VAR", true)
	if !boolResult {
		t.Error("Expected default value true, got false")
	}

	// Test NewFromEnv
	os.Setenv(EnvLogLevel, "debug")
	os.Setenv(EnvLogFormat, "pretty")
	os.Setenv(EnvLogCaller, "true")
	os.Setenv(EnvServiceName, "env-service")

	logger := NewFromEnv()
	if logger == nil {
		t.Error("NewFromEnv should return a non-nil logger")
	}

	if logger.ServiceName() != "env-service" {
		t.Errorf("Expected service name 'env-service', got '%s'", logger.ServiceName())
	}

	// Test InitGlobalFromEnv
	globalLogger := InitGlobalFromEnv()
	if globalLogger == nil {
		t.Error("InitGlobalFromEnv should return a non-nil logger")
	}

	if globalLogger.ServiceName() != "env-service" {
		t.Errorf("Expected service name 'env-service', got '%s'", globalLogger.ServiceName())
	}
}

// TestLogLevelParsing tests the ParseLevel and String functions
func TestLogLevelParsing(t *testing.T) {
	// Test all valid level strings
	testCases := []struct {
		levelStr string
		level    Level
	}{
		{"trace", TraceLevel},
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"error", ErrorLevel},
		{"fatal", FatalLevel},
		{"panic", PanicLevel},
		// Test uppercase versions too
		{"TRACE", TraceLevel},
		{"DEBUG", DebugLevel},
		{"INFO", InfoLevel},
		{"WARN", WarnLevel},
		{"WARNING", WarnLevel},
		{"ERROR", ErrorLevel},
		{"FATAL", FatalLevel},
		{"PANIC", PanicLevel},
	}

	for _, tc := range testCases {
		level, err := ParseLevel(tc.levelStr)
		if err != nil {
			t.Errorf("ParseLevel(%s) returned error: %v", tc.levelStr, err)
		}
		if level != tc.level {
			t.Errorf("ParseLevel(%s) = %v, want %v", tc.levelStr, level, tc.level)
		}
	}

	// Test invalid level string
	_, err := ParseLevel("invalid")
	if err == nil {
		t.Error("ParseLevel(invalid) should return an error")
	}

	// Test String method for each level
	stringTests := []struct {
		level    Level
		expected string
	}{
		{TraceLevel, "trace"},
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{PanicLevel, "panic"},
	}

	for _, tc := range stringTests {
		str := tc.level.String()
		if str != tc.expected {
			t.Errorf("Level(%d).String() = %s, want %s", tc.level, str, tc.expected)
		}
	}
}

// TestWithMethod tests the With method for adding context
func TestWithMethod(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Level:      InfoLevel,
		Pretty:     false,
		WithCaller: false,
		Output:     &buf,
	})

	// Use With to create a context
	ctx := log.With()
	// Add some fields
	ctx = ctx.Str("context_key", "context_value")

	// Create a new logger from context
	contextLogger := ctx.Logger()

	// Write a log with the context logger
	contextLogger.Info().Msg("context message")

	// Verify context fields are present
	assertLogContains(t, buf.String(), "context_key", "")
	assertLogContains(t, buf.String(), "context_value", "")
}

// TestDefaultConfig tests the DefaultConfig function
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != InfoLevel {
		t.Errorf("Expected default level to be InfoLevel, got %v", cfg.Level)
	}

	if cfg.Pretty != false {
		t.Error("Expected default Pretty to be false")
	}

	if cfg.WithCaller != true {
		t.Error("Expected default WithCaller to be true")
	}

	if cfg.Output != os.Stderr {
		t.Error("Expected default Output to be os.Stderr")
	}
}

// TestOptionsAPI tests the functional options API
func TestOptionsAPI(t *testing.T) {
	// Test individual options
	levelOption := WithLevel(DebugLevel)
	if levelOption == nil {
		t.Error("WithLevel should return a non-nil option function")
	}

	prettyOption := WithPrettyPrint(true)
	if prettyOption == nil {
		t.Error("WithPrettyPrint should return a non-nil option function")
	}

	callerOption := WithCaller(false)
	if callerOption == nil {
		t.Error("WithCaller should return a non-nil option function")
	}

	outputOption := WithOutput(os.Stdout)
	if outputOption == nil {
		t.Error("WithOutput should return a non-nil option function")
	}

	timeFormatOption := WithTimeFormat(time.RFC822)
	if timeFormatOption == nil {
		t.Error("WithTimeFormat should return a non-nil option function")
	}

	// Test NewWithOptions
	logger := NewWithOptions(
		WithLevel(DebugLevel),
		WithPrettyPrint(true),
		WithCaller(false),
		WithOutput(os.Stdout),
		WithTimeFormat(time.RFC822),
	)

	if logger == nil {
		t.Error("NewWithOptions should return a non-nil logger")
	}

	// Test predefined configurations
	defaultLogger := Default()
	if defaultLogger == nil {
		t.Error("Default should return a non-nil logger")
	}

	prodLogger := Production()
	if prodLogger == nil {
		t.Error("Production should return a non-nil logger")
	}

	devLogger := Development()
	if devLogger == nil {
		t.Error("Development should return a non-nil logger")
	}
}

// TestInvalidLogLevel tests the behavior with invalid log levels
func TestInvalidLogLevel(t *testing.T) {
	// Test with an invalid log level string in NewFromEnv
	origLevel := os.Getenv(EnvLogLevel)
	defer os.Setenv(EnvLogLevel, origLevel)

	os.Setenv(EnvLogLevel, "invalid_level")
	logger := NewFromEnv()

	// Should default to InfoLevel when level is invalid
	if logger == nil {
		t.Error("NewFromEnv should return a non-nil logger even with invalid level")
	}
}

// TestFormatter tests the formatter functionality
func TestFormatter(t *testing.T) {
	// Test JSONFormatter
	jsonFormatter := DefaultJSONFormatter()
	if jsonFormatter == nil {
		t.Error("DefaultJSONFormatter should return a non-nil formatter")
	}

	// Test Format method
	var buf bytes.Buffer
	result := jsonFormatter.Format(&buf)

	// The result should be an io.Writer
	if result == nil {
		t.Error("JSONFormatter.Format should return a non-nil io.Writer")
	}

	// Test PrettyFormatter
	prettyFormatter := DefaultPrettyFormatter()
	if prettyFormatter == nil {
		t.Error("DefaultPrettyFormatter should return a non-nil formatter")
	}

	// Test Format method
	buf.Reset()
	result = prettyFormatter.Format(&buf)

	// The result should be an io.Writer
	if result == nil {
		t.Error("PrettyFormatter.Format should return a non-nil io.Writer")
	}
}

// TestAddField tests the generic AddField method
func TestAddField(t *testing.T) {
	var buf bytes.Buffer

	log := New(Config{
		Level:      InfoLevel,
		Pretty:     false,
		WithCaller: false,
		Output:     &buf,
	})

	// Test using AddField for generic values
	log.Info().
		AddField("string", "value").
		AddField("int", 42).
		AddField("bool", true).
		AddField("float", 3.14).
		Msg("message with generic fields")

	logData := buf.String()

	// Verify the fields are present
	assertLogContains(t, logData, "value", "")
	assertLogContains(t, logData, "42", "")
	assertLogContains(t, logData, "true", "")
	assertLogContains(t, logData, "3.14", "")
}
