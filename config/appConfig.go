package config

import (
	"fmt"
	"os"

	"github.com/zebresel-com/mongodm"
)

type Config struct {
	DbConn *mongodm.Connection
	Port   string
}

func New() *Config {
	cnn := connDB()

	return &Config{
		DbConn: cnn,
		Port:   fmt.Sprint(":", os.Getenv("PORT")),
	}
}
