package logger

import (
	"os"
	"strconv"
)

// Config for the logger.
type Config struct {
	FileLogging     bool     // File logging (default: true)
	Dir             string   // Directory for log file (default: `/logs/app`)
	Filename        string   // Log file name (default: "app")
	FileMinLevel    LogLevel // Min level for file logging (default: Debug)
	ConsoleMinLevel LogLevel // Min level for stdout logging (default: Info)
}

func defaultConfig() Config {
	cfg := Config{
		Dir:             `/logs/app`,
		Filename:        "app",
		FileLogging:     true,
		FileMinLevel:    DebugLevel,
		ConsoleMinLevel: ErrorLevel,
	}

	if isFileLogging := os.Getenv("LOGGER_FILE_LOGGING"); isFileLogging != "" {
		cfg.FileLogging, _ = strconv.ParseBool(isFileLogging) 
	}

	if dir := os.Getenv("LOGGER_DIR"); dir != "" {
		cfg.Dir = dir
	}
	if filename := os.Getenv("LOGGER_FILENAME"); filename != "" {
		cfg.Filename = filename
	}
	if level := os.Getenv("LOGGER_FILE_MIN_LEVEL"); level != "" {
		cfg.FileMinLevel = ParseLevel(level)
	}
	if level := os.Getenv("LOGGER_CONSOLE_MIN_LEVEL"); level != "" {
		cfg.ConsoleMinLevel = ParseLevel(level)
	}

	return cfg
}
