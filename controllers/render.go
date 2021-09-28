package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/saruntaey/simple-auth/config"
)

func (c *Controller) render(w http.ResponseWriter, page string, data interface{}) {
	var tmpl *template.Template

	if c.appConfig.InProduction {
		t, ok := c.appConfig.Tmpls[page]
		if !ok {
			log.Fatal("No required page")
		}
		tmpl = t
	} else {
		fileName := fmt.Sprintf("page.%s.gohtml", page)
		tmpl = config.LoadTmpl(c.appConfig.ExPath, fileName)
	}

	tmpl.Execute(w, data)
}
