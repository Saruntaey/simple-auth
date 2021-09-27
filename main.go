package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/saruntaey/simple-auth/config"
	"github.com/saruntaey/simple-auth/controllers"
)

func main() {
	// setup app config
	appConfig := config.New()

	// init controllers
	c := controllers.New(appConfig)

	// mount routes
	http.HandleFunc("/register", c.Register)
	http.HandleFunc("/login", c.Login)
	http.HandleFunc("/me", c.GetMe)
	http.HandleFunc("/update", c.Update)
	http.HandleFunc("/logout", c.Logout)

	// serve the app
	fmt.Println("Listening on port: ", appConfig.Port)
	log.Fatal(http.ListenAndServe(appConfig.Port, nil))
}
