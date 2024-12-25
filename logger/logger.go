package logger

import (
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var currentLevel LogLevel

func init() {
	level := os.Getenv("CONSUMER_LOG_LEVEL")
	switch strings.ToLower(level) {
	case "debug":
		currentLevel = DebugLevel
	case "info":
		currentLevel = InfoLevel
	case "warn":
		currentLevel = WarnLevel
	case "error":
		currentLevel = ErrorLevel
	default:
		currentLevel = WarnLevel // Default log level
	}
}

func Debug(msg string) {
	if currentLevel <= DebugLevel {
		log.Printf("[DEBUG] %s", msg)
	}
}

func Info(msg string) {
	if currentLevel <= InfoLevel {
		log.Printf("[INFO] %s", msg)
	}
}

// func Warn(msg string) {
// 	if currentLevel <= WarnLevel {
// 		log.Printf("[WARN] %s", msg)
// 	}
// }

func Error(msg string) {
	if currentLevel <= ErrorLevel {
		log.Printf("[ERROR] %s", msg)
	}
}
