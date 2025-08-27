package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// LogLevel represents the logging level.
type LogLevel int

const (
	DebugLevel LogLevel = iota - 1 // Zap's Info is 0, so Trace is -1.
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

func ParseLevel(level string) LogLevel {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// toZapLevel converts LogLevel to zapcore.Level.
func (l LogLevel) toZapLevel() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return  zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
