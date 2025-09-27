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
	dir    string
	filename   string
	path  string
	buf   *bytes.Buffer
	mu    *sync.Mutex
	flush chan struct{}
}

const (
	maxflushbufferlength = 64 * 1024 // 64KB buffer
)

// createLogFile creates a new log file with timestamp in given dir/name.
func createLogFile(dir, baseName string) (*os.File, string, error) {
	timestamp := time.Now().Format("02-01-2006-15-04-05")
	filename := fmt.Sprintf("%s-%s.log", baseName, timestamp)
	path := filepath.Join(dir, filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, "", err
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, "", err
	}
	return file, path, nil
}

func NewFileAndSafeBufferedWriter(config *Config) (*safeBufferedWriter, error) {
	file, path, err := createLogFile(config.Dir, config.Filename)
	if err != nil {
		return nil, err
	}

	writer := &safeBufferedWriter{
		file: file,
		dir:  config.Dir,
		filename: config.Filename,
		path: path,
		buf:  bytes.NewBuffer(make([]byte, 0, maxflushbufferlength)),
		flush: make(chan struct{}, 1),
		mu:   &sync.Mutex{},
	}
	// Start flush goroutine.
	go writer.run()
	return writer, nil
}

func (w *safeBufferedWriter) run() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-w.flush:
			w.Flush()
		case <-ticker.C:
			w.Flush()
		}
	}
}

func (w *safeBufferedWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.buf.Len()+len(p) > maxflushbufferlength {
		select {
		case w.flush <- struct{}{}:
		default:
		}
	}
	return w.buf.Write(p)
}

func (w *safeBufferedWriter) Sync() error {
	w.Flush()
	return w.file.Sync()
}

func (w *safeBufferedWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buf.Len() == 0 {
		return
	}

	if _, err := os.Stat(w.path); os.IsNotExist(err) {
		newFile, newPath, err := createLogFile(w.dir, w.filename)
		if err != nil {
			os.Stderr.WriteString("Failed to recreate log file: " + err.Error() + "\n")
			return
		}
		w.file.Close()
		w.file = newFile
		w.path = newPath
	}

	if _, err := w.file.Write(w.buf.Bytes()); err != nil {
		os.Stderr.WriteString("Failed to write to log file: " + err.Error() + "\n")
	}
	w.buf.Reset()
}

func (w *safeBufferedWriter) Close() error {
	w.Sync()
	return w.file.Close()
}
