package utils

import (
	"log"
	"os"
)

// Logger struct to hold loggers for different log levels
type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
}

// NewLogger initializes and returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{
		InfoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		DebugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs informational messages
func (l *Logger) Info(message string) {
	l.InfoLogger.Println(message)
}

// Error logs error messages
func (l *Logger) Error(message string) {
	l.ErrorLogger.Println(message)
}

// Debug logs debug messages (only if debugging is enabled)
func (l *Logger) Debug(message string, debugEnabled bool) {
	if debugEnabled {
		l.DebugLogger.Println(message)
	}
}
