package logger

import (
	"bytes"
	"errors"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
func Init(configs ...Config) error {
	var config Config
	if len(configs) > 0 {
		config = configs[0]
		if config.FileMinLevel == 0 {
			config.FileMinLevel = DebugLevel
		}
		if config.ConsoleMinLevel == 0 {
			config.ConsoleMinLevel = InfoLevel
		}
	} else {
		config = defaultConfig()
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

	writer, err := NewFileAndSafeBufferedWriter(&config)
	if err != nil {
		return err
	}

	// Field pool for zap.Field slices.
	fieldPool := &sync.Pool{
		New: func() any {
			slice := make([]zapcore.Field, 0, 1)
			return &slice
		},
	}
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
		zapcore.AddSync(writer),
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
	log = &Logger{zap: zapLogger, file: writer.file, buf: writer.buf, flush: writer.flush, fieldPool: fieldPool}
	return nil
}

func Sync() error {
	if log == nil || log.zap == nil {
		return errors.New("logger not initialized")
	}
	Info("logger sync is being called ...")
	return log.zap.Sync()
}

func Recover() {
	if r := recover(); r != nil {
		Panic(r)
	}
	Sync()
}
