
# Go Logger

logger is a high-performance, thread-safe logging library for Go, built on zap (go.uber.org/zap). It provides structured logging with trace IDs, file and console output, and ensures no log loss during normal exits, panics, or fatal errors. The logger initializes automatically via environment variables and supports panic recovery in main and goroutines.
Features

## Performance
#### Fast Logging: ~4µs per log (242,900 logs/s, based on BenchmarkInfo).
#### Thread-Safe: Uses safeBufferedWriter for concurrent logging.
#### Structured Logs: Supports trace IDs and arbitrary fields.
#### No Log Loss: Flushes logs on exits, panics, and fatal errors.
#### Panic Recovery: Recover for main, Go for goroutines, with stack traces.
#### Environment Variables: Config via LOGGER_DIR, LOGGER_FILENAME, etc.

### Installation
```bash
go get github.com/krishnapal2545/logger@latest
```

## Quick Start
The logger initializes automatically using environment variables or defaults (/logs/, app, DEBUG for file, ERROR for console). Use defer Recover() in main and Go for goroutines:
package main

```
func main() {
    config := logger.Config{
		Dir:             "/logs/testing/logger",
		Filename:        "app",
		FileMinLevel:    logger.DebugLevel,
		ConsoleMinLevel: logger.InfoLevel,
	}
	if err := logger.Init(config); err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
    defer logger.Recover()
    logger.Go(func() {
        db.SomeOperation() // Panic caught, stack trace logged
        logger.InfoWithTraceID("abc123", "DB query", "success", 200)
    })
    logger.Info("Starting app", 42)
    logger.FatalWithTraceID("xyz789", "Failed", 500)
}
```

### Configuration
```bash
Environment Variables
set LOGGER_DIR=D:\logs
set LOGGER_FILENAME=app
set LOGGER_FILE_MIN_LEVEL=debug
set LOGGER_CONSOLE_MIN_LEVEL=error
```

#### Environment Variables Description Default

`LOGGER_DIR`
Log file directory
/logs 

`LOGGER_FILENAME`
Log file base name
app

`LOGGER_FILE_MIN_LEVEL`
File log level (debug, info, warn, error)
debug

`LOGGER_CONSOLE_MIN_LEVEL`
Console log level
error

### Manual Configuration

```bash
    config := logger.Config{
        Dir:             `/custom/logs`,
        Filename:        "myapp",
        FileMinLevel:    logger.InfoLevel,
        ConsoleMinLevel: logger.ErrorLevel,
    }
    if err := logger.Init(config); err != nil {
        panic(err)
    }
```

### Usage Logging Methods

```bash
logger.Info(message ...any)
logger.InfoWithTraceID(traceID string, message ...any)
logger.Fatal(message ...any): Logs and exits.
logger.FatalWithTraceID(traceID string, message ...any)
Similar methods for Debug, Warn, Error, Panic.

Example log:
27/08/2025 19:27:43.925 | INFO | main.go:25 | TRACE : abc123 | Starting app 42 extra
27/08/2025 19:27:43.929 | ERROR | redis/test.go:14 | Goroutine panic testing goroutine panic | stack: goroutine 18 [running]: ... redis/test.go:14 ...
```

### Log Flushing

Normal Exit: Recover calls Sync.

Panics: Recover and Go flush logs and log stack traces.

Fatal Errors: Fatal methods flush logs.

Manual Flush: logger.Sync().

### Log File
Logs are written to /logs/app-YYYY-MM-DD-HH-MM-SS.log. 

Check: wc -l /logs/app-2025-08-27-22-19-34.log


### Troubleshooting

`Log Loss` Use Recover and Go for panics; 

`Fatal` methods auto-flush.

`Panic` Stack Trace Go logs stack traces for goroutine panics.

`File Errors` Check os.Stderr for write/sync errors.

`Concurrency` Run go test -race.


### Contributing
File issues or PRs at github.com/krishnapal2545/logger.

License
MIT License.