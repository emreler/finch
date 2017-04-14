package logger

import (
	"encoding/json"
	"log"
)

const (
	levelInfo  = "INFO"
	levelError = "ERROR"
)

// Logger .
type Logger struct{}

// NewLogger returns a new Logger.
func NewLogger() *Logger {
	return &Logger{}
}

// Info logs messages with INFO level. Parameter must be either string or JSON serializable structs.
func (l *Logger) Info(data interface{}) {
	if str, ok := data.(string); ok {
		log.Printf("level=%s message='%s'", levelInfo, str)
	} else {
		jstring, _ := json.Marshal(data)
		log.Printf("level=%s message='%s'", levelInfo, jstring)
	}
}

// Error logs error with ERROR level.
func (l *Logger) Error(err error) {
	log.Printf("level=%s message='%s'", levelError, err)
}
