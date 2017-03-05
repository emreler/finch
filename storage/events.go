package storage

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/emreler/finch/models"
)

const eventsCollection = "events"
const typeCreateAlert = "create_alert"
const typeCreateUser = "create_user"
const typeProcessAlert = "process_alert"

func (s *Storage) CountProcessAlertLogs() (int, error) {
	ses := s.GetDBSession()
	defer ses.Close()

	data := struct {
		Type string
	}{
		typeProcessAlert,
	}

	count, err := ses.DB("finch").C(eventsCollection).Find(data).Count()

	if err != nil {
		return 0, err
	}

	return count, nil
}

// LogProcessAlert creates process_alert event
func (s *Storage) LogProcessAlert(alert *models.Alert, statusCode int) error {
	ses := s.GetDBSession()
	defer ses.Close()

	data := struct {
		Type       string
		Alert      bson.ObjectId
		StatusCode int
		CreatedAt  time.Time
	}{
		typeProcessAlert,
		alert.ID,
		statusCode,
		time.Now(),
	}

	err := ses.DB("finch").C(eventsCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}

// LogCreateAlert creates create_alert event
func (s *Storage) LogCreateAlert(alert *models.Alert) error {
	ses := s.GetDBSession()
	defer ses.Close()

	data := struct {
		Type      string
		Alert     bson.ObjectId
		CreatedAt time.Time
	}{
		typeCreateAlert,
		alert.ID,
		time.Now(),
	}

	err := ses.DB("finch").C(eventsCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}

// LogCreateUser creates create_user event
func (s *Storage) LogCreateUser(user *models.User) error {
	ses := s.GetDBSession()
	defer ses.Close()

	data := struct {
		Type      string
		User      bson.ObjectId
		CreatedAt time.Time
	}{
		typeCreateUser,
		user.ID,
		time.Now(),
	}

	err := ses.DB("finch").C(eventsCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}
