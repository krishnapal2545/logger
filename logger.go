package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config for the logger.
type Config struct {
	Dir             string   // Directory for log file
	Filename        string   // Log file name
	FileMinLevel    LogLevel // Min level for file logging (default: Debug)
	ConsoleMinLevel LogLevel // Min level for stdout logging (default: Info)
}

// Logger wraps the Zap logger.
type Logger struct {
	zap       *zap.Logger
	file      *os.File
	buf       *bytes.Buffer
	flush     chan struct{}
	fieldPool *sync.Pool
}

// Global logger instance.
var log *Logger

// Init initializes the global logger with the given config.
// Call this once, e.g., in main.
func Init(config Config) error {
	if config.FileMinLevel == 0 {
		config.FileMinLevel = DebugLevel
	}
	if config.ConsoleMinLevel == 0 {
		config.ConsoleMinLevel = InfoLevel
	}

	// Register custom encoder.
	if err := zap.RegisterEncoder("custom", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return &customEncoder{Encoder: zapcore.NewJSONEncoder(cfg)}, nil
	}); err != nil {
		return err
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = ""   // Handled in custom encoder.
	encCfg.CallerKey = "" // Handled in custom encoder.

	// File syncer with timestamped filename.
	timestamp := time.Now().Format("02-01-2006-15-04-05")
	filename := fmt.Sprintf("%s-%s.log", config.Filename, timestamp) // e.g., app-26-08-2025-15-04-05.log
	// File syncer.
	path := filepath.Join(config.Dir, filename)
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// Buffered writer for file.
	buf := bytes.NewBuffer(make([]byte, 0, 4096)) // 4KB buffer.
	flush := make(chan struct{}, 1)

	// Field pool for zap.Field slices.
	fieldPool := &sync.Pool{
		New: func() any {
			slice := make([]zapcore.Field, 0, 1)
			return &slice
		},
	}

	// Start flush goroutine.
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-flush:
				if buf.Len() > 0 {
					if _, err := file.Write(buf.Bytes()); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
					}
					buf.Reset()
				}
			case <-ticker.C:
				if buf.Len() > 0 {
					if _, err := file.Write(buf.Bytes()); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to write to log file: %v\n", err)
					}
					buf.Reset()
				}
			}
		}
	}()

	// Console syncer.
	consoleSyncer := zapcore.AddSync(os.Stdout)

	// Level enablers.
	fileEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= config.FileMinLevel.toZapLevel()
	})
	consoleEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= config.ConsoleMinLevel.toZapLevel()
	})

	// Cores.
	fileCore := zapcore.NewCore(
		&customEncoder{Encoder: zapcore.NewJSONEncoder(encCfg)},
		zapcore.AddSync(&bufferedWriter{file: file, buf: buf, flush: flush}),
		fileEnabler,
	)
	consoleCore := zapcore.NewCore(
		&customEncoder{Encoder: zapcore.NewJSONEncoder(encCfg)},
		consoleSyncer,
		consoleEnabler,
	)

	// Tee them.
	core := zapcore.NewTee(fileCore, consoleCore)

	// Build logger.
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	log = &Logger{zap: zapLogger, file: file, buf: buf, flush: flush, fieldPool: fieldPool}
	return nil
}

// bufferedWriter wraps a file with buffering.
type bufferedWriter struct {
	file  *os.File
	buf   *bytes.Buffer
	flush chan struct{}
}

func (w *bufferedWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if w.buf.Len()+n > 4096 { // Flush if buffer exceeds 4KB.
		select {
		case w.flush <- struct{}{}:
		default:
		}
	}
	return w.buf.Write(p)
}

func (w *bufferedWriter) Sync() error {
	select {
	case w.flush <- struct{}{}:
	default:
	}
	return nil
}
