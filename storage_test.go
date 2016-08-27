package main

import (
	"log"
	"testing"
)

func TestNewFoo(t *testing.T) {
	x := &Alert{Name: "emre"}
	s := NewStorage("mongodb://robocop:6Hi3QhgfWfmM@ds013162.mlab.com:13162/tmpmail-dev")
	id := s.AddAlert(x)
	log.Println(id)

	a := s.GetAlert("57c1f71289c75aed0825c405")
	log.Println(a)
}
