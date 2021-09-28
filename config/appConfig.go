package config

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/saruntaey/simple-auth/models"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2"
)

type Config struct {
	InProduction bool
	DbConn       *mongodm.Connection
	Session      *mgo.Collection
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

	// make field email in User model unique
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	err = cnn.Session.DB(os.Getenv("MONGO_DB")).C("users").EnsureIndex(index)
	if err != nil {
		log.Fatal("Fail to ensure index: ", err)
	}

	// mount session to DB
	session := cnn.Session.DB(os.Getenv("MONGO_DB")).C("sessions")

	// load template to memory
	tmpls := map[string]*template.Template{}
	if inProduction {
		tmpls = LoadTmpls(exPath)
	}

	return &Config{
		InProduction: inProduction,
		DbConn:       cnn,
		Session:      session,
		Port:         fmt.Sprint(":", os.Getenv("PORT")),
		Tmpls:        tmpls,
		ExPath:       exPath,
	}
}
