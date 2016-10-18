package storage

import (
	"os"
	"testing"
	"time"

	"gitlab.com/emreler/finch/config"
	"gitlab.com/emreler/finch/models"
	"gopkg.in/mgo.v2/bson"
)

var s *Storage
var userID string
var alertID string

func TestMain(m *testing.M) {
	config := config.NewConfig("../config.json")

	s = NewStorage(config.Mongo)
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	user := &models.User{Name: "foo", Email: "bar@usefinch.co"}

	ID, err := s.CreateUser(user)

	if err != nil {
		t.Error(err)
		return
	}

	userID = ID
}

func TestCreateAlert(t *testing.T) {
	alert := &models.Alert{Name: "foo's alert", User: bson.ObjectIdHex(userID), AlertDate: time.Now().Add(10 * time.Second), Data: "somedata"}

	ID, err := s.CreateAlert(alert)

	if err != nil {
		t.Error(err)
		return
	}

	alertID = ID
}

func TestGetAlert(t *testing.T) {
	alert, err := s.GetAlert(alertID)

	if err != nil {
		t.Error(err)
		return
	}

	if alert.Data != "somedata" {
		t.Errorf("Alert data is wrong")
		return
	}
}

func TestGetUserAlerts(t *testing.T) {
	alerts, err := s.GetUserAlerts(userID)

	if err != nil {
		t.Error(err)
		return
	}

	if len(alerts) == 1 && alerts[0].Data == "somedata" {
		return
	}

	t.Errorf("Invalid user alerts data")
}
