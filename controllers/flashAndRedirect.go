package controllers

import (
	"net/http"

	"github.com/saruntaey/simple-auth/models"
)

func (c *Controller) flashAndRedirect(w http.ResponseWriter, r *http.Request, errGetSession error, session *models.Session, msg string, add string) {
	if errGetSession == nil {
		c.flashMsg(session.Id, "danger", msg)
	} else {
		sessionId := c.genSessionId("")
		c.flashMsg(sessionId, "danger", msg)
		c.sendSession(w, sessionId.Hex())
	}
	http.Redirect(w, r, add, http.StatusSeeOther)
}
