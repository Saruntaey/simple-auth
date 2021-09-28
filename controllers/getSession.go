package controllers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/saruntaey/simple-auth/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (c *Controller) getSession(w http.ResponseWriter, r *http.Request) (*models.Session, *models.FlashMsg, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, nil, err
	}
	sid := cookie.Value
	query := bson.M{
		"_id": bson.ObjectIdHex(sid),
	}

	session := &models.Session{}
	q := c.appConfig.Session.Find(query)
	err = q.One(session)
	if err != nil {
		c.delCookie(w)
		return nil, nil, err
	}
	// check if the session expired
	if time.Now().After(session.Expired) {
		// delete session from DB
		c.appConfig.Session.RemoveId(session.Id)
		c.delCookie(w)
		return nil, nil, errors.New("session is expired")
	}

	// extend session expiration and delete flash message
	exp, err := strconv.Atoi(os.Getenv("SESSION_EXPIRED"))
	if err != nil {
		log.Println("SESSION_EXPIRED should be number in minute: ", err)
	}
	session.Expired = time.Now().Add(time.Minute * time.Duration(exp))
	flashMsg := session.FlashMsg
	session.FlashMsg = nil
	change := mgo.Change{
		Update:    session,
		ReturnNew: true,
	}
	// save session to DB
	q.Apply(change, session)
	return session, flashMsg, nil
}

func (c *Controller) delCookie(w http.ResponseWriter) {
	delCookie := http.Cookie{
		Name:     "session",
		Value:    "none",
		HttpOnly: true,
		MaxAge:   -1,
	}
	if c.appConfig.InProduction {
		delCookie.Secure = true
	}

	http.SetCookie(w, &delCookie)
}
