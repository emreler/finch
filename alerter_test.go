package main

import "testing"

func TestAsd(t *testing.T) {
	alerter := NewAlerter()
	alerter.AddAlert("foo", 10)
}
