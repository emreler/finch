package logger

import (
	"fmt"
	"os"
	"testing"
)

type SomeStruct struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

var l *Logger

func TestMain(m *testing.M) {
	l = NewLogger()
	res := m.Run()
	os.Exit(res)
}

func TestLog(t *testing.T) {
	l.Info("hello again")
	l.Info(&SomeStruct{"Emığre Kağçıyan", 5234})
	l.Error(fmt.Errorf("some error occured with just error"))
}
