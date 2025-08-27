package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Public logging methods without traceID.
func Debug(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	zapLog.zap.Debug(fmt.Sprint(message...))
}
func Info(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	zapLog.zap.Info(fmt.Sprint(message...))
}
func Warn(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	zapLog.zap.Warn(fmt.Sprint(message...))
}
func Error(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	zapLog.zap.Error(fmt.Sprint(message...))
}
func Fatal(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	Sync()
	zapLog.zap.Fatal(fmt.Sprint(message...))
}
func Panic(message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	Sync()
	zapLog.zap.WithOptions(zap.AddCaller()).Panic(fmt.Sprint(message...))

	// zapLog.zap.Panic(fmt.Sprint(message...))
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
	zapLog.zap.Debug(fmt.Sprint(message...), *fields...)
}
func InfoWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	zapLog.zap.Info(fmt.Sprint(message...), *fields...)
}
func WarnWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	zapLog.zap.Warn(fmt.Sprint(message...), *fields...)
}
func ErrorWithTraceID(traceID string, message ...any) {
	if zapLog == nil || zapLog.zap == nil {
		panic("logger is not initialized...")
	}
	fields := zapLog.fieldPool.Get().(*[]zapcore.Field)
	defer zapLog.fieldPool.Put(fields)
	*fields = (*fields)[:0]
	*fields = append(*fields, zap.String("traceid", traceID))
	zapLog.zap.Error(fmt.Sprint(message...), *fields...)
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
	zapLog.zap.Panic(fmt.Sprint(message...), *fields...)
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
	zapLog.zap.Fatal(fmt.Sprint(message...), *fields...)
}
