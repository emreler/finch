package storage

import (
	"fmt"
	"log"
	"regexp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/emreler/finch/config"
	"github.com/emreler/finch/models"
)

// Storage struct is used for storeing persistant data of alerts
type Storage struct {
	Session *mgo.Session
}

// NewStorage creates and returns new Storage instance
func NewStorage(url config.MongoConfig) *Storage {
	ses, err := mgo.Dial(string(url))

	if err != nil {
		log.Fatal(err)
	}

	return &Storage{Session: ses}
}

// GetDBSession returns a new connection from the pool
func (s *Storage) GetDBSession() *mgo.Session {
	return s.Session.Copy()
}

// CreateUser creates new user
func (s *Storage) CreateUser(user *models.User) (string, error) {
	ses := s.GetDBSession()
	defer ses.Close()

	user.ID = bson.NewObjectId()

	err := ses.DB("tmpmail-dev").C("users").Insert(user)

	if err != nil {
		return "", err
	}

	return user.ID.Hex(), nil
}

// CreateAlert adds new alert to storage
func (s *Storage) CreateAlert(a *models.Alert) (string, error) {
	ses := s.GetDBSession()
	defer ses.Close()

	a.ID = bson.NewObjectId()
	err := ses.DB("tmpmail-dev").C("alerts").Insert(a)

	if err != nil {
		return "", err
	}

	return a.ID.Hex(), nil
}

// GetAlert Finds and returns alert data from storage
func (s *Storage) GetAlert(alertID string) (*models.Alert, error) {
	if match, _ := regexp.Match(`(?i)^[a-f\d]{24}$`, []byte(alertID)); !match {
		return nil, fmt.Errorf("Alert ID '%s' is not a valid MongoDB ObjectID", alertID)
	}

	ID := bson.ObjectIdHex(alertID)

	ses := s.GetDBSession()
	defer ses.Close()

	alert := &models.Alert{}
	err := ses.DB("tmpmail-dev").C("alerts").Find(bson.M{"_id": ID}).One(alert)

	if err != nil {
		return nil, err
	}

	return alert, nil
}

// UpdateAlert .
func (s *Storage) UpdateAlert(alert *models.Alert) error {
	ses := s.GetDBSession()
	defer ses.Close()

	err := ses.DB("tmpmail-dev").C("alerts").Update(bson.M{"_id": alert.ID}, alert)

	if err != nil {
		return err
	}

	return nil
}

// GetUserAlerts .
func (s *Storage) GetUserAlerts(userID string) ([]*models.Alert, error) {
	if match, _ := regexp.Match(`(?i)^[a-f\d]{24}$`, []byte(userID)); !match {
		return nil, fmt.Errorf("User ID '%s' is not a valid MongoDB ObjectID", userID)
	}

	ses := s.GetDBSession()
	defer ses.Close()

	ID := bson.ObjectIdHex(userID)

	var alerts []*models.Alert

	err := ses.DB("tmpmail-dev").C("alerts").Find(bson.M{"user": ID}).All(&alerts)

	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		alert.AlertDate = alert.AlertDate.UTC()
	}

	return alerts, nil
}
