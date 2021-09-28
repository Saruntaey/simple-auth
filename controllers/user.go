package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/saruntaey/simple-auth/config"
	"github.com/saruntaey/simple-auth/models"
	"github.com/zebresel-com/mongodm"
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
		r.ParseForm()

		User := c.appConfig.DbConn.Model("User")
		user := &models.User{}

		User.New(user)

		// fill data
		user.Email = r.PostForm.Get("email")
		user.PasswordRaw = r.PostForm.Get("password")
		user.Name = r.PostForm.Get("name")

		if valid, issues := user.ValidateCreate(); !valid {
			msg := ""
			for _, v := range issues {
				msg += v.Error()
				msg += ", "
			}
			msg = strings.TrimRight(msg, ", ")
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		err := user.HashPassword()
		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		err = user.Save()
		if _, ok := err.(*mongodm.DuplicateError); ok {
			msg := fmt.Sprintf("the email %s is already taken", user.Email)
			http.Error(w, msg, http.StatusBadRequest)
			return
		} else if err != nil {
			fmt.Print(err)
			msg := "server error"
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		sessionId := c.genSessionId(user.Id)
		c.sendSession(w, sessionId.Hex())
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("registered"))

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
	session, err := c.getSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	User := c.appConfig.DbConn.Model("User")
	user := &models.User{}
	User.FindId(session.User).Exec(user)
	w.Write([]byte(user.Name))
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
