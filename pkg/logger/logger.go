package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Global logger instance
var log = logrus.New()

// InitLogger initializes the logger settings such as format and log level
func InitLogger() {
	// Set the log format (using text formatter with full timestamps)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// Get the log level from the config file
	logLevel := viper.GetString("log_level")
	if logLevel == "" {
		logLevel = "debug" // Default to Debug level if not set
	}

	// Set the log level
	switch logLevel {
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.DebugLevel) // Default to Debug level
	}
}

// Info logs an informational message
func Info(msg string) {
	log.Info(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	log.Debug(msg)
}

// Error logs an error message
func Error(msg string) {
	log.Error(msg)
}
