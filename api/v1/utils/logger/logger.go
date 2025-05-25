package logger

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

// Log levels
const (
	InfoLevel  = "INFO"
	ErrorLevel = "ERROR"
)

// log formats a message with optional fields
func logMessage(level string, msg string, fields map[string]interface{}) {
	if fields != nil {
		logger.Printf("[%s] %s | %+v\n", level, msg, fields)
	} else {
		logger.Printf("[%s] %s\n", level, msg)
	}
}

// Info logs an info message
func Info(msg string) {
	logMessage(InfoLevel, msg, nil)
}

// Error logs an error message
func Error(err error, msg string) {
	if err != nil {
		logMessage(ErrorLevel, msg+": "+err.Error(), nil)
	} else {
		logMessage(ErrorLevel, msg, nil)
	}
}

// LogWithFields logs a message with fields
func LogWithFields(level string, msg string, fields map[string]interface{}) {
	logMessage(level, msg, fields)
}
