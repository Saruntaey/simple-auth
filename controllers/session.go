package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/saruntaey/simple-auth/models"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	InProduction bool
	DbConn       *mongodm.Connection
	Session      *mgo.Collection
	SessionExp   int
	SessionModel *models.Session
	FlashMsg     *models.FlashMsg
}

func (c *Controller) NewSession() *Session {
	return &Session{
		InProduction: c.appConfig.InProduction,
		DbConn:       c.appConfig.DbConn,
		Session:      c.appConfig.Session,
		SessionExp:   c.appConfig.SessionExp,
	}
}

func (s *Session) InitModel() *Session {
	s.SessionModel = &models.Session{
		Id: bson.NewObjectId(),
	}
	return s
}

func (s *Session) GetFromCookie(w http.ResponseWriter, r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return s, err
	}
	sid := cookie.Value
	query := bson.M{
		"_id": bson.ObjectIdHex(sid),
	}

	session := &models.Session{}
	err = s.Session.Find(query).One(session)
	if err != nil {
		s.delCookie(w)
		return s, err
	}
	// check if the session expired
	if time.Now().After(session.Expired) {
		// delete session from DB
		s.Session.RemoveId(session.Id)
		s.delCookie(w)
		return s, errors.New("session is expired")
	}

	// extend session expiration and delete flash message in DB
	s.SessionModel = session
	s.FlashMsg = session.FlashMsg
	s.SessionModel.FlashMsg = nil
	s.Save()
	return s, nil
}

func (s *Session) delCookie(w http.ResponseWriter) {
	delCookie := http.Cookie{
		Name:     "session",
		Value:    "none",
		HttpOnly: true,
		MaxAge:   -1,
	}
	if s.InProduction {
		delCookie.Secure = true
	}

	http.SetCookie(w, &delCookie)
}

func (s *Session) AddUser(userId bson.ObjectId) *Session {
	s.SessionModel.User = userId
	return s
}

func (s *Session) AddFlashMsg(status string, msg string) *Session {
	s.SessionModel.FlashMsg = &models.FlashMsg{
		Status: status,
		Msg:    msg,
	}
	return s
}

func (s *Session) Save() (*models.Session, error) {
	s.SessionModel.Expired = time.Now().Add(time.Minute * time.Duration(s.SessionExp))
	query := bson.M{
		"_id": s.SessionModel.Id,
	}
	_, err := s.Session.Upsert(query, s.SessionModel)
	if err != nil {
		return nil, fmt.Errorf("fail to save session to DB: %w", err)
	}
	return s.SessionModel, err
}

func (s *Session) FlashAndRedirect(w http.ResponseWriter, r *http.Request, status string, msg string, address string) {
	if s.SessionModel == nil {
		s.InitModel()
	}
	s.AddFlashMsg(status, msg)
	s.Save()
	s.SendSession(w)
	http.Redirect(w, r, address, http.StatusSeeOther)
}

func (s *Session) SendSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    s.SessionModel.Id.Hex(),
		HttpOnly: true,
	}
	if s.InProduction {
		cookie.Secure = true
	}
	http.SetCookie(w, &cookie)
}
