package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// safeBufferedWriter ensures thread-safe buffer writes with no log loss.
type safeBufferedWriter struct {
	file  *os.File
	buf   *bytes.Buffer
	mu    *sync.Mutex
	flush chan struct{}
}

func NewFileAndSafeBufferedWriter(config *Config) (*safeBufferedWriter, error) {
	// File syncer with timestamped filename.
	timestamp := time.Now().Format("02-01-2006-15-04-05")
	filename := fmt.Sprintf("%s-%s.log", config.Filename, timestamp) // e.g., app-26-08-2025-15-04-05.log
	// File syncer.
	path := filepath.Join(config.Dir, filename)
	if err := os.MkdirAll(config.Dir, 0755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer := &safeBufferedWriter{file: file, buf: bytes.NewBuffer(make([]byte, 0, 4096)), flush: make(chan struct{}, 1), mu: &sync.Mutex{}}
	// Start flush goroutine.
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-writer.flush:
				writer.mu.Lock()
				if writer.buf.Len() > 0 {
					if _, err := writer.file.Write(writer.buf.Bytes()); err != nil {
						os.Stderr.WriteString("Failed to write to log file: " + err.Error() + "\n")
					}
					writer.buf.Reset()
				}
				writer.mu.Unlock()
			case <-ticker.C:
				writer.mu.Lock()
				if writer.buf.Len() > 0 {
					if _, err := writer.file.Write(writer.buf.Bytes()); err != nil {
						os.Stderr.WriteString("Failed to write to log file: " + err.Error() + "\n")
					}
					writer.buf.Reset()
				}
				writer.mu.Unlock()
			}
		}
	}()
	return writer, nil
}

func (w *safeBufferedWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	n = len(p)
	// If buffer would exceed 4KB, flush first.
	if w.buf.Len()+n > 4096 {
		if _, err := w.file.Write(w.buf.Bytes()); err != nil {
			os.Stderr.WriteString("Failed to write to log file: " + err.Error() + "\n")
		}
		w.buf.Reset()
		select {
		case w.flush <- struct{}{}:
		default:
		}
	}
	return w.buf.Write(p)
}

func (w *safeBufferedWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buf.Len() > 0 {
		if _, err := w.file.Write(w.buf.Bytes()); err != nil {
			os.Stderr.WriteString("Failed to write to log file: " + err.Error() + "\n")
		}
		w.buf.Reset()
	}
	if err := w.file.Sync(); err != nil {
		os.Stderr.WriteString("Failed to sync log file: " + err.Error() + "\n")
	}
	select {
	case w.flush <- struct{}{}:
	default:
	}
	return nil
}
