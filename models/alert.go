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
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	AlertDate   time.Time     `bson:"alertDate,omitempty" json:"alertDate"`
	Channel     string        `bson:"channel,omitempty" json:"channel"`
	Method      string        `bson:"method,omitempty" json:"method,omitempty"`
	ContentType string        `bson:"contentType,omitempty" json:"contentType,omitempty"`
	URL         string        `bson:"url,omitempty" json:"url,omitempty"`
	Data        string        `bson:"data,omitempty" json:"data,omitempty"`
	Schedule    *Schedule     `bson:"schedule,omitempty" json:"schedule,omitempty"`
	Enabled     bool          `json:"enabled"`
	User        bson.ObjectId `json:"-"`
}

func NewAlert() *Alert {
	return &Alert{
		Enabled: true,
		Channel: "http",
		Method:  "GET",
	}
}
