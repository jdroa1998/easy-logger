# Easy Logger

A high-performance structured logging library for Go based on [zerolog](https://github.com/rs/zerolog), designed to provide a clean API, efficient logs, and a fluid interface for your applications.

## Version

Current version: v0.1.1

### Changelog

#### v0.1.1
- Added `SetServiceName` method to dynamically update service name after logger creation

#### v0.1.0
- Initial release
- Basic logging functionality
- Structured logging support
- Environment variable configuration
- Builder pattern for logger creation

## Features

- **High Performance**: Optimized for speed and low memory usage
- **Structured JSON Logs**: Ready for systems like ELK, Loki, or Graylog
- **Fluid Interface**: Clean, chainable API for improved developer experience
- **Multiple Log Styles**: Support for both structured and printf-style logs
- **Service Identification**: Automatic tagging of logs with service name
- **Environment Variable Configuration**: Simple setup in different environments
- **Pretty Format**: Human-readable output for development environments
- **Customizable Timestamp Format**: Flexible time formatting options
- **Context Support**: Add structured fields to all logs
- **Production & Development Modes**: Pre-configured settings for different environments

## Installation

```bash
go get -u github.com/jdroa1998/easy-logger
```

## Quick Start

```go
package main

import (
    "errors"
    "github.com/jdroa1998/easy-logger/logger"
)

func main() {
    // Create a logger with default configuration
    log := logger.New(logger.DefaultConfig())
    
    // Log a simple message
    log.InfoMsg("Application started")
    
    // Log with formatting
    log.Info().Msg("Server listening on port %d", 8080)
    
    // Log with structured fields
    log.Info().
        Str("service", "api").
        Int("port", 8080).
        Bool("active", true).
        Msg("Server configured")
        
    // Log errors
    err := errors.New("connection refused")
    if err != nil {
        log.Error().
            WithError(err).
            Str("operation", "startup").
            Msg("Failed to initialize service")
    }
}
```

## Example Output

A simple log in JSON format:
```json
{"level":"info","service":"UNKNOWN-SERVICE","time":"2023-04-01T12:34:56Z","message":"Application started"}
```

A structured log in JSON format:
```json
{"level":"info","service":"api","port":8080,"active":true,"time":"2023-04-01T12:34:56Z","message":"Server configured"}
```

A log in pretty format (for development):
```
12:34:56 INF Server configured service=api port=8080 active=true
```

An error log with additional context:
```json
{"level":"error","service":"UNKNOWN-SERVICE","error":"connection refused","operation":"startup","time":"2023-04-01T12:34:56Z","message":"Failed to initialize service"}
```

## Using the Builder Pattern

```go
package main

import (
    "os"
    "time"
    "github.com/jdroa1998/easy-logger/logger"
)

func main() {
    // Create a logger with the builder pattern
    log := logger.NewBuilder().
        WithLevel(logger.DebugLevel).
        WithPrettyPrint(true).
        WithServiceName("my-service").
        WithCaller(true).
        Build()
    
    // Use predefined configurations
    prodLog := logger.NewBuilder().
        Production().
        WithServiceName("api-production").
        Build()
    
    devLog := logger.NewBuilder().
        Development().
        WithServiceName("api-development").
        Build()
    
    // Custom configuration
    customLog := logger.NewBuilder().
        WithLevel(logger.InfoLevel).
        WithPrettyPrint(false).
        WithCaller(true).
        WithOutput(os.Stdout).
        WithTimeFormat(time.RFC3339).
        WithServiceName("custom-service").
        Build()
    
    // Use the loggers
    log.InfoMsg("Logger configured with builder")
    prodLog.InfoMsg("Production logger ready")
    devLog.Debug().Msg("Value: %v", 42)
    customLog.Info().Str("type", "custom").Msg("Custom logger initialized")
}
```

## Performance

Easy Logger is designed to be efficient with minimal overhead above the underlying zerolog library:

| Operation | Easy Logger | zerolog | Overhead |
|-----------|-------------|---------|----------|
| Simple Log | 150.9 ns/op, 32 B/op, 1 allocs/op | 103.8 ns/op, 0 B/op, 0 allocs/op | ~45% |
| Structured Log | 208.5 ns/op, 32 B/op, 1 allocs/op | 142.6 ns/op, 0 B/op, 0 allocs/op | ~46% |
| With Caller | 754.0 ns/op, 368 B/op, 5 allocs/op | 556.3 ns/op, 344 B/op, 3 allocs/op | ~35% |

The overhead is reasonable considering the additional convenience features provided.

## Environment Variables

Configure the logger easily with environment variables:

```go
package main

import "github.com/jdroa1998/easy-logger/logger"

func main() {
    // Initialize from environment variables
    log := logger.NewFromEnv()
    
    log.InfoMsg("Logger configured from environment")
}
```

Available environment variables:
- `LOG_LEVEL`: Log level (trace, debug, info, warn, error, fatal, panic)
- `LOG_FORMAT`: Log format (json, pretty)
- `LOG_CALLER`: Enable/disable caller information (true, false)
- `SERVICE_NAME`: Service name to add to all logs

## Logging Styles

Easy Logger supports multiple logging styles to fit different coding preferences:

### 1. Direct Style (Simple)
```go
log.InfoMsg("User authenticated")
log.ErrorMsg("Database error: %v", err)
```

### 2. Builder Style (Structured)
```go
log.Info().
    Str("user", "admin").
    Int("id", 1234).
    Bool("active", true).
    Msg("User logged in")
    
log.Error().
    WithError(err).
    Str("query", "SELECT * FROM users").
    Int("attempt", 3).
    Msg("Database query failed")
```

### 3. With Context Fields
```go
// Create a logger with predefined fields
reqLogger := log.WithFields(map[string]any{
    "request_id": "req-123",
    "ip":         "192.168.1.10",
    "user_agent": "Mozilla/5.0...",
})

// All logs from this logger will include those fields
reqLogger.InfoMsg("Request started")
reqLogger.Warn().Int("response_time_ms", 500).Msg("Slow response")
```

## API Reference

### Core Types

#### Logger
The main logger instance that provides logging capabilities.

```go
type Logger struct {
    // Internal fields not exposed
}

// Methods
func (l *Logger) SetServiceName(name string)  // Set service name after creation
func (l *Logger) ServiceName() string         // Get current service name
```

#### LogBuilder
Provides a fluid interface for creating structured logs.

```go
type LogBuilder struct {
    // Internal fields not exposed
}
```

#### Level
Represents log severity levels.

```go
type Level int8

const (
    TraceLevel Level = -1
    DebugLevel Level = 0
    InfoLevel  Level = 1
    WarnLevel  Level = 2
    ErrorLevel Level = 3
    FatalLevel Level = 4
    PanicLevel Level = 5
)
```

### Creating Loggers

- `New(cfg Config) *Logger`: Create a new logger with the given configuration
- `NewBuilder() *LoggerBuilder`: Start building a logger with the builder pattern
- `NewFromEnv() *Logger`: Create a logger configured from environment variables
- `Default() *Logger`: Create a logger with default settings
- `Development() *Logger`: Create a logger optimized for development
- `Production() *Logger`: Create a logger optimized for production

### Log Levels

Each log level has methods in two styles:
1. Builder style: `logger.Debug()`, `logger.Info()`, etc.
2. Direct style: `logger.DebugMsg()`, `logger.InfoMsg()`, etc.

Available levels:
- `Trace`: `logger.Trace()`, `logger.TraceMsg()`
- `Debug`: `logger.Debug()`, `logger.DebugMsg()`
- `Info`: `logger.Info()`, `logger.InfoMsg()`
- `Warn`: `logger.Warn()`, `logger.WarnMsg()`
- `Error`: `logger.Error()`, `logger.ErrorMsg()`
- `Fatal`: `logger.Fatal()`, `logger.FatalMsg()` (terminates the program)
- `Panic`: `logger.Panic()`, `logger.PanicMsg()` (causes a panic)

### LogBuilder Methods

- `Str(key string, value string) *LogBuilder`: Add a string field
- `Int(key string, value int) *LogBuilder`: Add an integer field
- `Bool(key string, value bool) *LogBuilder`: Add a boolean field
- `AddField(key string, value any) *LogBuilder`: Add a generic field
- `WithError(err error) *LogBuilder`: Add an error
- `Msg(msg string, values ...any)`: Finalize the log with a message

### Context and Fields

- `WithFields(fields map[string]any) *Logger`: Create a new logger with predefined fields
- `With() zerolog.Context`: Access the underlying zerolog context
- `ServiceName() string`: Get the current service name

### Configuration

#### Config Struct
```go
type Config struct {
    Level       Level      // Minimum level to log
    Pretty      bool       // Enable pretty (human-readable) output
    WithCaller  bool       // Include caller information
    Output      io.Writer  // Destination for logs
    TimeFormat  string     // Format for timestamps
    ServiceName string     // Name to identify service in logs
}
```

#### Builder Methods
- `WithLevel(level Level) *LoggerBuilder`: Set minimum log level
- `WithPrettyPrint(enabled bool) *LoggerBuilder`: Enable/disable pretty format
- `WithCaller(enabled bool) *LoggerBuilder`: Include caller information
- `WithOutput(output io.Writer) *LoggerBuilder`: Set output destination
- `WithTimeFormat(format string) *LoggerBuilder`: Set timestamp format
- `WithServiceName(name string) *LoggerBuilder`: Set service name
- `Development() *LoggerBuilder`: Configure for development environment
- `Production() *LoggerBuilder`: Configure for production environment
- `Build() *Logger`: Create logger with configured settings

### Other Utilities

- `ParseLevel(levelStr string) (Level, error)`: Parse level from string
- `Level.String() string`: Convert level to string
- `DefaultConfig() Config`: Get default configuration
- `DefaultJSONFormatter() Formatter`: Get default JSON formatter
- `DefaultPrettyFormatter() Formatter`: Get default pretty formatter

## Best Practices

### Service Identification
Always set a meaningful service name to identify the source of logs:

```go
// During logger creation
logger := logger.NewBuilder().
    WithServiceName("payment-service").
    Build()

// Or after logger creation
logger.SetServiceName("payment-service")
```

### Structured Logging
Prefer structured logging over simple messages when possible:

```go
// Good: structured and searchable
log.Info().
    Str("user", user.ID).
    Int("items", cart.Count).
    Float64("total", cart.Total).
    Msg("Order completed")

// Less Useful: information buried in message
log.InfoMsg("User %s completed order with %d items for $%.2f", user.ID, cart.Count, cart.Total)
```

### Consistent Field Names
Use consistent field names across your application:

```go
// Consistent naming convention
log.Info().
    Str("user_id", "123").
    Str("order_id", "456").
    Int("item_count", 5).
    Msg("Order processed")
```

### Log Levels
Use appropriate log levels:

- `Trace`: Extremely detailed information, useful for debugging complex sequence of events
- `Debug`: Detailed information, useful during development
- `Info`: General operational information about system behavior
- `Warn`: Potentially harmful situations that might require attention
- `Error`: Error conditions that should be investigated
- `Fatal`: Severe error conditions that cause termination
- `Panic`: Catastrophic failures that require immediate attention

### Context Propagation
Use context loggers to maintain request context:

```go
// Create a request-scoped logger
reqLogger := log.WithFields(map[string]any{
    "request_id": requestID,
    "user_id": userID,
})

// Pass it to other functions
processOrder(reqLogger, order)

// In other functions
func processOrder(log *logger.Logger, order Order) {
    log.Info().
        Str("order_id", order.ID).
        Msg("Processing order")
}
```

## Complete Examples

See the [examples directory](https://github.com/jdroa1998/easy-logger/tree/main/examples) for more detailed usage examples.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.