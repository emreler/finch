package models

import "gopkg.in/mgo.v2/bson"

// User .
type User struct {
	ID    bson.ObjectId `bson:"_id"`
	Name  string
	Email string
}
