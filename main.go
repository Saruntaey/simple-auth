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

	// handle favicon
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// redirect root
	http.HandleFunc("/", redirectRoot)

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

func redirectRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}
