package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	User     bson.ObjectId `json:"-" bson:"user,omitempty"`
	Expired  time.Time     `json:"-" bson:"expired"`
	FlashMsg *FlashMsg     `json:"-" bson:"flashMsg,omitempty"`
}

type FlashMsg struct {
	Msg    string `json:"-" bson:"msg,omitempty"`
	Status string `json:"-" bson:"status,omitempty"`
}
