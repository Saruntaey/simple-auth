package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saruntaey/simple-auth/config"
	"github.com/saruntaey/simple-auth/models"
)

type data struct {
	Title string
	Data  interface{}
}

type Controller struct {
	appConfig *config.Config
}

func New(appConfig *config.Config) *Controller {
	return &Controller{
		appConfig: appConfig,
	}
}

// @desc    Register user
// @route   < GET | POST > /register
// @access  Public
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := data{
			Title: "register",
		}

		c.render(w, "register", data)

	case http.MethodPost:
		User := c.appConfig.DbConn.Model("User")
		user := &models.User{}

		User.New(user)

		json.NewDecoder(r.Body).Decode(user)
		err := user.Save()
		if err != nil {
			fmt.Print(err)
		}
		w.Write([]byte("register"))

	default:
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
	}
}

// @desc    Login user
// @route   < GET | POST > /login
// @access  Public
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := data{
			Title: "login",
		}

		c.render(w, "login", data)

	case http.MethodPost:
		// User := c.appConfig.DbConn.Model("User")
		// user := &models.User{}

		// User.New(user)

		// json.NewDecoder(r.Body).Decode(user)
		// err := user.Save()
		// if err != nil {
		// 	fmt.Print(err)
		// }
		w.Write([]byte("login"))

	default:
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
	}
}

// @desc    Get current user
// @route   GET /me
// @access  Private
func (c *Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get me"))
}

// @desc    Update user
// @route   < GET | PUT > /update
// @access  Private
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update details"))
}

// @desc    Logout user
// @route   GET /logout
// @access  Private
func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logout"))
}
