package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Public logging methods without caller by default.
// Use AddCallerSkip(2) to point to the caller in the main application.

func Debug(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Debug(fmt.Sprint(message...))
}

func Info(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	// zapLog.zap.Sugar().WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Info(message...)
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Info(fmt.Sprint(message...))
}

func Warn(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Warn(fmt.Sprint(message...))
}

func Error(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Error(fmt.Sprint(message...))
}

func Fatal(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	Sync()
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Fatal(fmt.Sprint(message...))
}

// Panic explicitly logs the caller and stacktrace.
func Panic(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	Sync()
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2)).Panic(fmt.Sprint(message...))
}

// Public logging methods with traceID.
func DebugWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Debug(fmt.Sprint(message...), *fields...)
}

func InfoWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Info(fmt.Sprint(message...), *fields...)
}

func WarnWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Warn(fmt.Sprint(message...), *fields...)
}

func ErrorWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Error(fmt.Sprint(message...), *fields...)
}

func PanicWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	Sync()
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2)).Panic(fmt.Sprint(message...), *fields...)
}

func FatalWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	Sync()
	zapLog.zap.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Fatal(fmt.Sprint(message...), *fields...)
}
