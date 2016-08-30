package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Storage struct is used for storeing persistant data of alerts
type Storage struct {
	Session *mgo.Session
}

// NewStorage creates and returns new Storage instance
func NewStorage(url MongoConfig) *Storage {
	ses, err := mgo.Dial(string(url))

	if err != nil {
		log.Fatal(err)
	}

	return &Storage{Session: ses}
}

// Alert struct represents alert data stored in storage
type Alert struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	AlertDate time.Time
	Channel   string
	URL       string
	Data      string
}

// AddAlert adds new alert to storage
func (s *Storage) AddAlert(a *Alert) string {
	a.ID = bson.NewObjectId()
	err := s.Session.DB("tmpmail-dev").C("alerts").Insert(a)

	if err != nil {
		log.Fatal(err)
	}

	return a.ID.Hex()
}

// GetAlert Finds and returns alert data from storage
func (s *Storage) GetAlert(alertID string) *Alert {
	ID := bson.ObjectIdHex(alertID)

	alert := new(Alert)
	err := s.Session.DB("tmpmail-dev").C("alerts").Find(bson.M{"_id": ID}).One(alert)

	if err != nil {
		log.Fatal(err)
	}

	return alert
}
