package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ProcessAlert .
type ProcessAlert struct {
	Type       string        `json:"type"`
	Alert      bson.ObjectId `json:"-"`
	StatusCode int           `json:"responseStatusCode"`
	CreatedAt  time.Time     `json:"date"`
}
