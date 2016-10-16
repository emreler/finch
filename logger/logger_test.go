package logger

import (
	"fmt"
	"os"
	"testing"

	"gitlab.com/emreler/finch/config"
)

type SomeStruct struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

var l *Logger

func TestMain(m *testing.M) {
	config := config.NewConfig("../config.json")
	l = NewLogger(config.Logentries)
	res := m.Run()
	l.conn.Close()
	os.Exit(res)
}

func TestLog(t *testing.T) {
	l.Info("hello again")
	l.Info(&SomeStruct{"Emığre Kağçıyan", 5234})
	l.Error(fmt.Errorf("some error occured with just error"))
}
