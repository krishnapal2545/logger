package logger

import (
	"fmt"
	"os"
	"testing"
)

func TestLogging(t *testing.T){
	if err := Init(); err != nil {
		panic(err)
	}
	defer Recover()

	Debug("testing debug")
	DebugWithTraceID("debug","test")
	Info("testing info")
	InfoWithTraceID("info","test")
	Warn("testing warn")
	WarnWithTraceID("warn","test")
	Error("testing error")
	ErrorWithTraceID("error","testing")
	Panic("testing panic")
	Fatal("testing fatal")
	FatalWithTraceID("fatal", "testing")
}

func TestOneMillionLogging(t *testing.T) {
	if err := Init(); err != nil {
		panic(err)
	}
	defer Recover()

	for i := range 10_00_000 {
		InfoWithTraceID(fmt.Sprintf("trace-id-%v", i), "test message", i, "extra", "extra1", "extra2", 100, map[string]any{"key1": 12, "key2": "krishna"})
	}
}

// BenchmarkInfo tests the performance of the Info method (no traceID).
func BenchmarkInfo(b *testing.B) {
	// Setup logger.
	dir := "/logs/benchmark/logger"
	filename := "app"
	if err := os.MkdirAll(dir, 0755); err != nil {
		b.Fatal(err)
	}
	config := Config{
		Dir:             dir,
		Filename:        filename,
		FileMinLevel:    DebugLevel,
		ConsoleMinLevel: ErrorLevel,
	}
	if err := Init(config); err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir) // Clean up after test.

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		Info("test message", i, "extra")
	}
}

// BenchmarkInfoWithTraceID tests the performance of the InfoWithTraceID method.
func BenchmarkInfoWithTraceID(b *testing.B) {
	// Setup logger.
	dir := "/logs/benchmark/logger"
	filename := "app"
	if err := os.MkdirAll(dir, 0755); err != nil {
		b.Fatal(err)
	}
	config := Config{
		Dir:             dir,
		Filename:        filename,
		FileMinLevel:    DebugLevel,
		ConsoleMinLevel: ErrorLevel,
	}
	if err := Init(config); err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		InfoWithTraceID("abc123", "test message", i, "extra")
	}
}

// BenchmarkErrorWithTraceID tests the performance of the ErrorWithTraceID method.
func BenchmarkErrorWithTraceID(b *testing.B) {
	// Setup logger.
	dir := "D:\\logs\\benchmark\\logger"
	filename := "app"
	if err := os.MkdirAll(dir, 0755); err != nil {
		b.Fatal(err)
	}
	config := Config{
		Dir:             dir,
		Filename:        filename,
		FileMinLevel:    DebugLevel,
		ConsoleMinLevel: ErrorLevel,
	}
	if err := Init(config); err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		ErrorWithTraceID("abc123", "error message", i, "extra")
	}
}
