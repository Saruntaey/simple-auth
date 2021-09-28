package controllers

import (
	"log"
	"net/http"
)

func (c *Controller) render(w http.ResponseWriter, page string, data interface{}) {
	tmpl, ok := c.appConfig.Tmpl[page]
	if !ok {
		log.Fatal("No required page")
	}

	tmpl.Execute(w, data)
}
