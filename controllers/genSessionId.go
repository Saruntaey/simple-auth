package controllers

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/saruntaey/simple-auth/models"
	"gopkg.in/mgo.v2/bson"
)

func (c *Controller) genSessionId(userId bson.ObjectId) bson.ObjectId {
	exp, err := strconv.Atoi(os.Getenv("SESSION_EXPIRED"))
	if err != nil {
		log.Println("SESSION_EXPIRED should be number in minute: ", err)
	}
	session := models.Session{
		Id:      bson.NewObjectId(),
		Expired: time.Now().Add(time.Minute * time.Duration(exp)),
	}
	if len(userId) > 0 {
		session.User = userId
	}
	err = c.appConfig.Session.Insert(session)
	if err != nil {
		log.Print("Fail to create session: ", err)
	}
	return session.Id
}
