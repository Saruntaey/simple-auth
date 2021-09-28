package config

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/saruntaey/simple-auth/models"
	"github.com/zebresel-com/mongodm"
)

type Config struct {
	InProduction bool
	DbConn       *mongodm.Connection
	Port         string
	Tmpls        map[string]*template.Template
	ExPath       string
}

func New() *Config {
	// check runing mode
	var inProduction bool

	mode := os.Getenv("GO_ENV")
	if mode == "production" {
		inProduction = true
	} else if mode == "development" {
		inProduction = false
	} else {
		log.Fatal(`GO_ENV should be "production" or "development"`)
	}

	// get program path
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)

	// connect to DB
	cnn := connDB(exPath)

	// mouth models to DB
	cnn.Register(&models.User{}, "users")

	// load template to memory
	tmpls := map[string]*template.Template{}
	if inProduction {
		tmpls = LoadTmpls(exPath)
	}

	return &Config{
		InProduction: inProduction,
		DbConn:       cnn,
		Port:         fmt.Sprint(":", os.Getenv("PORT")),
		Tmpls:        tmpls,
		ExPath:       exPath,
	}
}
