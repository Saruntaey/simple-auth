package controllers

import (
	"net/http"
	"time"

	"github.com/saruntaey/simple-auth/models"
	"gopkg.in/mgo.v2/bson"
)

func (c *Controller) getSession(r *http.Request) (*models.Session, error) {
	session := &models.Session{}
	cookie, err := r.Cookie("session")
	if err != nil {
		return session, err
	}
	sid := cookie.Value
	query := bson.M{
		"_id": bson.ObjectIdHex(sid),
		"expired": bson.M{
			"$gt": time.Now(),
		},
	}

	err = c.appConfig.Session.Find(query).One(session)
	if err != nil {
		return session, err
	}
	return session, nil
}
