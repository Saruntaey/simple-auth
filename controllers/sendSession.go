package controllers

import (
	"net/http"
)

func (c *Controller) sendSession(w http.ResponseWriter, msg string) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    msg,
		HttpOnly: true,
	}
	if c.appConfig.InProduction {
		cookie.Secure = true
	}
	http.SetCookie(w, &cookie)
}
