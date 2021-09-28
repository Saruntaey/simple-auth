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
	DbConn *mongodm.Connection
	Port   string
	Tmpl   map[string]*template.Template
	ExPath string
}

func New() *Config {
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
	tmpl := LoadTmpl(exPath)

	return &Config{
		DbConn: cnn,
		Port:   fmt.Sprint(":", os.Getenv("PORT")),
		Tmpl:   tmpl,
		ExPath: exPath,
	}
}
