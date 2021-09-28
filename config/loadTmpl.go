package config

import (
	"html/template"
	"log"
	"path/filepath"
	"strings"
)

func LoadTmpls(exPath string) map[string]*template.Template {
	tmplCache := map[string]*template.Template{}

	pagesDir, err := filepath.Glob(filepath.Join(exPath, "views/page*"))
	if err != nil {
		log.Fatal("Fail to load pages: ", err)
	}

	for _, pageDir := range pagesDir {
		fileName := filepath.Base(pageDir)
		name := strings.Split(fileName, ".")[1]

		tmpl := LoadTmpl(exPath, fileName)

		tmplCache[name] = tmpl
	}
	return tmplCache
}

func LoadTmpl(exPath string, fileName string) *template.Template {
	pageDir := filepath.Join(exPath, "views", fileName)
	tmpl := template.Must(template.ParseFiles(pageDir))
	tmpl = template.Must(tmpl.ParseGlob(filepath.Join(exPath, "views/layout*")))
	return tmpl
}
