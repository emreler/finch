package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Schedule represents schedule of alert
type Schedule struct {
	RepeatEvery int
}

// Alert struct represents alert data stored in storage
type Alert struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name,omitempty"`
	AlertDate   time.Time     `json:"alertDate"`
	Channel     string        `json:"channel"`
	Method      string        `json:"method,omitempty"`
	ContentType string        `json:"contentType,omitempty"`
	URL         string        `json:"url,omitempty"`
	Data        string        `json:"data,omitempty"`
	Schedule    *Schedule     `json:"schedule,omitempty"`
	User        bson.ObjectId `json:"-"`
}
