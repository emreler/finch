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

// LogMessage .
type LogMessage struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// NewLogger .
func NewLogger() *Logger {
	return &Logger{}
}

// Info .
func (l *Logger) Info(data interface{}) {
	var j []byte
	if str, ok := data.(string); ok {
		logMsg := &LogMessage{
			Level:   levelInfo,
			Message: str,
		}

		j, _ = json.Marshal(logMsg)

		log.Println(string(j))
	} else {
		jstring, _ := json.Marshal(data)

		logMsg := &LogMessage{
			Level:   levelInfo,
			Message: string(jstring),
		}

		j, _ = json.Marshal(logMsg)

		log.Println(string(j))
	}
}

func (l *Logger) Error(err error) {
	logMsg := &LogMessage{
		Level:   levelError,
		Message: err.Error(),
	}

	j, _ := json.Marshal(logMsg)

	log.Println(string(j))
}
