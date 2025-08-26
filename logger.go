package logger

import (
	"fmt"
	"os"
	"path/filepath"

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
	zap *zap.Logger
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

	// File syncer.
	path := filepath.Join(config.Dir, config.Filename)
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fileSyncer := zapcore.AddSync(file)

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
		fileSyncer,
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
	log = &Logger{zap: zapLogger}
	return nil
}

// Public logging methods without traceID.
func Trace(message ...any) { log.zap.Debug(fmt.Sprint(message...)) }
func Debug(message ...any) { log.zap.Debug(fmt.Sprint(message...)) }
func Info(message ...any)  { log.zap.Info(fmt.Sprint(message...)) }
func Warn(message ...any)  { log.zap.Warn(fmt.Sprint(message...)) }
func Error(message ...any) { log.zap.Error(fmt.Sprint(message...)) }
func Fatal(message ...any) { log.zap.Fatal(fmt.Sprint(message...)) }

// Public logging methods with traceID.
func TraceWithTraceID(traceID string, message ...any) {
	log.zap.Debug(fmt.Sprint(message...), zap.String("traceid", traceID))
}
func DebugWithTraceID(traceID string, message ...any) {
	log.zap.Debug(fmt.Sprint(message...), zap.String("traceid", traceID))
}
func InfoWithTraceID(traceID string, message ...any) {
	log.zap.Info(fmt.Sprint(message...), zap.String("traceid", traceID))
}
func WarnWithTraceID(traceID string, message ...any) {
	log.zap.Warn(fmt.Sprint(message...), zap.String("traceid", traceID))
}
func ErrorWithTraceID(traceID string, message ...any) {
	log.zap.Error(fmt.Sprint(message...), zap.String("traceid", traceID))
}
func FatalWithTraceID(traceID string, message ...any) {
	log.zap.Fatal(fmt.Sprint(message...), zap.String("traceid", traceID))
}
