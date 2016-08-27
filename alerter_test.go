package main

import "testing"

func TestAsd(t *testing.T) {
	alerter := InitAlerter()
	alerter.AddAlert("foo", 10)
}
