package controllers

import (
	"github.com/saruntaey/simple-auth/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (c *Controller) flashMsg(sessionId bson.ObjectId, status string, msg string) {
	q := c.appConfig.Session.FindId(sessionId)
	session := &models.Session{}
	q.One(&session)

	flashMsg := &models.FlashMsg{
		Status: status,
		Msg:    msg,
	}

	session.FlashMsg = flashMsg
	change := mgo.Change{
		Update:    session,
		ReturnNew: true,
	}
	// save session to DB
	q.Apply(change, session)
}
