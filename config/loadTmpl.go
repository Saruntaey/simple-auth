package config

import (
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

func LoadTmpl(exPath string) map[string]*template.Template {
	tmplCache := map[string]*template.Template{}

	pagesDir, err := filepath.Glob(filepath.Join(exPath, "views/page*"))
	if err != nil {
		log.Fatal("Fail to load pages: ", err)
	}

	for _, pageDir := range pagesDir {
		fileName := filepath.Base(pageDir)
		name := strings.Split(fileName, ".")[1]

		tmpl := template.Must(template.ParseFiles(pageDir))
		tmpl, err = tmpl.ParseGlob(filepath.Join(exPath, "views/layout*"))
		if err != nil {
			log.Fatal("Fail to load layout: ", err)
		}

		tmplCache[name] = tmpl
	}
	return tmplCache
}
