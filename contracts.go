package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Public logging methods without traceID.
func Trace(message ...any) { log.zap.Debug(fmt.Sprint(message...)) }
func Debug(message ...any) { log.zap.Debug(fmt.Sprint(message...)) }
func Info(message ...any)  { log.zap.Info(fmt.Sprint(message...)) }
func Warn(message ...any)  { log.zap.Warn(fmt.Sprint(message...)) }
func Error(message ...any) { log.zap.Error(fmt.Sprint(message...)) }
func Fatal(message ...any) { log.zap.Fatal(fmt.Sprint(message...)) }

// Public logging methods with traceID.
func TraceWithTraceID(traceID string, message ...any) {
	fields := log.fieldPool.Get().([]zapcore.Field)
	defer log.fieldPool.Put(&fields)
	fields = fields[:0]
	fields = append(fields, zap.String("traceid", traceID))
	log.zap.Debug(fmt.Sprint(message...), fields...)
}
func DebugWithTraceID(traceID string, message ...any) {
	fields := log.fieldPool.Get().([]zapcore.Field)
	defer log.fieldPool.Put(&fields)
	fields = fields[:0]
	fields = append(fields, zap.String("traceid", traceID))
	log.zap.Debug(fmt.Sprint(message...), fields...)
}
func InfoWithTraceID(traceID string, message ...any) {
	fields := log.fieldPool.Get().([]zapcore.Field)
	defer log.fieldPool.Put(&fields)
	fields = fields[:0]
	fields = append(fields, zap.String("traceid", traceID))
	log.zap.Info(fmt.Sprint(message...), fields...)
}
func WarnWithTraceID(traceID string, message ...any) {
	fields := log.fieldPool.Get().([]zapcore.Field)
	defer log.fieldPool.Put(&fields)
	fields = fields[:0]
	fields = append(fields, zap.String("traceid", traceID))
	log.zap.Warn(fmt.Sprint(message...), fields...)
}
func ErrorWithTraceID(traceID string, message ...any) {
	fields := log.fieldPool.Get().([]zapcore.Field)
	defer log.fieldPool.Put(&fields)
	fields = fields[:0]
	fields = append(fields, zap.String("traceid", traceID))
	log.zap.Error(fmt.Sprint(message...), fields...)
}
func FatalWithTraceID(traceID string, message ...any) {
	fields := log.fieldPool.Get().([]zapcore.Field)
	defer log.fieldPool.Put(&fields)
	fields = fields[:0]
	fields = append(fields, zap.String("traceid", traceID))
	log.zap.Fatal(fmt.Sprint(message...), fields...)
}
