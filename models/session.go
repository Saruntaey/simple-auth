package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	User    bson.ObjectId `json:"-" bson:"user"`
	Expired time.Time     `json:"-" bson:"expired"`
}
