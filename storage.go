package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/emreler/finch/config"
	"gitlab.com/emreler/finch/models"
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

// User .
type User struct {
	ID    bson.ObjectId `bson:"_id"`
	Name  string
	Email string
	Token string
}

func generateToken() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

// CheckToken checks token
func (s *Storage) CheckToken(token string) (*bson.ObjectId, error) {
	user := &User{}
	err := s.Session.DB("tmpmail-dev").C("users").Find(bson.M{"token": token}).One(user)

	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}

// CreateUser creates new user
func (s *Storage) CreateUser(user *User) (string, error) {
	user.ID = bson.NewObjectId()
	user.Token = generateToken()

	err := s.Session.DB("tmpmail-dev").C("users").Insert(user)

	return user.Token, err
}

// CreateAlert adds new alert to storage
func (s *Storage) CreateAlert(a *models.Alert) string {
	a.ID = bson.NewObjectId()
	err := s.Session.DB("tmpmail-dev").C("alerts").Insert(a)

	if err != nil {
		log.Fatal(err)
	}

	return a.ID.Hex()
}

// GetAlert Finds and returns alert data from storage
func (s *Storage) GetAlert(alertID string) *models.Alert {
	ID := bson.ObjectIdHex(alertID)

	alert := &models.Alert{}
	err := s.Session.DB("tmpmail-dev").C("alerts").Find(bson.M{"_id": ID}).One(alert)

	if err != nil {
		log.Fatal(err)
	}

	return alert
}
