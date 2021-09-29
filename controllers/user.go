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
	Title    string
	FlashMsg *models.FlashMsg
	Data     interface{}
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
	session, errGetSession := c.NewSession().GetFromCookie(w, r)
	switch r.Method {

	case http.MethodGet:
		// check if user already login
		if errGetSession == nil && len(session.SessionModel.User) > 0 {
			http.Redirect(w, r, "/me", http.StatusSeeOther)
			return
		}

		data := data{
			Title: "register",
		}
		// add flash message to data
		if errGetSession == nil {
			data.FlashMsg = session.FlashMsg
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
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return

		}

		err := user.HashPassword()
		if err != nil {
			msg := err.Error()
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return
		}

		err = user.Save()
		if _, ok := err.(*mongodm.DuplicateError); ok {
			msg := fmt.Sprintf("the email %s is already taken", user.Email)
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return
		} else if err != nil {
			fmt.Print(err)
			msg := "server error"
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return
		}

		if errGetSession != nil {
			session.InitModel()
		}
		session.AddUser(user.Id)
		session.FlashAndRedirect(w, r, "success", "Welcome to Simple-Auth", "/me")

	default:
		session.FlashAndRedirect(w, r, "danger", "Method not allow", "/register")
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
	session, errGetSession := c.NewSession().GetFromCookie(w, r)
	if errGetSession != nil || len(session.SessionModel.User) == 0 {
		session.FlashAndRedirect(w, r, "danger", "Please login first", "/login")
		return
	}

	User := c.appConfig.DbConn.Model("User")
	user := &models.User{}
	User.FindId(session.SessionModel.User).Exec(user)
	data := data{
		Title:    "Profile",
		FlashMsg: session.FlashMsg,
		Data:     user,
	}
	c.render(w, "me", data)
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
