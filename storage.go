package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Storage struct {
	Session *mgo.Session
}

func NewStorage(url MongoConfig) *Storage {
	ses, err := mgo.Dial(string(url))

	if err != nil {
		panic(err)
	}

	return &Storage{Session: ses}
}

type Alert struct {
	ID        bson.ObjectId `bson:"_id"`
	Name      string
	AlertDate time.Time
	Channel   string
	URL       string
	Data      string
}

func (s *Storage) AddAlert(a *Alert) string {
	a.ID = bson.NewObjectId()
	err := s.Session.DB("tmpmail-dev").C("alerts").Insert(a)

	if err != nil {
		log.Fatal(err)
	}

	return a.ID.Hex()
}

func (s *Storage) GetAlert(alertID string) *Alert {
	ID := bson.ObjectIdHex(alertID)

	alert := new(Alert)
	err := s.Session.DB("tmpmail-dev").C("alerts").Find(bson.M{"_id": ID}).One(alert)

	if err != nil {
		log.Fatal(err)
	}

	return alert
}
