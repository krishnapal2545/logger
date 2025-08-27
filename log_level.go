package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

// LogLevel represents the logging level.
type LogLevel int

const (
	TraceLevel LogLevel = iota - 1 // Zap's Debug is 0, so Trace is -1.
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func ParseLevel(level string) LogLevel {
	level = strings.ToUpper(level)
	switch level {
	case "TRACE":
		return TraceLevel
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
	case TraceLevel:
		return zapcore.Level(-1) // Custom trace level.
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
