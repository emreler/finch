package storage

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/emreler/finch/config"
	"github.com/emreler/finch/errors"
	"github.com/emreler/finch/models"
)

const eventsCollection = "events"
const typeCreateAlert = "create_alert"
const typeCreateUser = "create_user"
const typeProcessAlert = "process_alert"

// MongoStorage struct is used for storeing persistant data of alerts
type MongoStorage struct {
	session *mgo.Session
}

// NewStorage creates and returns new Storage instance
func NewStorage(url config.MongoConfig) *MongoStorage {
	ses, err := mgo.Dial(string(url))

	if err != nil {
		log.Fatal(err)
	}

	return &MongoStorage{session: ses}
}

// GetDBSession returns a new connection from the pool
func (s *MongoStorage) GetDBSession() *mgo.Session {
	return s.session.Copy()
}

// CreateUser creates new user
func (s *MongoStorage) CreateUser(user *models.User) error {
	ses := s.GetDBSession()
	defer ses.Close()

	user.ID = bson.NewObjectId()

	err := ses.DB("finch").C("users").Insert(user)

	if err != nil {
		return err
	}

	return nil
}

// CreateAlert adds new alert to storage
func (s *MongoStorage) CreateAlert(a *models.Alert) error {
	ses := s.GetDBSession()
	defer ses.Close()

	a.ID = bson.NewObjectId()
	err := ses.DB("finch").C("alerts").Insert(a)

	if err != nil {
		return err
	}

	return nil
}

// GetAlert Finds and returns alert data from storage
func (s *MongoStorage) GetAlert(alertID string) (*models.Alert, error) {
	ID := bson.ObjectIdHex(alertID)

	ses := s.GetDBSession()
	ses.SetSocketTimeout(time.Second * 10)
	ses.SetSyncTimeout(time.Second * 10)
	defer ses.Close()

	alert := &models.Alert{}
	err := ses.DB("finch").C("alerts").Find(bson.M{"_id": ID}).One(alert)

	if err != nil {
		// don't retry if it's a "not found" error
		if err.Error() != mgo.ErrNotFound.Error() {
			return nil, &errors.RetryProcessError{Msg: err.Error()}
		}

		return nil, err
	}

	return alert, nil
}

// UpdateAlert .
func (s *MongoStorage) UpdateAlert(alert *models.Alert) error {
	ses := s.GetDBSession()
	ses.SetSocketTimeout(time.Second * 10)
	ses.SetSyncTimeout(time.Second * 10)
	defer ses.Close()

	err := ses.DB("finch").C("alerts").Update(bson.M{"_id": alert.ID}, alert)

	if err != nil {
		return err
	}

	return nil
}

// GetUserAlerts .
func (s *MongoStorage) GetUserAlerts(userID string) ([]*models.Alert, error) {
	if match, _ := regexp.Match(`(?i)^[a-f\d]{24}$`, []byte(userID)); !match {
		return nil, fmt.Errorf("User ID '%s' is not a valid MongoDB ObjectID", userID)
	}

	ses := s.GetDBSession()
	defer ses.Close()

	ID := bson.ObjectIdHex(userID)

	var alerts []*models.Alert

	err := ses.DB("finch").C("alerts").Find(bson.M{"user": ID}).All(&alerts)

	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		alert.AlertDate = alert.AlertDate.UTC()
	}

	return alerts, nil
}

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
