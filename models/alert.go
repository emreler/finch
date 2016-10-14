package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Schedule struct {
	RepeatEvery int
}

// Alert struct represents alert data stored in storage
type Alert struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string
	AlertDate   time.Time
	Channel     string
	Method      string
	ContentType string
	URL         string
	Data        string
	Schedule    *Schedule
	User        bson.ObjectId
}
