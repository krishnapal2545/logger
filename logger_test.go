package logger

import (
	"os"
	"testing"
)

// BenchmarkInfo tests the performance of the Info method (no traceID).
func BenchmarkInfo(b *testing.B) {
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
	defer os.RemoveAll(dir) // Clean up after test.

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		Info("test message", i, "extra")
	}
}

// BenchmarkInfoWithTraceID tests the performance of the InfoWithTraceID method.
func BenchmarkInfoWithTraceID(b *testing.B) {
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
