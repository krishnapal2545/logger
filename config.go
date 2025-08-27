package logger

import "os"

// Config for the logger.
type Config struct {
	Dir             string   // Directory for log file
	Filename        string   // Log file name
	FileMinLevel    LogLevel // Min level for file logging (default: Debug)
	ConsoleMinLevel LogLevel // Min level for stdout logging (default: Info)
}

func defaultConfig() Config {
	cfg := Config{
		Dir:             `/logs/app`,
		Filename:        "app",
		FileMinLevel:    DebugLevel,
		ConsoleMinLevel: ErrorLevel,
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
