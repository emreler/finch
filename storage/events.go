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

// GetAlertHistory returns event history for passed alertID
func (s *MongoStorage) GetAlertHistory(alertID string, limit int) ([]*models.ProcessAlert, error) {
	ses := s.GetDBSession()
	defer ses.Close()

	find := struct {
		Type  string
		Alert bson.ObjectId
	}{
		typeProcessAlert,
		bson.ObjectIdHex(alertID),
	}

	var res []*models.ProcessAlert

	query := ses.DB("finch").C(eventsCollection).Find(find)

	if limit != 0 && limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Sort("-createdat").All(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CountProcessAlertLogs counts total number of "process_alert" events
func (s *MongoStorage) CountProcessAlertLogs() (int, error) {
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
func (s *MongoStorage) LogProcessAlert(alert *models.Alert, statusCode int) error {
	ses := s.GetDBSession()
	defer ses.Close()

	data := &models.ProcessAlert{
		Type:       typeProcessAlert,
		Alert:      alert.ID,
		StatusCode: statusCode,
		CreatedAt:  time.Now(),
	}

	err := ses.DB("finch").C(eventsCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}

// LogCreateAlert creates create_alert event
func (s *MongoStorage) LogCreateAlert(alert *models.Alert) error {
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
func (s *MongoStorage) LogCreateUser(user *models.User) error {
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
